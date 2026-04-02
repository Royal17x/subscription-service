# Subscription service

REST-сервис для агрегации данных об онлайн подписках пользователей.

## Стек

- **Go 1.25**
- **PostgreSQL 16**
- **chi** - HTTP роутер
- **pgx/v5** - драйвер PostgreSQL
- **golang-migrate** - миграции БД
- **swaggo/swag** - Swagger документация
- **log/slog** - структурированные логи

## Запуск

### Требования
- Docker
- Docker compose

### Шаги

1. Клонировать реопзиторий:
```bash
git clone https://github.com/Royal17x/subscription-service
cd subscription-service
```
2. Создать `.env` файл из примера:
```bash
cp .env.example .env
```

3. Запустить:
```bash
docker compose up --build
```

Сервис будет доступен на `http://localhost:8080`

Swagger UI: `http://localhost:8080/swagger/index.html`

## API

| Метод | Путь | Описание |
|-------|------|----------|
| POST | /api/v1/subscriptions | Создать подписку |
| GET | /api/v1/subscriptions | Список подписок |
| GET | /api/v1/subscriptions/{id} | Получить по ID |
| PUT | /api/v1/subscriptions/{id} | Обновить |
| DELETE | /api/v1/subscriptions/{id} | Удалить |
| GET | /api/v1/subscriptions/total-cost | Суммарная стоимость за период |

### Примеры запросов

Создать подписку:
```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

Получить суммарную стоимость за период:
```bash
curl "http://localhost:8080/api/v1/subscriptions/total-cost?date_from=01-2025&date_to=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
```

## Конфигурация

| Переменная | Описание | Пример |
|------------|----------|--------|
| APP_PORT | Порт сервера | 8080 |
| APP_ENV | Окружение | local |
| DB_HOST | Хост БД | postgres |
| DB_PORT | Порт БД | 5432 |
| DB_USER | Пользователь БД | postgres |
| DB_PASSWORD | Пароль БД | postgres |
| DB_NAME | Имя БД | subscriptions |
| DB_SSLMODE | SSL режим | disable |