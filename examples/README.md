# Примеры ботов для Telegram Emulator

Этот каталог содержит примеры ботов, которые демонстрируют работу с Telegram Emulator API.

## Простой Python бот

### Требования
- Python 3.7+
- requests

### Установка зависимостей
```bash
pip install requests
```

### Использование

1. **Запустите Telegram Emulator**
   ```bash
   cd /path/to/telegram-emulator
   make dev
   ```

2. **Создайте бота в веб-интерфейсе**
   - Откройте http://localhost:3001
   - Перейдите в раздел "Управление ботами"
   - Создайте нового бота с токеном `1234567890:ABCdefGHIjklMNOpqrsTUVwxyz`

3. **Запустите пример бота**
   ```bash
   cd examples
   python simple_bot.py
   ```

4. **Протестируйте бота**
   - Отправьте сообщение в чат с ботом через веб-интерфейс
   - Бот должен ответить на ваше сообщение

### Доступные команды бота
- `/start` - Начать работу с ботом
- `/help` - Показать справку
- `/echo <текст>` - Эхо-ответ

## Структура API

Telegram Emulator поддерживает следующие основные методы Telegram Bot API:

### Получение информации о боте
```http
GET /bot{token}/getMe
```

### Получение обновлений
```http
GET /bot{token}/getUpdates?offset=0&limit=100&timeout=30
```

### Отправка сообщения
```http
POST /bot{token}/sendMessage
Content-Type: application/json

{
  "chat_id": "123456789",
  "text": "Привет, мир!",
  "parse_mode": "HTML"
}
```

### Webhook
```http
POST /bot{token}/setWebhook
Content-Type: application/json

{
  "url": "https://your-server.com/webhook"
}
```

## Создание собственного бота

1. **Создайте бота в эмуляторе**
   - Используйте веб-интерфейс для создания бота
   - Запишите токен бота

2. **Используйте любой Telegram Bot API библиотеку**
   - python-telegram-bot
   - aiogram
   - telebot
   - Или создайте собственный клиент

3. **Настройте URL эмулятора**
   - Замените `api.telegram.org` на `localhost:3001`
   - Используйте токен вашего бота

## Пример с python-telegram-bot

```python
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters

def start(update, context):
    update.message.reply_text('Привет! Я бот в эмуляторе.')

def echo(update, context):
    update.message.reply_text(update.message.text)

# Настройка бота для эмулятора
updater = Updater(
    token='YOUR_BOT_TOKEN',
    base_url='http://localhost:3001',
    use_context=True
)

dispatcher = updater.dispatcher
dispatcher.add_handler(CommandHandler("start", start))
dispatcher.add_handler(MessageHandler(Filters.text & ~Filters.command, echo))

updater.start_polling()
updater.idle()
```

## Отладка

- Проверьте логи эмулятора в консоли
- Используйте веб-интерфейс для мониторинга сообщений
- Проверьте статус бота в разделе "Управление ботами"

## Поддерживаемые функции

- ✅ Отправка и получение текстовых сообщений
- ✅ Команды бота
- ✅ Polling режим
- ✅ Webhook режим
- ✅ Эмуляция статусов сообщений
- ✅ Мультипользовательские чаты
- ✅ Групповые чаты

## Ограничения

- Медиафайлы (фото, видео, документы) пока не поддерживаются
- Inline кнопки и клавиатуры в разработке
- Некоторые продвинутые функции Telegram API могут не работать
