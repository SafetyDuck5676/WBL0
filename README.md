# WBL0
В сервисе реализовано:

* Подключение и подписка на канал в Nats-streaming
* Сохранение полученных данных в Postgres
* Хранение и выдача данных по id из кеша (с помощью http-сервера)
* В случае падения сервиса Кеш восстанваливается из Postgres
* Сделан простейший интерфейс отображения полученных данных, для их запроса по id.

# Начало работы

### Требования
Установите компилятор Go: https://golang.org/doc/tutorial/getting-started

Также необходим доступ к рабочему Nats server: https://docs.nats.io/running-a-nats-service/introduction/installation

Доступ к Postgres. https://www.postgresql.org/download/

Для удобства, в wbl0-server также существует docker-compose файл, с помощью которого можно сразу поднять при необходимости Nats-Streaming и Postgres. Для взаимодействия с DB также использован Adminer, включенный в этот docker-compose файл, но можно и использовать любой другой (DBeaver к примеру)

### Установка
Склонируйте репозиторий и запустите сервер, который будет доступен по URL http://localhost:8081 в вашем браузере. Если вы используете docker-compose для быстрого старта, обратите внимание на его порты. Если порты не заняты, то Nats и Postgres с Adminer должны работать. Иначе измените порты. Обязательно проверьте также .env файл - в нем находятся данные для Postgres и Nats, включая порты. Это также необходимо сделать, если у вас все установлено локально (вместо docker-compose)

В wbl0-server также есть db.sql, который нужно выполнить в приложенном Admirer или вашем менеджером ДБ по выбору.

### Взаимодействие с сервером

Запустите две программы из этих двух папок, wbl0-server и wbl0-storage, в них лежит main.go. Получить доступ к сайту можно через http://localhost:8081. Введите 1, вы получите нужные данные.


### Завершение работы с сервером
Обязательно завершайте все программы через CTRL+C, "Graceful Finish". Это корректно завершает работу, не выводя проблем в целом.
