# OpenAPI

На 4-м занятии мы:

- дописали сервис сокращения ссылок, который начали писать на [3-м занятии](../go-2)
- написали [OpenAPI](https://swagger.io/docs/specification/about/) спецификацию его API
- написали тесты реализации API с помощью библиотечки
[stretchr/testify](https://github.com/stretchr/testify)
- проверили соответствие реализации и спецификации в помощью библиотечки [getkin/kin-openapi](https://github.com/getkin/kin-openapi#validating-http-requestsresponses)

# Полезные ссылки

- полезная [статья](https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702) про концепцию оборачивания обработчиков (aka Middleware). Будет полезно при выполнении большой домашки
- [Effective Go](https://golang.org/doc/effective_go) -- свод "правил" написания идеоматичного кода на Go. В частности, в разделе [Interfaces and methods](https://golang.org/doc/effective_go#interface_methods) разбирается трюк, который я проделал в тестах с `RoundTripperFunc` на примере `http.HandlerFunc`.
