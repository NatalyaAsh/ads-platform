# Ad-Search Service

Сервис управления объявлениями и поиска.

## 🚀 Быстрый старт

### Запуск сервиса

```bash
cd services/ad-search
go run cmd/main.go
```
## Запуск тестов
```bash
go test ./... -v
```
## Запуск через Makefile (из корня проекта)
```bash
make run-ad        # запустить сервис
make test-ad       # запустить тесты (если добавить)
```
## 📦 Конфигурация
|Переменная|Значение по умолчанию|Описание|
|-|-|-|
|GRPC_PORT|50052|Порт gRPC сервера|
|DB_HOST|localhost|PostgreSQL хост|
|DB_PORT|5433|PostgreSQL порт|
|DB_USER|postgres|Пользователь БД|
|DB_PASSWORD|postgres|Пароль БД|
|DB_NAME|ad_db|Имя БД|
|MONGODB_URI|mongodb://localhost:27017|MongoDB URI|
|MONGODB_DATABASE|ad_media|Имя БД для медиа|
|ENV|dev|Окружение (dev/stage/prod)|

## 🔌 gRPC API
### Методы объявлений
**Метод	Описание**
CreateAd	Создать объявление
GetAd	Получить объявление по ID
GetUserAds	Получить объявления пользователя
UpdateAd	Обновить объявление
DeleteAd	Удалить объявление
ListAds	Список объявлений с фильтрацией
**Методы категорий**
Метод	Описание
ListCategories	Получить все категории
GetCategory	Получить категорию по ID
## Пример запроса (CreateAd)
```go
CreateAdRequest {
    title: "iPhone 14 Pro"
    description: "Отличное состояние"
    price: 80000.00
    user_id: 1
    category_id: 1
}
```
## 🗄️ База данных
### PostgreSQL (основные данные)
Таблица	Описание
ads	Объявления (id, title, price, status, user_id, category_id)
categories	Категории (id, name, slug)
### MongoDB (медиа-файлы)
Коллекция	Описание
ad_media	Метаданные изображений (ad_id, file_path, file_name)
### Миграции
```bash
# Автоматически накатываются при запуске
# Ручной запуск (из папки сервиса):
migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/ad_db?sslmode=disable" up
```
## 📁 Структура проекта
```text
ad-search/
├── cmd/
│   └── main.go                 # Точка входа (gRPC сервер)
├── internal/
│   ├── model/                  # Модели данных
│   │   ├── ad.go
│   │   ├── category.go
│   │   └── media.go
│   ├── repository/             # Репозитории
│   │   ├── ad_repo.go
│   │   ├── category_repo.go
│   │   ├── media_repo.go
│   │   ├── postgres.go
│   │   └── mongo_client.go
│   ├── service/                # Бизнес-логика
│   │   └── ad_service.go
│   ├── handler/                # gRPC хендлеры
│   │   └── grpc_handler.go
│   └── server/                 # Запуск gRPC сервера
│       └── grpc.go
├── pkg/
│   ├── config/                 # Конфигурация
│   │   └── config.go
│   └── db/                     # Миграции
│       └── migrations.go
├── migrations/                 # SQL миграции
│   ├── 001_create_categories_table.up.sql
│   ├── 001_create_categories_table.down.sql
│   ├── 002_create_ads_table.up.sql
│   └── 002_create_ads_table.down.sql
├── proto_gen/                  # Символическая ссылка на сгенерированные proto
├── uploads/                    # Папка для загруженных файлов
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

# С покрытием
go test ./... -cover
```
## 🐳 Запуск с Docker
Базы данных запускаются через Docker Compose из корня проекта:
```bash
# Из корня проекта
make up

# Или
docker-compose up -d
```
## 📊 Статус
Компонент	Статус
Модели данных	✅
Миграции PostgreSQL	✅
Репозитории	✅
Сервисный слой	✅
MongoDB (медиа)	✅
gRPC API	✅
Тесты	✅
## 🔗 Связь с другими сервисами
GraphQL Gateway — будет вызывать gRPC методы
Auth Service — проверка прав (в будущем)

## 📝 TODO
- Добавить RabbitMQ для событий
- Добавить Elasticsearch для поиска
- Подключить проверку прав через Auth Service