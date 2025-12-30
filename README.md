# ecom_go_test

## Требования

- Go 1.25 или выше

## Установка и запуск

### Локально

1. Клонируйте репозиторий:

`git clone <https://github.com/Pryflas/ecom_go_test.git`
`cd ecom_go_test`

2. Запустите приложение:

`go run .`

Сервер будет доступен по адресу `http://localhost:8080`


### Docker

#### Сборка образа

`docker build -t ecom_go_test .`
#### Запуск контейнера

`docker run -p 8080:8080 ecom_go_test`

## Запуск тестов

`go test -v`

`go test -cover`


## API Эндпоинты

### Создание задачи

**POST** `/todos`

 `curl -X POST http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Моя первая задача","description":"Тестирование API"}'`


### Получение всех задач

**GET** `/todos`

`curl http://localhost:8080/todos`


### Получение задачи по ID

**GET** `/todos/{id}`

`curl http://localhost:8080/todos/1`

### Обновление задачи

**PUT** `/todos/{id}`

 `curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Моя первая задача","description":"НЕ Тестирование API"}'`


### Удаление задачи

**DELETE** `/todos/{id}`

`curl -X DELETE http://localhost:8080/todos/1`


