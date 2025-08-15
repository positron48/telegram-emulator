.PHONY: help dev run-backend run-frontend test build clean docker-build docker-run install-deps install-frontend-deps init migrate lint fmt

# Переменные
BINARY_NAME=telegram-emulator
BUILD_DIR=build
WEB_DIR=web

# Цвета для вывода
GREEN=\033[32m
BLUE=\033[34m
GRAY=\033[90m
RESET=\033[0m
BOLD=\033[1m

# =============================================================================
# ИНФОРМАЦИЯ
# =============================================================================

# Команда помощи по умолчанию
.DEFAULT_GOAL := help

help: ## Показать справку по командам
	@printf "$(BOLD)Telegram Emulator - Доступные команды:$(RESET)\n\n"
	@awk 'BEGIN { \
		FS = ":.*?## "; \
		category = ""; \
	} \
	/^# =============================================================================/ { \
		getline; \
		if ($$0 ~ /^# [^=]/) { \
			category = substr($$0, 3); \
			printf "\n$(GREEN)%s$(RESET)\n", category; \
		} \
	} \
	/^[a-zA-Z_-]+:.*?## / { \
		if (category != "") { \
			printf "  $(BLUE)%-20s$(RESET) $(GRAY)%s$(RESET)\n", $$1, $$2; \
		} \
	}' $(MAKEFILE_LIST)

# =============================================================================
# РАЗРАБОТКА
# =============================================================================

dev: ## Запустить backend и frontend в режиме разработки
	@echo "$(BOLD)🚀 Запуск Telegram Emulator в режиме разработки...$(RESET)"
	@echo ""
	@echo "$(GREEN)📋 Статус запуска:$(RESET)"
	@echo ""
	@echo "$(BLUE)1. Backend сервер (Go):$(RESET)"
	@echo "   Запуск на порту 3001..."
	@echo "   API: http://localhost:3001"
	@echo "   WebSocket: ws://localhost:3001/ws"
	@echo ""
	@echo "$(BLUE)2. Frontend сервер (React):$(RESET)"
	@echo "   Запуск на порту 3000..."
	@echo "   Web UI: http://localhost:3000"
	@echo ""
	@echo "$(GRAY)⏳ Запуск сервисов в screen сессиях...$(RESET)"
	@echo ""
	@echo "$(BOLD)Запуск backend в screen сессии 'telegram-backend'...$(RESET)"
	@screen -dmS telegram-backend bash -c "cd $(CURDIR) && go run cmd/emulator/main.go; exec bash"
	@echo "$(GREEN)✅ Backend запущен на http://localhost:3001$(RESET)"
	@echo ""
	@echo "$(BOLD)Запуск frontend в screen сессии 'telegram-frontend'...$(RESET)"
	@screen -dmS telegram-frontend bash -c "cd $(WEB_DIR) && npm run dev; exec bash"
	@echo "$(GREEN)✅ Frontend запущен на http://localhost:3000$(RESET)"
	@echo ""
	@echo "$(BOLD)🎉 Telegram Emulator готов к работе!$(RESET)"
	@echo "$(GRAY)Откройте http://localhost:3000 в браузере$(RESET)"
	@echo ""
	@echo "$(GREEN)📺 Управление screen сессиями:$(RESET)"
	@echo "   screen -r telegram-backend    # Подключиться к backend логам"
	@echo "   screen -r telegram-frontend   # Подключиться к frontend логам"
	@echo "   screen -ls                    # Список всех сессий"
	@echo ""
	@echo "$(GRAY)Для остановки используйте: make stop$(RESET)"
	@echo ""
	@echo "$(BOLD)⏳ Ожидание запуска сервисов...$(RESET)"
	@sleep 3
	@echo "$(GREEN)✅ Все сервисы запущены и готовы к работе!$(RESET)"

run-backend: ## Запустить только backend сервер
	@echo "🚀 Запуск backend сервера..."
	go run cmd/emulator/main.go

run-frontend: ## Запустить только frontend (React dev сервер)
	@echo "🌐 Запуск frontend..."
	cd $(WEB_DIR) && npm run dev

test: ## Запустить тесты
	@echo "🧪 Запуск тестов..."
	go test ./...

build: ## Собрать бинарный файл
	@echo "🔨 Сборка проекта..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/emulator/main.go

clean: ## Очистить build директории
	@echo "🧹 Очистка..."
	rm -rf $(BUILD_DIR)
	go clean

stop: ## Остановить все процессы разработки
	@echo "$(BOLD)🛑 Остановка процессов разработки...$(RESET)"
	@echo ""
	@echo "$(BLUE)Остановка screen сессий...$(RESET)"
	@screen -S telegram-backend -X quit 2>/dev/null || echo "$(GRAY)Backend сессия не найдена$(RESET)"
	@echo "$(GREEN)✅ Backend сессия остановлена$(RESET)"
	@screen -S telegram-frontend -X quit 2>/dev/null || echo "$(GRAY)Frontend сессия не найдена$(RESET)"
	@echo "$(GREEN)✅ Frontend сессия остановлена$(RESET)"
	@echo ""
	@echo "$(BLUE)Очистка процессов...$(RESET)"
	@-pkill -f 'go run cmd/emulator/main.go' 2>/dev/null || echo "$(GRAY)Backend процессы не найдены$(RESET)"
	@-pkill -f 'npm run dev' 2>/dev/null || echo "$(GRAY)Frontend процессы не найдены$(RESET)"
	@-pkill -f 'vite' 2>/dev/null || echo "$(GRAY)Vite процессы не найдены$(RESET)"
	@echo "$(GREEN)✅ Все процессы остановлены$(RESET)"
	@echo ""
	@echo "$(BOLD)🎉 Telegram Emulator полностью остановлен!$(RESET)"

logs-backend: ## Показать логи backend сервера
	@echo "$(BOLD)📺 Подключение к backend логам...$(RESET)"
	@echo "$(GRAY)Для выхода из screen нажмите Ctrl+A, затем D$(RESET)"
	@screen -r telegram-backend

logs-frontend: ## Показать логи frontend сервера
	@echo "$(BOLD)📺 Подключение к frontend логам...$(RESET)"
	@echo "$(GRAY)Для выхода из screen нажмите Ctrl+A, затем D$(RESET)"
	@screen -r telegram-frontend

logs: ## Показать список всех screen сессий
	@echo "$(BOLD)📺 Активные screen сессии:$(RESET)"
	@screen -ls

# =============================================================================
# DOCKER
# =============================================================================

docker-build: ## Собрать Docker образ
	@echo "🐳 Сборка Docker образа..."
	docker build -t $(BINARY_NAME) .

docker-run: ## Запустить Docker контейнер
	@echo "🐳 Запуск Docker контейнера..."
	docker run -p 3001:3001 $(BINARY_NAME)

# =============================================================================
# УСТАНОВКА И НАСТРОЙКА
# =============================================================================

install-deps: ## Установить Go зависимости
	@echo "📦 Установка Go зависимостей..."
	go mod download
	go mod tidy

install-frontend-deps: ## Установить frontend зависимости (npm)
	@echo "📦 Установка frontend зависимостей..."
	cd $(WEB_DIR) && npm install

# Инициализация проекта
init: install-deps install-frontend-deps ## Полная инициализация проекта
	@echo "✅ Проект инициализирован!"

# =============================================================================
# БАЗА ДАННЫХ
# =============================================================================

migrate: ## Запустить миграции базы данных
	@echo "🗄️ Запуск миграций..."
	go run cmd/emulator/main.go migrate

# =============================================================================
# ТЕСТИРОВАНИЕ И КАЧЕСТВО КОДА
# =============================================================================

lint: ## Проверить код линтером
	@echo "🔍 Проверка кода..."
	golangci-lint run

# Форматирование
fmt: ## Форматировать код
	@echo "🎨 Форматирование кода..."
	go fmt ./...
	go vet ./...
