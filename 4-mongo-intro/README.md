# Материалы к 5-му занятию

## Пагинация и MongoDB

На семинаре мы обсуждали разницу между двумя форматами пагинации:

- явная нумерация страниц в стиле `/api/v1/users/{userId}/posts?page=2&size=10`
- cursor-like схема, в которой вместе с i-й страницей возвращается pageToken, по котороу можно получить (i+1)-ю:
  ```
  GET /api/v1/users/{userId}/posts?size=10
  ... headers...
  
  200 OK
  ... headers ...
  
  {"posts": [...], nextPage: "biba"}
  
  GET /api/v1/users/{userId}/posts?page=biba&size=10
  ... headers ...
  
  200 OK
  ... headers ...
  
  {"posts": [...], nextPage: "kuka"}
  ```

Мы обсудили два важных недостатка явной нумерации страниц в сценарии использования "скроллинг ленты вниз":

- при публикации нового поста страницы "съезжают", 
  и последний пост с текущей страницы оказывается первым постом следущей страницы.

- для извлечения из базы данных страницы, заданной номером, 
  нам приходится использовать метод курсора `skip` (`OFFSET` в SQL),
  из-за чего базе данных приходится на каждый запрос итерироваться по всем постам пользователя, начиная с самого старого.

По второму недостатку возникли сомнения: высказывалось мнение, что сканирования всех постов автора можно легко избежать:

- найти в B-дереве по составному ключу (`userId`, `time`) за логарифм минимальный элемент с заданным `authorId`
- от него проитерироваться к `(page*size)`-му порядковому элементу среди всех с заданным `authorId` по убыванию `time`   
- начиная с этого элемента выгрузить `size` подряд идущих элементов

Однако, к сожалению, эта процедура слишком сложна для планировщика базы данных. Давайте проверим на примере:

1. Запустим контейнер с MongoDB 4.4:
   ```
   docker run -p 27017:27017 mongo:4.4
   ```

2. Создадим коллекцию постов:
   ```
   db.posts.insert([
       {text: "First post a1", authorId: "author1", "metainfo": {"foo": "bar"}},
       {text: "Second post a1", authorId: "author1", "metainfo": {"foo": "baz"}},
       {text: "First post a2", authorId: "author2", "metainfo": {"foo": "bar"}},
       {text: "Third post a1", authorId: "author1", "metainfo": {"foo": "baz"}},
       {text: "Fourth post a1", authorId: "author1", "metainfo": {"foo": "bar"}},
       {text: "Second post a2", authorId: "author2", "metainfo": {"foo": "baz"}}
   ])
   ```
   
3. Создадим индекс по ключам `authorId` и `_id` (
   мы используем `_id` для сортировки по времени создания, т.к. он имеет тип ObjectId, 
   первые 4 байта которого --- unix timestamp создания документа в Big Endian):
   ```
   db.posts.createIndex({authorId: 1, _id: -1})
   ```
   
   **Обратите внимание**: в индексах по нескольким полям важно направление отдельных полей.
   Так, созданный выше индекс [не получится](https://docs.mongodb.com/v4.4/core/index-compound/#sort-order) 
   использовать для сортировки в таком запросе: 
   ```
   db.posts.find().sort({authorId: 1, _id: 1})
   ```

4. Посмотрим на план исполнения запроса с использованием ``skip``
   (результат `explain` подчищен от нерелевантных данных)

   ```
   db.posts.find({authorId: "author1"}).sort({_id: -1}).skip(2).limit(2).explain()
   
   {
     "queryPlanner" : {
       "winningPlan" : {
         "stage" : "LIMIT",
         "limitAmount" : 2,
         "inputStage" : {
           "stage" : "FETCH",
           "inputStage" : {
             "stage" : "SKIP",
             "skipAmount" : 2, // <--- отбрасывание первых двух 2-х постов
             "inputStage" : {
               "stage" : "IXSCAN",
               "indexName" : "authorId_1__id_-1",
               "direction" : "forward",
               "indexBounds" : {
                 "authorId" : ["[\"author1\", \"author1\"]"], 
                 "_id" : ["[MaxKey, MinKey]"] // <---- итерация по всем постам автора
               }
             }
           }
         }
       }
     }
   }
   ```
   
   Видно, что несмотря на теоретическую возможность избежать сканирования всех постов автора, 
   на практике оно происходит.
   Можно еще запустить `explain` с отображением статистики исполнения. Если отбросить визуальный шум,
   то можно увидеть, что всего при сканировании индекса было использовано 4 ключа --- все посты автора:
   ```
   db.posts.find({authorId: "author1"}).sort({_id: -1}).skip(2).limit(2).explain("executionStats")
   {
     "queryPlanner" : {
       "executionStats": {
         "totalKeysExamined": 4 // <--- по индексу была упорядоченная итерация с начала
       }
     }
   }
   ```
   
   В то же время, если использовать в качестве pageToken идентификатор последнего поста предыдущей страницы,
   можно добиться итерации по индексу только начиная с нужного документа:

   ```
   db.posts.find({authorId: "author1", _id: {$gt: ObjectId("61438b9091af65f80f0a22d9")}}).sort({_id: -1}).limit(2).explain()
   
   {
       "queryPlanner" : {
           "winningPlan" : {
               "stage" : "LIMIT",
               "limitAmount" : 2,
               "inputStage" : {
                   "stage" : "FETCH", // < --- стадия SKIP отсутствует 
                   "inputStage" : {
                       "stage" : "IXSCAN",
                       "indexName" : "authorId_1__id_-1",
                       "direction" : "forward",
                       "indexBounds" : {
                           "authorId" : ["[\"author1\", \"author1\"]"],
                           "_id" : ["[ObjectId('ffffffffffffffffffffffff'), ObjectId('61438b9091af65f80f0a22d9'))"] // < --- итерация только с первого лемента страницы 
                       }
                   }
               }
           }
       }
   }
   ```
   
   и в статистике исполнения также видно разницу:
   ```
   db.posts.find({authorId: "author1", _id: {$gt: ObjectId("61438b9091af65f80f0a22d9")}}).sort({_id: -1}).limit(2).explain("executionStats")
   {
     "queryPlanner" : {
       "executionStats": {
         "totalKeysExamined": 4 // <--- по индексу была упорядоченная итерация с начала
       }
     }
   }
   ```