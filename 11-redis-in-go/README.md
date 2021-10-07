# Redis in Go

Изменения по сравнению с прошлой версией сокращателя ссылок:

- Прикрутили кеш к хранилищу сокращенных ссылок
- Прикрутили rate limit на эндпоинты

Полезные ссылки:

- Офицальный сайт Redis --- [redis.io](https://redis.io)
- Мы использовали самый популярный драйвер [go-redis](https://github.com/go-redis/redis).
- Офицальный сайт фзыка программирования Lua --- [lua.org](https://www.lua.org/)
- [Data Types short summary](https://redis.io/topics/data-types) --- какие в Redis есть типы значений помимо строк.
- [Conversion between Lua and Redis data types](https://redis.io/commands/eval#conversion-between-lua-and-redis-data-types) --- 
  тонкости конвертации значений из Redis в Lua и обратно.