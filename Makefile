.PHONY: dev run-backend run-frontend test build clean docker-build docker-run

# Переменные
BINARY_NAME=telegram-emulator
BUILD_DIR=build
WEB_DIR=web

# Основные команды
dev: run-backend

run-backend:
	@echo "🚀 Запуск backend сервера..."
	go run cmd/emulator/main.go

run-frontend:
	@echo "🌐 Запуск frontend..."
	cd $(WEB_DIR) && npm run dev

test:
	@echo "🧪 Запуск тестов..."
	go test ./...

build:
	@echo "🔨 Сборка проекта..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/emulator/main.go

clean:
	@echo "🧹 Очистка..."
	rm -rf $(BUILD_DIR)
	go clean

# Docker команды
docker-build:
	@echo "🐳 Сборка Docker образа..."
	docker build -t $(BINARY_NAME) .

docker-run:
	@echo "🐳 Запуск Docker контейнера..."
	docker run -p 3001:3001 $(BINARY_NAME)

# Установка зависимостей
install-deps:
	@echo "📦 Установка Go зависимостей..."
	go mod download
	go mod tidy

install-frontend-deps:
	@echo "📦 Установка frontend зависимостей..."
	cd $(WEB_DIR) && npm install

# Инициализация проекта
init: install-deps install-frontend-deps
	@echo "✅ Проект инициализирован!"

# Миграции
migrate:
	@echo "🗄️ Запуск миграций..."
	go run cmd/emulator/main.go migrate

# Линтинг
lint:
	@echo "🔍 Проверка кода..."
	golangci-lint run

# Форматирование
fmt:
	@echo "🎨 Форматирование кода..."
	go fmt ./...
	go vet ./...
