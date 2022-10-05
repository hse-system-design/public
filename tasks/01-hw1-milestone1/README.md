# Домашнее задание №1. Часть 1.

**Решения до 2 октября 23:59:59.(9) МСК включительно**.

Моментом сдачи считается timestamp успешного решения в Jenkins. 

## Задача

Разработать web-сервис, реализующий HTTP API с функциональностью минималистичного микроблога типа Twitter.
Сервис должен предоставлять следующий функционал:
- создание нового поста
- получение поста по уникальному идентификатору
- получение всех постов пользователя в обратном хронологическом порядке с пагинацией

Формальное описание API можно найти в [microblog.yaml](./microblog.yaml).

## Рекоммендации

- Для хранения постов можно использовать встроенные в Go структуры данных map и slice,
  а также примитивы синхронизации из пакета ``sync``, например, ``sync.RWMutex``.
  Если у вас есть сложности в реализацией хранилища, рекомендуется сначала сделать задание [golang-task](../00-golang-task).

- Обратите внимание, что последующие домашние задания в этом курсе будут являться доработками данного сервиса,
  поэтому рекомендуем позаботиться о том, чтобы вы сами смогли разобраться в своем коде через 2 недели.

- Не поленитесь написать тесты к своему решению, чтобы не тратить время на ожидание сборки в Jenkins.
  В этом вам помогут библиотечки [stretchr/testify](https://github.com/stretchr/testify) и [getkin/kin-openapi](https://github.com/getkin/kin-openapi#validating-http-requestsresponses).
  При выполнении следующих заданий скажете себе спасибо!

- Сразу выделите интерфейс хранилища постов,
  поскольку в одной из следующих частей этой домашки вам придется перейти с in-memory хранилища
  на использование персистентной базы данных.

## Формат сдачи

- Сервис должен быть оформлен в виде git-репозитория с реализацией на Go.
  В **корне** репозитория должен находиться ``Dockerfile``, описывающий образ сервиса.

- При запуске Docker-образа сервис должен автоматически стартовать на порту,
  указанном в переменной окружения `SERVER_PORT`.

- Критерий сдачи --- успешная сборка Jenkins-джобы http://51.250.71.53/job/hw1-milestone1/
  (логин `student` пароль `student1234!`)

- Успешную сборку нужно сохранить, нажав кнопку "Keep the build forever",
  и загрузить в форму https://forms.gle/Xkgsj9tEu2cncq1D8

- Проверить, что вы успешно отправили посылку, можно по [этой таблице](https://docs.google.com/spreadsheets/d/1549wsZqTzb3p4Vw6c4-cCWeI1T5eR-kx4MUNNKqsECU/edit?usp=sharing)