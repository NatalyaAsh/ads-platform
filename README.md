# Ads Platform

Микросервисная платформа для публикации и поиска объявлений.

## 🏗️ Архитектура

- **Auth User Service** — аутентификация, регистрация, JWT
- **Ad Search Service** — управление объявлениями, категории, медиа (MongoDB)
- **GraphQL Gateway** — единый GraphQL API, авторизация, маршрутизация к gRPC-сервисам
- **Базы данных**: PostgreSQL (постянные данные), MongoDB (метаданные медиа)
- **Фронтенд** — простой интерфейс на Vue 3 (опционально)

## 🛠 Технологии

| Компонент | Технологии |
|-----------|------------|
| **Backend** | Go, gRPC, GraphQL, JWT |
| **Базы данных** | PostgreSQL, MongoDB |
| **API Gateway** | GraphQL (graphql-go) |
| **Фронтенд** | Vue 3, Tailwind CSS, Fetch API |
| **Контейнеризация** | Docker, Docker Compose |

## 🚀 Быстрый старт

### Предварительные требования

- Go 1.26+
- Docker & Docker Compose
- Make (для команд)

### Запуск всех баз данных

```bash
make up
```

### Запуск микросервисов
В отдельных терминалах:
```bash
# Терминал 1: Auth Service
make run-auth

# Терминал 2: Ad Search Service
make run-ad

# Терминал 3: GraphQL Gateway
make run-gateway
```

## Доступ к API
GraphQL Playground: http://localhost:8080/graphql

Auth gRPC: порт 50051

Ad-search gRPC: порт 50052

## Фронтенд
Просто открой frontend/index.html в браузере.

📦 Основные команды Makefile
|Команда|Описание|
|-|-|
|make up|Запустить все базы данных (Docker Compose)|
|make down|Остановить все контейнеры|
|make run-auth|Запустить Auth Service|
|make run-ad|Запустить Ad Search Service|
|make run-gateway|Запустить GraphQL Gateway|
|make test-auth|Запустить тесты Auth Service|
|make test-ad|Запустить тесты Ad Search Service|
|make test-gateway|Запустить тесты Gateway|
|make tidy|Обновить зависимости во всех сервисах|

🗂 Структура проекта
```text
ads-platform/
├── services/
│   ├── auth-user/          # Auth User Service
│   ├── ad-search/          # Ad Search Service
│   └── gateway/            # GraphQL Gateway
├── frontend/               # Vue 3 фронт (index.html)
├── proto/                  # Протофайлы (.proto)
├── docker-compose.yml      # Базы данных
├── Makefile                # Команды для запуска
└── README.md               # Этот файл
```

## 🧪 Тестирование
```bash
# Все тесты
make test-auth
make test-ad
make test-gateway
```

## 📄 Документация по сервисам
- [Auth User Service](./services/auth-user/README.md) — регистрация, логин, JWT, роли
- [Ad Search Service](./services/ad-search/README.md) — объявления, категории, медиа
- [GraphQL Gateway](./services/gateway/README.md) — единая точка входа, авторизация
- [Frontend](./frontend/README.md) — интерфейс на Vue 3

## 💡 Пример GraphQL запроса
```go
mutation {
  login(email: "user@example.com", password: "123456") {
    accessToken
  }
}

mutation {
  createAd(input: {
    title: "iPhone 15 Pro",
    price: 99900,
    categoryId: "1"
  }) {
    id
    title
    price
  }
}
```
