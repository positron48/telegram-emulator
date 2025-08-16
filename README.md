# 🤖 Telegram Emulator

Веб-эмулятор Telegram для локального тестирования и разработки ботов. Эмулятор предоставляет графический интерфейс, имитирующий Telegram, и позволяет тестировать ботов без необходимости использования реального Telegram API.

## 🎯 Возможности

### ✅ Реализовано
- **Локальное тестирование** ботов без интернета
- **Современный React интерфейс** в стиле Telegram
- **Real-time обновления** через WebSocket (Gorilla WebSocket)
- **Полный контроль** над сообщениями и состояниями
- **Отладка** ботов с детальным логированием
- **Автоматизированное тестирование** через API
- **Мультипользовательское тестирование** с несколькими чатами
- **Панель отладки** с событиями и статистикой
- **Адаптивный дизайн** для всех устройств
- **Выбор пользователя** из списка доступных
- **Поиск чатов** с фильтрацией
- **Статусы сообщений** (отправлено, доставлено, прочитано)
- **Автоскролл** к новым сообщениям
- **Создание пользователей** через интерфейс (модальное окно)
- **Создание чатов** через интерфейс (приватные, группы, каналы)
- **Управление ботами** через интерфейс (создание, редактирование, удаление)
- **API для ботов** с webhook поддержкой
- **Улучшенная логика чатов** - фильтрация по участникам
- **Удаление пользователей и чатов** через интерфейс
- **Редактирование ботов** (токен, webhook URL, активность)
- **Исправление дублирования сообщений** - оптимистичные обновления
- **Улучшенное WebSocket переподключение** с экспоненциальной задержкой
- **Визуальные улучшения** - компактные чаты, контрастные аватары
- **Автоматическое добавление в чаты** при отправке сообщений
- **Различные типы чатов** с уникальными иконками и поведением
- **Панель настроек** с темами оформления и экспортом/импортом данных
- **Темная и светлая темы** с правильными CSS переменными
- **Поддержка русского и английского языков** с полной локализацией
- **Экспорт/импорт данных** в JSON формате

### 🚧 В разработке
- **Управление участниками чатов** (добавление/удаление)
- **Интеграция с реальными Telegram ботами**
- **Тестовые сценарии** для автоматизированного тестирования
- **Настройки безопасности**
- **Синхронизация настроек**

## 🚀 Быстрый старт

### Предварительные требования

- Go 1.23+
- Node.js 18+
- SQLite (встроен)

### Установка и запуск

1. **Клонирование репозитория**
```bash
git clone <repository-url>
cd telegram-emulator
```

2. **Установка зависимостей**
```bash
make install-deps
make install-frontend-deps
```

3. **Запуск эмулятора**
```bash
make dev
```

4. **Открытие в браузере**
```
http://localhost:3000
```

## 📋 API Endpoints

### Пользователи

- `GET /api/users` - Получить всех пользователей
- `POST /api/users` - Создать пользователя
- `GET /api/users/:id` - Получить пользователя по ID
- `PUT /api/users/:id` - Обновить пользователя
- `DELETE /api/users/:id` - Удалить пользователя
- `GET /api/users/:id/chats` - Получить чаты пользователя

### Чаты

- `GET /api/chats?user_id=<user_id>` - Получить чаты пользователя
- `POST /api/chats` - Создать чат
- `GET /api/chats/:id` - Получить чат по ID
- `PUT /api/chats/:id` - Обновить чат
- `DELETE /api/chats/:id` - Удалить чат
- `GET /api/chats/:id/messages` - Получить сообщения чата
- `POST /api/chats/:id/members` - Добавить участника
- `DELETE /api/chats/:id/members/:userID` - Удалить участника

### Сообщения

- `GET /api/messages/:id` - Получить сообщение по ID
- `PUT /api/messages/:id/status` - Обновить статус сообщения

### 🆕 Боты

- `GET /api/bots` - Получить всех ботов
- `POST /api/bots` - Создать бота
- `GET /api/bots/:id` - Получить бота по ID
- `PUT /api/bots/:id` - Обновить бота
- `DELETE /api/bots/:id` - Удалить бота
- `POST /api/bots/:id/sendMessage` - Отправить сообщение через бота
- `GET /api/bots/:id/updates` - Получить обновления для бота
- `POST /api/bots/:id/webhook` - Webhook для бота

## 🔧 Примеры использования

### Создание пользователя

```bash
curl -X POST http://localhost:3001/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_user",
    "first_name": "Тестовый",
    "last_name": "Пользователь",
    "is_bot": false
  }'
```

### Создание чата

```bash
curl -X POST http://localhost:3001/api/chats \
  -H "Content-Type: application/json" \
  -d '{
    "type": "group",
    "title": "Тестовая группа",
    "user_ids": ["user1_id", "user2_id"]
  }'
```

### 🆕 Создание бота

```bash
curl -X POST http://localhost:3001/api/bots \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Bot",
    "username": "test_bot",
    "token": "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz",
    "webhook_url": "http://localhost:8080/webhook"
  }'
```

### Получение чатов пользователя

```bash
curl "http://localhost:3001/api/chats?user_id=user_id"
```

## 🏗️ Архитектура

### Backend (Go)
```
telegram-emulator/
├── cmd/emulator/          # Точка входа приложения
├── internal/
│   ├── api/              # HTTP API и обработчики
│   ├── emulator/         # Основная логика эмулятора
│   ├── models/           # Модели данных
│   ├── repository/       # Репозитории для работы с БД
│   ├── websocket/        # WebSocket сервер (Gorilla WebSocket)
│   └── pkg/              # Общие пакеты
├── configs/              # Конфигурационные файлы
└── migrations/           # Миграции базы данных
```

### Frontend (React)
```
web/
├── src/
│   ├── components/       # React компоненты
│   │   ├── Sidebar.jsx   # Боковая панель с чатами
│   │   ├── ChatWindow.jsx # Окно чата
│   │   ├── MessageBubble.jsx # Пузырьки сообщений
│   │   ├── DebugPanel.jsx # Панель отладки
│   │   ├── UserSelector.jsx # Выбор пользователя
│   │   ├── CreateUserModal.jsx # Создание пользователя
│   │   ├── CreateChatModal.jsx # Создание чата
│   │   ├── BotManager.jsx # Управление ботами
│   │   ├── CreateBotModal.jsx # Создание бота
│   │   ├── EditBotModal.jsx # Редактирование бота
│   │   └── SettingsPanel.jsx # Панель настроек
│   ├── services/         # Сервисы для работы с API
│   │   ├── api.js        # REST API клиент
│   │   └── websocket.js  # WebSocket клиент (нативный WebSocket)
│   ├── store/            # Управление состоянием (Zustand)
│   └── types/            # Типы данных (JSDoc)
├── package.json          # Зависимости
└── vite.config.js        # Конфигурация Vite
```

### Технологический стек
- **Backend**: Go 1.23+, Gin, Gorilla WebSocket, GORM, SQLite
- **Frontend**: React 18, Vite, Tailwind CSS, Zustand, нативный WebSocket
- **Инструменты**: Make, Screen (для управления процессами)

## 🔧 Конфигурация

Основные настройки находятся в файле `configs/config.yaml`:

```yaml
emulator:
  port: 3001
  host: localhost
  debug: true

database:
  url: sqlite:///data/emulator.db
  max_connections: 10

websocket:
  heartbeat_interval: 30s
  max_connections: 1000

bots:
  webhook_timeout: 30s
  max_connections: 100

logging:
  level: debug
  format: console
  file: ""
```

## 🛠️ Команды Make

- `make dev` - Запуск в режиме разработки
- `make run-backend` - Запуск только backend
- `make run-frontend` - Запуск только frontend
- `make test` - Запуск тестов
- `make build` - Сборка проекта
- `make clean` - Очистка
- `make install-deps` - Установка зависимостей
- `make migrate` - Запуск миграций

## 📊 Модели данных

### User
```go
type User struct {
    ID        string    `json:"id"`
    Username  string    `json:"username"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    IsBot     bool      `json:"is_bot"`
    IsOnline  bool      `json:"is_online"`
    LastSeen  time.Time `json:"last_seen"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Chat
```go
type Chat struct {
    ID          string    `json:"id"`
    Type        string    `json:"type"` // private, group, channel
    Title       string    `json:"title"`
    Username    string    `json:"username"`
    Description string    `json:"description"`
    Members     []User    `json:"members"`
    LastMessage *Message  `json:"last_message"`
    UnreadCount int       `json:"unread_count"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Message
```go
type Message struct {
    ID        string    `json:"id"`
    ChatID    string    `json:"chat_id"`
    FromID    string    `json:"from_id"`
    From      User      `json:"from"`
    Text      string    `json:"text"`
    Type      string    `json:"type"` // text, file, voice, photo
    Status    string    `json:"status"` // sending, sent, delivered, read
    IsOutgoing bool     `json:"is_outgoing"`
    Timestamp time.Time `json:"timestamp"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 🆕 Bot
```go
type Bot struct {
    ID         string    `json:"id"`
    Name       string    `json:"name"`
    Username   string    `json:"username"`
    Token      string    `json:"token"`
    WebhookURL string    `json:"webhook_url"`
    IsActive   bool      `json:"is_active"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

## 🔮 Планы развития

### Этап 1: Базовая инфраструктура ✅
- [x] Настройка проекта и зависимостей
- [x] Создание моделей данных
- [x] Настройка базы данных
- [x] Базовый REST API

### Этап 2: WebSocket и UI ✅
- [x] WebSocket сервер (Gorilla WebSocket)
- [x] Современный React интерфейс
- [x] Компоненты чата
- [x] Real-time обновления
- [x] Панель отладки
- [x] Выбор пользователя

### Этап 3: Интеграция с ботами ✅
- [x] Webhook поддержка
- [x] Прямая интеграция через API
- [x] Эмуляция Telegram Bot API
- [x] Тестирование с реальными ботами
- [x] Управление ботами через интерфейс

### Этап 4: Продвинутые функции 📋
- [x] Создание пользователей через интерфейс
- [x] Создание чатов через интерфейс
- [x] Управление ботами через интерфейс
- [x] Удаление пользователей и чатов
- [x] Редактирование ботов
- [x] Исправление дублирования сообщений
- [x] Улучшенное WebSocket переподключение
- [x] Визуальные улучшения интерфейса
- [x] Панель настроек с темами и экспортом/импортом
- [x] Темы оформления (светлая/темная/системная)
- [x] Компактный режим
- [ ] Расширенная статистика и мониторинг
- [ ] Тестовые сценарии
- [ ] Управление участниками чатов

## 🛠️ Разработка

### Команды для разработки
```bash
# Запуск в режиме разработки
make dev

# Просмотр логов
make logs-backend    # Логи backend
make logs-frontend   # Логи frontend
make logs           # Список screen сессий

# Остановка
make stop

# Только backend
make run-backend

# Только frontend
make run-frontend
```

### WebSocket коммуникация
- Backend WebSocket: `ws://localhost:3001/ws?user_id=<user_id>`
- Используется нативный WebSocket (не Socket.io)
- Требуется user_id параметр для подключения к WebSocket
- Real-time обновления сообщений и чатов
- Отправка сообщений через WebSocket

### Структура компонентов
- `Sidebar` - боковая панель с чатами и выбором пользователя
- `ChatWindow` - основное окно чата
- `MessageBubble` - пузырьки сообщений
- `DebugPanel` - панель отладки с событиями
- `UserSelector` - выбор пользователя
- `CreateUserModal` - создание пользователей
- `CreateChatModal` - создание чатов
- `BotManager` - управление ботами

### Известные проблемы и решения
1. **WebSocket подключения** - исправлено добавлением user_id параметра в URL
2. **Дублирование событий в отладке** - исправлено защитой от повторных вызовов useEffect и отключением StrictMode
3. **Отправка сообщений** - исправлено дублирование через оптимистичные обновления
4. **API ботов** - исправлено обновлением URL в api.js
5. **Логика чатов** - добавлена фильтрация по участникам и автоматическое добавление
6. **WebSocket переподключение** - исправлено с экспоненциальной задержкой и очисткой старых соединений
7. **Статус подключения** - исправлено правильным обновлением состояния isConnected
8. **Визуальные проблемы** - исправлены контрастные аватары и компактные чаты

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте ветку для новой функции (`git checkout -b feature/amazing-feature`)
3. Зафиксируйте изменения (`git commit -m 'Add amazing feature'`)
4. Отправьте в ветку (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📄 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 📞 Поддержка

Если у вас есть вопросы или предложения, создайте Issue в репозитории.

---

**Статус проекта**: 🚧 В разработке

Данный эмулятор находится в активной разработке. Приветствуются любые предложения по улучшению функциональности!

## 📚 Дополнительная документация

- [TELEGRAM_EMULATOR_SPECIFICATION.md](TELEGRAM_EMULATOR_SPECIFICATION.md) - Техническая спецификация
- [SETTINGS_PANEL_FEATURES.md](SETTINGS_PANEL_FEATURES.md) - Документация по панели настроек

## 🎯 Последние обновления

### v1.3.0 - Панель настроек и темы оформления
- ✅ Добавлена полнофункциональная панель настроек
- ✅ Реализованы темы оформления (светлая/темная/системная) с правильными CSS переменными
- ✅ Добавлена поддержка русского и английского языков с полной локализацией
- ✅ Реализован экспорт/импорт данных в JSON формате
- ✅ Добавлены настройки WebSocket и сообщений
- ✅ Интеграция настроек с localStorage
- ✅ Удалены неработающие функции (уведомления, компактный режим)

### v1.2.0 - Улучшения интерфейса и стабильности
- ✅ Исправлено дублирование сообщений при отправке
- ✅ Улучшено WebSocket переподключение с экспоненциальной задержкой
- ✅ Добавлено удаление пользователей и чатов через интерфейс
- ✅ Добавлено редактирование ботов (токен, webhook URL)
- ✅ Улучшены визуальные элементы (компактные чаты, контрастные аватары)
- ✅ Исправлен статус подключения WebSocket
- ✅ Добавлено автоматическое добавление в чаты при отправке сообщений

### v1.1.0 - Управление ботами и чатами
- ✅ Добавлено создание пользователей через модальное окно
- ✅ Добавлено создание чатов (приватные, группы, каналы)
- ✅ Добавлено управление ботами через интерфейс
- ✅ Реализован API для ботов с webhook поддержкой
- ✅ Улучшена логика чатов с фильтрацией по участникам

### v1.0.0 - Базовая функциональность
- ✅ Основной интерфейс в стиле Telegram
- ✅ WebSocket для real-time обновлений
- ✅ Панель отладки
- ✅ Выбор пользователя
- ✅ Отправка и получение сообщений
