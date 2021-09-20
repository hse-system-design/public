# Mongo in Go

Используем стандартный драйвер для работы с Монгой из Го

Более подробную информацию можно всегда найти на [официальном сайте драйвера](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures)

Помимо этого все упаковано в docker-compose. Чтобы собрать и запустить все приложение, достаточно запустить команду `docker-compose up --build` в директории с кодом.

Проверить, что все успешно собралось и работае можно например через curl:

```bash
KEY=`curl -X POST http://localhost:8080/api/urls -d '{"url": "http://yandex.ru"}' | jq -r ".key"`
curl http://localhost:8080/$KEY -L
```