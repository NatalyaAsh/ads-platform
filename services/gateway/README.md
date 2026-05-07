# GraphQL Gateway

Единая точка входа для фронтенда, объединяющая gRPC микросервисы.

## 🚀 Быстрый старт

### Запуск сервиса

```bash
cd services/gateway
go run cmd/server.go
```
## Запуск через Makefile (из корня проекта)
```bash
make run-gateway        # запустить gateway (если добавить в Makefile)
```
## 📦 Конфигурация
Переменные окружения (можно задать в .env):

|Переменная|Значение по умолчанию|Описание|
|-|-|-|
|PORT|8080|Порт HTTP сервера|
|AUTH_SERVICE_ADDR|localhost:50051|Адрес Auth User Service (gRPC)|
|AD_SERVICE_ADDR|localhost:50052|Адрес Ad Search Service (gRPC)|

## 🔌 GraphQL API
Gateway поддерживает следующие типы запросов.

### Query
|Поле|Описание|Аргументы|
|-|-|-|
|ad|Получить объявление по ID|id: ID!|
|categories|Список всех категорий|–|
### Mutation
Поле	Описание	Аргументы
register	Регистрация пользователя	email: String!, password: String!
login	Вход в систему	email: String!, password: String!
createAd	Создать объявление (требуется авторизация)	input: CreateAdInput!
updateAd	Обновить объявление (требуется авторизация)	id: ID!, input: UpdateAdInput!
deleteAd	Удалить объявление (требуется авторизация)	id: ID!
### Входные типы
```bash
input CreateAdInput {
    title: String!
    description: String
    price: Float!
    categoryId: ID!
}

input UpdateAdInput {
    title: String!
    description: String
    price: Float!
    categoryId: ID!
}
```
## Типы ответов
```bash
type AuthPayload {
    userId: ID!
    role: String!
    accessToken: String!
    refreshToken: String!
}

type Ad {
    id: ID!
    title: String!
    description: String
    price: Float!
    userId: ID!
    categoryId: ID!
    category: Category
    status: String!
    views: Int!
    createdAt: String!
    updatedAt: String!
}

type Category {
    id: ID!
    name: String!
    slug: String!
}
```
## 🔐 Авторизация
Для вызовов мутаций createAd, updateAd, deleteAd требуется передавать JWT access token в заголовке Authorization:
```bash
Authorization: Bearer <access_token>
Токен можно получить после успешного login или register.
```
## 📁 Структура проекта
```text
gateway/
├── cmd/
│   └── server.go               # Точка входа
├── pkg/
│   └── grpc/
│       ├── auth/               # gRPC клиент Auth Service
│       │   └── client.go
│       └── ad/                 # gRPC клиент Ad Search Service
│           └── client.go
├── internal/
│   └── pb/                     # Скопированные proto файлы
│       ├── auth/
│       └── ad/
├── go.mod
├── go.sum
└── README.md
```
## 🧪 Примеры запросов
### Регистрация
```go
mutation {
  register(email: "user@example.com", password: "123456") {
    userId
    role
    accessToken
    refreshToken
  }
}
```
### Вход
```go
mutation {
  login(email: "user@example.com", password: "123456") {
    accessToken
  }
}
```
### Создание объявления (с авторизацией)
```go
mutation {
  createAd(input: {
    title: "Ноутбук Apple MacBook Pro",
    description: "16GB RAM, 512GB SSD",
    price: 120000,
    categoryId: "1"
  }) {
    id
    title
    price
    category {
      name
    }
  }
}
```
### Обновление объявления
```go
mutation {
  updateAd(id: "1", input: {
    title: "MacBook Pro (обновлён)",
    price: 115000,
    categoryId: "1"
  }) {
    id
    title
    price
  }
}
```
### Удаление объявления
```go
mutation {
  deleteAd(id: "1")
}
```
### Получение объявления
```go
{
  ad(id: "1") {
    id
    title
    price
    category {
      name
    }
  }
}
```
### Список категорий
```go
{
  categories {
    id
    name
    slug
  }
}
```
## 🐳 Запуск с Docker
Gateway ожидает работающие сервисы auth-user и ad-search. Все сервисы можно поднять через Docker Compose из корня проекта:

```bash
make up
```
## 📊 Зависимости
- Auth User Service (gRPC, порт 50051)
- Ad Search Service (gRPC, порт 50052)
- PostgreSQL (для Auth и Ad, поднимается через Compose)
- MongoDB (для Ad, поднимается через Compose)

## ✅ Статус
Компонент	Статус
GraphQL схема	✅
Query ad	✅
Query categories	✅
Mutation register	✅
Mutation login	✅
Mutation createAd	✅ (с авторизацией)
Mutation updateAd	✅ (с авторизацией)
Mutation deleteAd	✅ (с авторизацией)
Авторизация JWT	✅
gRPC клиенты	✅

## 📝 TODO
- Добавить пагинацию в ads query
- Добавить мутацию uploadMedia (загрузка изображений через MongoDB)
- Добавить админ-мутации (blockUser, moderateAd)
- Интеграция с Elasticsearch + RabbitMQ (поиск)
- Helm чарты для Kubernetes