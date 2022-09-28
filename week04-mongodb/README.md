# Содержание Занятия

На занятии мы

- поверхностно познакомились с MongoDB и обсудили её основные отличия от популярных реляционных баз данных
- обсудили ключевые моменты API сервиа, который необходимо реализовать в [hw1-milestone1](../tasks/01-hw1-milestone1/microblog.yaml)
- обсудили плюсы и минусы двух разных подходов к реализации пагинации с точки зрения эффективности работы с базой данных
- посмотрели, как работать с MongoDB в Go

# Полезные ссылки

- **[самое важное, что вы должны знать про MongoDB](https://youtu.be/b2F-DItXtZs)**
- сокращатель ссылок, хранящий данные в MongoDB в репозитории [live coding](https://github.com/hse-system-design/live-coding/tree/ab10d180b2f4a4b2468d612828bc74741f2d6993)
  (рекоммендуется знакомиться с кодом в этом репозитории по-коммитно, см. раздел History / git log на вашем компьютере)
- [Go Mongo Tutorial](https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial)
- [don't use "offset"](https://www.youtube.com/watch?v=WDJRRNCGIRs&ab_channel=HusseinNasser) --- 10-минутное видео,
  с объяснением проблемы использования offset в SQL (то же, что мы обсуждали на занятии про пагинацию)
- [jsonb in postgres](https://dzone.com/articles/using-jsonb-in-postgresql-how-to-effectively-store) --- как
  получить в постгресе что-то типа mongo experience.