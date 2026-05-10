# Auth User Service

Сервис аутентификации и управления пользователями.

## 🚀 Быстрый старт

### Запуск сервиса

```bash
cd services/auth-user
go run cmd/main.go

## Запуск тестов
```bash
go test ./... -v
```
## Тестовый gRPC клиент
```bash
go run cmd/client/main.go
```

## 📦 Конфигурация
|Переменная|Значение по умолчанию|Описание
|-|--------|---|
|GRPC_PORT|50051|Порт gRPC сервера|
|DB_HOST|localhost|PostgreSQL хост|
|DB_PORT|5432|PostgreSQL порт|
|DB_USER|postgres|Пользователь БД|
|DB_PASSWORD|postgres|Пароль БД|
|DB_NAME|auth_user_db|Имя БД|
|JWT_ACCESS_SECRET|default-access-secret-change-me|Секрет для access токенов|
|JWT_REFRESH_SECRET|default-refresh-secret-change-me|Секрет для refresh токенов|
|JWT_ACCESS_TTL|24h|Время жизни access токена|
|JWT_REFRESH_TTL|720h|Время жизни refresh токена|
|ENV|dev|Окружение (dev/stage/prod)|

## 🔌 gRPC API
### Методы
|Метод|Описание|
|-|-|
|Register|Регистрация пользователя|
|Login|Вход в систему|
|ValidateToken|Проверка access токена|
|RefreshToken|Обновление пары токенов|
|Logout|Выход (отзыв refresh токена)|
|GetUser|Получение информации о пользователе|

### Примеры запросов
Register:
```protobuf
RegisterRequest {
    email: "user@example.com"
    password: "password123"
}
```
Login:
```protobuf
LoginRequest {
    email: "user@example.com"
    password: "password123"
}
```

## 🗄️ База данных
### Таблицы
|Таблица|Описание|
|-|-|
|users|Пользователи (id, email, password_hash, role, is_blocked)|
|refresh_tokens|Refresh токены (id, user_id, token, expires_at, revoked)|

### Миграции
Миграции накатываются автоматически при запуске сервиса.

Ручной запуск:
```bash
# Накатить все миграции
migrate -path migrations -database "$DSN" up

# Откатить одну миграцию
migrate -path migrations -database "$DSN" down 1
```

## 📁 Структура проекта
```text
auth-user/
├── cmd/
│   ├── main.go           # Точка входа (gRPC сервер)
│   └── client/           # Тестовый gRPC клиент
│       └── main.go
├── internal/
│   ├── model/            # Модели данных
│   │   ├── user.go
│   │   ├── refresh_token.go
│   │   └── errors.go
│   ├── repository/       # Работа с БД
│   │   ├── postgres.go
│   │   ├── user_repo.go
│   │   └── refresh_repo.go
│   ├── service/          # Бизнес-логика
│   │   └── auth_service.go
│   ├── handler/          # gRPC хендлеры
│   │   └── grpc_handler.go
│   └── server/           # Запуск сервера
│       └── grpc.go
├── pkg/
│   ├── config/           # Конфигурация
│   │   └── config.go
│   ├── db/               # Миграции
│   │   └── migrations.go
│   └── jwt/              # JWT утилиты
│       └── jwt.go
├── migrations/           # SQL миграции
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_refresh_tokens_table.up.sql
│   └── 002_create_refresh_tokens_table.down.sql
├── proto_gen/            # Символическая ссылка на сгенерированные proto файлы
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## 🧪 Тестирование
```bash
# Все тесты
go test ./... -v

# Только репозитории
go test ./internal/repository/... -v

# Только сервис
go test ./internal/service/... -v

# Только JWT
go test ./pkg/jwt/... -v

# С покрытием
go test ./... -cover
```

## 🐳 Запуск с Docker
```bash
# Собрать образ
docker build -t auth-user-service .

# Запустить контейнер
docker run -p 50051:50051 --env-file .env auth-user-service
```

## 📊 Зависимости
|Зависимость|Версия|Назначение|
|-|-|-|
|Go|1.26+|Язык программирования|
|PostgreSQL|15+|Хранение данных|
|gRPC|-|RPC фреймворк|
|GORM|v1.25+|ORM для PostgreSQL|
|golang-migrate|v4.17+|Миграции БД|
|JWT|v5|JSON Web Tokens|

