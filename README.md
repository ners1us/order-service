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

## Запуск приложения

```bash
docker compose up -d --build
```

## Остановка приложения

```bash
docker compose down
```

## Просмотр логов приложения

```bash
docker logs order-service-app-1
```

## Очистка базы данных

```bash
docker volume rm order-service_postgres-data
```

## Данные для авторизации в БД

- PostgreSQL, порт - 5432:
    - username: user
    - password: password
