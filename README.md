# order-service

Приложение, работающее с сервисом пункта выдачи заказов (ПВЗ).

## Инструменты

- Go
- Gin
- SQL
- gRPC
- Docker
- PostgreSQL
- Prometheus
- Testcontainers

## Сущности

### 1. **Users** — Пользователи
- **id**: Уникальный идентификатор пользователя.
- **email**: Электронная почта пользователя.
- **password**: Пароль пользователя.
- **role**: Роль пользователя в системе:
  - `employee` — Сотрудник.
  - `moderator` — Модератор.

### 2. **PVZs** — Пункты выдачи заказов (ПВЗ)
- **id**: Уникальный идентификатор ПВЗ.
- **registration_date**: Дата регистрации ПВЗ в системе.
- **city**: Город, в котором расположен ПВЗ:
  - `Москва`
  - `Санкт-Петербург`
  - `Казань`

### 3. **Receptions** — Приемки товаров в ПВЗ
- **id**: Уникальный идентификатор приемки.
- **date_time**: Дата и время создания приемки.
- **pvz_id**: Идентификатор связанного ПВЗ.
- **status**: Статус приемки:
  - `in_progress` — В процессе.
  - `closed` — Завершена.

### 4. **Products** — Товары в приемках
- **id**: Уникальный идентификатор товара.
- **date_time**: Дата и время добавления товара.
- **type**: Тип товара:
  - `электроника`
  - `одежда`
  - `обувь`
- **reception_id**: Идентификатор связанной приемки.

## Серверы

### gRPC (порт: 3000)

gRPC-сервер для получения списка всех ПВЗ.

### Metrics (порт: 9000)

Просмотр Prometheus-метрик для мониторинга операций с ПВЗ, приемками и товарами.

### HTTP (порт: 8080)

HTTP-сервер для управления ПВЗ, приемками, товарами и пользователями с JWT-авторизацией.

## Возможности API

### REST API

- **/register** (POST) — Регистрация нового пользователя.
- **/dummyLogin** (POST) — Получение тестового JWT-токена для роли.
- **/login** (POST) — Авторизация пользователя с выдачей JWT-токена.
- **/receptions** (POST) — Создание новой приемки товаров в ПВЗ (только для сотрудников).
- **/products** (POST) — Добавление товара в текущую приемку (только для сотрудников).
- **/pvz** (POST) — Создание нового ПВЗ (только для модераторов).
- **/pvz** (GET) — Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией.
- **/pvz/{pvzId}/close_last_reception** (POST) — Закрытие последней открытой приемки в ПВЗ (только для сотрудников).
- **/pvz/{pvzId}/delete_last_product** (POST) — Удаление последнего добавленного товара из текущей приемки (только для
  сотрудников).

### gRPC API

- **GetPVZList** — Получение списка всех ПВЗ.

### Metrics

- **/metrics** (GET) — метрики Prometheus: количество созданных ПВЗ, приемок, добавленных товаров, а также другие метрики.

## Команды

### Запуск контейнеров

```bash
make run
```

### Остановка контейнеров

```bash
make stop
```

### Просмотр логов

#### rest-app

```bash
make rest-logs
```

#### grpc-app

```bash
make grpc-logs
```

### Очистка базы данных

```bash
make db-clean
```

### Просмотреть список всех команд
```bash
make help
```

## Данные для авторизации в БД

- PostgreSQL, порт - 5432:
    - name: order-service-db
    - username: user
    - password: password

## Примечания

- Для gRPC сервера включена рефлексия.
- Работу endpoint'ов рекомендуется проверять в Postman.
- Protobuf-файл для сущности **пункта выдачи заказов** можно
  просмотреть [тут](https://github.com/ners1us/order-service/blob/main/internal/api/grpc/proto/pvz.proto).