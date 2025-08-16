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

### 🚧 В разработке
- **Создание новых пользователей** через интерфейс
- **Настройки приложения** (панель настроек)
- **Создание новых чатов** через интерфейс
- **Управление ботами** через интерфейс
- **Экспорт/импорт данных**
- **Темы оформления** (светлая/темная)

### 🐛 Известные проблемы
- **WebSocket подключения** - исправлено добавлением user_id параметра
- **Дублирование событий** - исправлено разделением useEffect и правильной очисткой обработчиков
- **Сообщения не появляются в реальном времени** - исправлено отправкой сообщений всем участникам чата, включая отправителя
- **Статусы сообщений не обновляются** - исправлено добавлением обработчика message_status_update
- **Лоадер сообщений не исчезает** - исправлено fallback таймером и правильным обновлением состояния
- **Отсутствует функциональность** кнопки настроек
- **Нет создания пользователей** через интерфейс

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

- `GET /api/chats` - Получить все чаты
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
    "type": "private",
    "title": "Тестовый чат",
    "username": "test_chat",
    "description": "Описание чата",
    "user_ids": ["user1_id", "user2_id"]
  }'
```

### Получение всех пользователей

```bash
curl http://localhost:3001/api/users
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
│   │   └── UserSelector.jsx # Выбор пользователя
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

### Этап 3: Интеграция с ботами 📋
- [ ] Webhook поддержка
- [ ] Прямая интеграция через API
- [ ] Эмуляция Telegram Bot API
- [ ] Тестирование с реальными ботами

### Этап 4: Продвинутые функции 📋
- [ ] Создание пользователей через интерфейс
- [ ] Панель настроек
- [ ] Создание чатов через интерфейс
- [ ] Управление ботами через интерфейс
- [ ] Экспорт/импорт данных
- [ ] Темы оформления
- [ ] Статистика и мониторинг
- [ ] Тестовые сценарии

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

### Известные проблемы и решения
1. **WebSocket подключения** - исправлено добавлением user_id параметра в URL
2. **Дублирование событий в отладке** - исправлено защитой от повторных вызовов useEffect и отключением StrictMode
3. **Отправка сообщений** - теперь полностью через WebSocket
4. **Кнопка настроек** - пока не реализована функциональность
5. **Создание пользователей** - требует реализации модального окна

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
