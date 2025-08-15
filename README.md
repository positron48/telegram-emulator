# 🤖 Telegram Emulator

Веб-эмулятор Telegram для локального тестирования и разработки ботов. Эмулятор предоставляет графический интерфейс, имитирующий Telegram, и позволяет тестировать ботов без необходимости использования реального Telegram API.

## 🎯 Возможности

- ✅ **Локальное тестирование** ботов без интернета
- ✅ **Визуальный интерфейс** похожий на Telegram
- ✅ **Полный контроль** над сообщениями и состояниями
- ✅ **Отладка** ботов с детальным логированием
- ✅ **Автоматизированное тестирование** через API
- ✅ **Мультипользовательское тестирование** с несколькими чатами

## 🚀 Быстрый старт

### Предварительные требования

- Go 1.23+
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
```

3. **Запуск эмулятора**
```bash
make dev
```

4. **Открытие в браузере**
```
http://localhost:3001
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

```
telegram-emulator/
├── cmd/emulator/          # Точка входа приложения
├── internal/
│   ├── api/              # HTTP API и обработчики
│   ├── emulator/         # Основная логика эмулятора
│   ├── models/           # Модели данных
│   ├── repository/       # Репозитории для работы с БД
│   ├── websocket/        # WebSocket сервер
│   └── pkg/              # Общие пакеты
├── web/                  # Веб-интерфейс
├── configs/              # Конфигурационные файлы
└── migrations/           # Миграции базы данных
```

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

### Этап 2: WebSocket и UI 🚧
- [ ] WebSocket сервер
- [ ] Базовый веб-интерфейс
- [ ] Компоненты чата
- [ ] Real-time обновления

### Этап 3: Интеграция с ботами 📋
- [ ] Webhook поддержка
- [ ] Прямая интеграция через API
- [ ] Эмуляция Telegram Bot API
- [ ] Тестирование с реальными ботами

### Этап 4: Продвинутые функции 📋
- [ ] Отладочная панель
- [ ] Статистика и мониторинг
- [ ] Тестовые сценарии
- [ ] Документация

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
