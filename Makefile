.PHONY: help
help:
	@echo "Доступные команды:"
	@echo "  make run-auth         - Запустить auth-user сервис"
	@echo "  make test-auth        - Запустить все тесты auth-user"
	@echo "  make test-auth-cover  - Запустить тесты с покрытием"
	@echo "  make clean-db         - Очистить БД auth-user (удалить все данные)"
	@echo "  make gen-proto        - Сгенерировать protobuf код"
	@echo "  make build-auth       - Собрать бинарник auth-user"
	@echo "  make tidy             - Обновить зависимости всех сервисов"
	@echo "  make lint             - Запустить линтер"

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
	@echo "✅ Готово"

# ========== LINT ==========

.PHONY: lint
lint:
	@echo "🔍 Запуск линтера..."
	cd services/auth-user && golangci-lint run ./...

.PHONY: vet
vet:
	@echo "🔍 Запуск go vet..."
	cd services/auth-user && go vet ./...	