# Telegram Emulator

Веб-эмулятор Telegram для локального тестирования и разработки ботов. Эмулятор предоставляет графический интерфейс, имитирующий Telegram, и позволяет тестировать ботов без необходимости использования реального Telegram API.

## 🚀 Возможности

### ✅ Реализовано
- **Полноценный Telegram Bot API** - совместимый с официальным API
- **Веб-интерфейс** - похожий на Telegram
- **Управление ботами** - создание, редактирование, активация/деактивация
- **Система чатов** - приватные и групповые чаты
- **Отправка сообщений** - текстовые сообщения между пользователями
- **Получение обновлений** - боты получают обновления через polling
- **Отправка сообщений через ботов** - боты могут отвечать на сообщения
- **WebSocket поддержка** - real-time обновления интерфейса
- **База данных SQLite** - хранение всех данных

### 🔧 Telegram Bot API методы

Эмулятор поддерживает следующие методы Telegram Bot API:

#### Основные методы
- `getMe` - получение информации о боте
- `getUpdates` - получение обновлений
- `sendMessage` - отправка сообщения
- `setWebhook` - установка webhook
- `deleteWebhook` - удаление webhook
- `getWebhookInfo` - информация о webhook

#### Поддерживаемые типы обновлений
- Сообщения (`message`)
- Отредактированные сообщения (`edited_message`)
- Callback queries (`callback_query`)
- Inline queries (`inline_query`)
- И другие типы обновлений

## 🛠️ Установка и запуск

### Требования
- Go 1.23+
- Node.js 18+
- SQLite

### Быстрый старт

1. **Клонирование репозитория**
   ```bash
   git clone <repository-url>
   cd telegram-emulator
   ```

2. **Запуск эмулятора**
   ```bash
   make run
   ```

3. **Открытие веб-интерфейса**
   ```
   http://localhost:3001
   ```

### Ручная установка

1. **Backend (Go)**
   ```bash
   cd cmd/emulator
   go run main.go
   ```

2. **Frontend (React)**
   ```bash
   cd web
   npm install
   npm run dev
   ```

## 🤖 Создание и тестирование ботов

### 1. Создание бота через веб-интерфейс

1. Откройте http://localhost:3001
2. Перейдите в раздел "Управление ботами"
3. Нажмите "Создать бота"
4. Заполните форму:
   - **Имя**: Test Bot
   - **Username**: test_bot
   - **Токен**: 1234567890:ABCdefGHIjklMNOpqrsTUVwxyz

### 2. Тестирование через API

#### Получение информации о боте
```bash
curl "http://localhost:3001/bot1234567890:ABCdefGHIjklMNOpqrsTUVwxyz/getMe"
```

#### Получение обновлений
```bash
curl "http://localhost:3001/bot1234567890:ABCdefGHIjklMNOpqrsTUVwxyz/getUpdates"
```

#### Отправка сообщения
```bash
curl -X POST "http://localhost:3001/bot1234567890:ABCdefGHIjklMNOpqrsTUVwxyz/sendMessage" \
  -H "Content-Type: application/json" \
  -d '{"chat_id": "2773246093156", "text": "Привет!"}'
```

### 3. Python бот

Создайте файл `bot.py`:

```python
import requests
import time

class TelegramEmulatorBot:
    def __init__(self, token):
        self.token = token
        self.api_url = f"http://localhost:3001/bot{token}"
        self.offset = 0
    
    def get_updates(self):
        params = {'offset': self.offset, 'limit': 100}
        response = requests.get(f"{self.api_url}/getUpdates", params=params)
        return response.json()
    
    def send_message(self, chat_id, text):
        data = {'chat_id': chat_id, 'text': text}
        response = requests.post(f"{self.api_url}/sendMessage", json=data)
        return response.json()
    
    def run(self):
        print("Бот запущен...")
        while True:
            updates = self.get_updates()
            if updates.get('ok') and updates.get('result'):
                for update in updates['result']:
                    self.offset = max(self.offset, update['update_id'] + 1)
                    
                    if 'message' in update:
                        message = update['message']
                        chat_id = message['chat']['id']
                        text = message.get('text', '')
                        
                        print(f"Получено: {text}")
                        
                        # Простая логика бота
                        if text == '/start':
                            self.send_message(chat_id, "Привет! Я бот в эмуляторе.")
                        else:
                            self.send_message(chat_id, f"Вы написали: {text}")
            
            time.sleep(1)

# Запуск бота
bot = TelegramEmulatorBot("1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
bot.run()
```

### 4. Тестирование полного цикла

1. **Создайте пользователя и чат**
   ```bash
   # Создание пользователя
   curl -X POST http://localhost:3001/api/users \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "first_name": "Test", "last_name": "User"}'
   
   # Создание чата
   curl -X POST http://localhost:3001/api/chats \
     -H "Content-Type: application/json" \
     -d '{"type": "private", "title": "Test Chat", "user_ids": ["USER_ID"]}'
   ```

2. **Отправьте сообщение в чат**
   ```bash
   curl -X POST http://localhost:3001/api/chats/CHAT_ID/messages \
     -H "Content-Type: application/json" \
     -d '{"text": "Привет, бот!", "from_user_id": "USER_ID"}'
   ```

3. **Запустите бота**
   ```bash
   python bot.py
   ```

4. **Проверьте ответы бота в веб-интерфейсе**

## 📁 Структура проекта

```
telegram-emulator/
├── cmd/emulator/          # Точка входа приложения
├── internal/
│   ├── api/              # HTTP API и Telegram Bot API
│   ├── emulator/         # Основная логика эмулятора
│   ├── models/           # Модели данных
│   ├── repository/       # Слой доступа к данным
│   ├── websocket/        # WebSocket сервер
│   └── pkg/              # Общие пакеты
├── web/                  # React фронтенд
├── examples/             # Примеры ботов
├── configs/              # Конфигурационные файлы
└── migrations/           # Миграции базы данных
```

## 🔧 Конфигурация

Основные настройки в `configs/config.yaml`:

```yaml
emulator:
  port: 3001
  host: localhost
  debug: true

database:
  url: sqlite:///data/emulator.db

logging:
  level: debug
  format: console
```

## 🧪 Тестирование

### Запуск тестов
```bash
make test
```

### Тестирование API
```bash
# Тест Telegram Bot API
python examples/simple_bot.py

# Тест веб-интерфейса
open http://localhost:3001
```

## 📚 Документация

- [Спецификация](TELEGRAM_EMULATOR_SPECIFICATION.md) - полная техническая спецификация
- [Примеры ботов](examples/) - примеры использования
- [API документация](docs/api.md) - описание API методов

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📄 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 🆘 Поддержка

Если у вас есть вопросы или проблемы:

1. Проверьте [Issues](https://github.com/your-repo/issues)
2. Создайте новый Issue с описанием проблемы
3. Приложите логи и конфигурацию

---

**Статус**: ✅ Готов к использованию

Telegram Emulator предоставляет полноценную среду для тестирования и разработки Telegram ботов с совместимым API и удобным веб-интерфейсом.
