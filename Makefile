.PHONY: help
help:
	@echo "Доступные команды:"
	@echo ""
	@echo "🐳 Управление контейнерами:"
	@echo "  make up              - Запустить все контейнеры"
	@echo "  make down            - Остановить все контейнеры"
	@echo "  make down-volumes    - Остановить и удалить данные"
	@echo "  make status          - Проверить статус контейнеров"
	@echo "  make logs            - Посмотреть логи"
	@echo ""
	@echo "🚀 Запуск сервисов:"
	@echo "  make run-auth        - Запустить auth-user сервис"
	@echo "  make run-ad          - Запустить ad-search сервис"
	@echo "  make run-all         - Запустить базы + сервисы"
	@echo ""
	@echo "🧪 Тестирование:"
	@echo "  make test-auth       - Запустить тесты auth-user"
	@echo "  make test-auth-cover - Тесты auth-user с покрытием"
	@echo ""
	@echo "🔧 Утилиты:"
	@echo "  make clean-db        - Очистить БД auth-user"
	@echo "  make gen-proto       - Сгенерировать protobuf код"
	@echo "  make tidy            - Обновить зависимости"
	@echo "  make vet             - Запустить go vet"

# ========== DOCKER COMPOSE ==========

.PHONY: up
up:
	@echo "🐳 Запуск всех контейнеров..."
	docker-compose up -d
	@echo "✅ Контейнеры запущены"

.PHONY: down
down:
	@echo "🛑 Остановка всех контейнеров..."
	docker-compose down
	@echo "✅ Контейнеры остановлены"

.PHONY: down-volumes
down-volumes:
	@echo "⚠️  Остановка контейнеров и удаление данных..."
	docker-compose down -v
	@echo "✅ Контейнеры и данные удалены"

.PHONY: status
status:
	@echo "📊 Статус контейнеров:"
	docker-compose ps

.PHONY: logs
logs:
	@echo "📋 Логи контейнеров:"
	docker-compose logs -f

.PHONY: restart
restart:
	@echo "🔄 Перезапуск контейнеров..."
	docker-compose restart
	@echo "✅ Контейнеры перезапущены"

# ========== AUTH-USER SERVICE ==========

.PHONY: run-auth
run-auth:
	@echo "🚀 Запуск auth-user сервиса..."
	cd services/auth-user && go run cmd/main.go

.PHONY: test-auth
test-auth:
	@echo "🧪 Запуск тестов auth-user..."
	cd services/auth-user && go test ./... -v

.PHONY: test-auth-cover
test-auth-cover:
	@echo "📊 Запуск тестов auth-user с покрытием..."
	cd services/auth-user && go test ./... -cover

.PHONY: build-auth
build-auth:
	@echo "🔨 Сборка auth-user сервиса..."
	cd services/auth-user && go build -o auth-user cmd/main.go

.PHONY: clean-db
clean-db:
	@echo "🗑️ Очистка БД auth-user..."
	docker exec -it test-postgres psql -U postgres -d auth_user_db -c "DELETE FROM refresh_tokens; DELETE FROM users;"
	@echo "✅ БД очищена"

# ========== AD-SEARCH SERVICE ==========

.PHONY: run-ad
run-ad:
	@echo "🚀 Запуск ad-search сервиса..."
	cd services/ad-search && go run cmd/main.go

# ========== BOTH SERVICES ==========

.PHONY: run-all
run-all: up
	@echo "🚀 Базы данных запущены"
	@echo "Теперь можно запускать сервисы:"
	@echo "  make run-auth"
	@echo "  make run-ad"

# ========== PROTO ==========

.PHONY: gen-proto
gen-proto:
	@echo "🔧 Генерация protobuf кода..."
	./proto/generate.sh

# ========== DEPENDENCIES ==========

.PHONY: tidy
tidy:
	@echo "📦 Обновление зависимостей..."
	cd services/auth-user && go mod tidy
	cd services/ad-search && go mod tidy
	@echo "✅ Готово"

# ========== LINT ==========

.PHONY: vet
vet:
	@echo "🔍 Запуск go vet..."
	cd services/auth-user && go vet ./...
	cd services/ad-search && go vet ./...