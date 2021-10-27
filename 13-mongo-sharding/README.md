# Шардирование MongoDB

На семинаре обсуждали шардирование MongoDB.

В папке [src](./src) можно найти код с семинара, в котором в сервис сокращения ссылок,
был добавлен эндпоинт ``/maintenance/createIndices``, который шардирует коллекцию коротких ссылок.

Полезные ссылки:

- [Mongo DB Is Web Scale](https://www.youtube.com/watch?v=b2F-DItXtZs) --- короткий информативный ролик,
  в котором рассматриваются плюсы и минусы шардирования.
- [Mongo Sharding](https://docs.mongodb.com/manual/sharding/) --- офицальная документация MongoDB по теме шардирования
- [docker-compose с шардированной монгой](https://github.com/lfyuomr-gylo/mongodb-sharding-docker-compose)
- [sh.enableSharding](https://docs.mongodb.com/manual/reference/method/sh.enableSharding/) --- 
  команда, разрешающая шардировать коллекции в заданной базе данных
- [sh.shardCollection](https://docs.mongodb.com/manual/reference/method/sh.shardCollection/) ---
  команда включения шардирования конкретной коллекции