# order-service

Приложение, работающее с сервисом пункта выдачи заказов.

## Инструменты

- Go
- Gin
- SQL
- gRPC
- Docker
- PostgreSQL
- Testcontainers

## Генерация кода из protobuf-файла
```bash
make generate-proto
```

## Запуск приложения

```bash
docker compose up -d --build
```

## Остановка приложения

```bash
docker compose down
```

## Просмотр логов приложения

### rest-app
```bash
docker logs order-service-rest-app-1
```

### grpc-app
```bash
docker logs order-service-grpc-app-1
```

## Очистка базы данных

```bash
docker volume rm order-service_postgres-data
```

## Данные для авторизации в БД

- PostgreSQL, порт - 5432:
    - username: user
    - password: password

## Примечания
- Работу endpoint'ов рекомендуется проверять в Postman.
- Для gRPC сервера включена рефлексия.
- Protobuf-файл можно просмотреть [тут](https://github.com/ners1us/order-service/blob/main/internal/api/grpc/proto/pvz.proto)