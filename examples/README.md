# Примеры ботов для Telegram Emulator

[English](README_EN.md) | Русский

Этот каталог содержит примеры ботов, которые демонстрируют работу с Telegram Emulator API.

## Интерактивный Python бот

### Требования
- Python 3.7+
- requests
- flask (для webhook режима)

### Установка зависимостей
```bash
pip install requests flask
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

3. **Запустите интерактивный бот**
   ```bash
   cd examples
   python simple_bot.py
   ```

4. **Выберите режим работы**
   ```
   🤖 Telegram Emulator Bot
   ========================================
   Выберите режим работы бота:
   1. Polling (обычный)
   2. Long Polling (30s)
   3. Webhook
   4. Выход

   Введите номер режима (1-4):
   ```

### Режимы работы

#### 1. Polling (обычный)
- Бот запрашивает обновления каждую секунду
- Подходит для тестирования и разработки
- Использует меньше ресурсов

#### 2. Long Polling (30s)
- Бот ждет новые обновления до 30 секунд
- Более эффективен для продакшена
- Соответствует стандартам Telegram Bot API

#### 3. Webhook
- Эмулятор отправляет обновления на ваш сервер
- Бот запускает встроенный Flask сервер
- Максимальная производительность
- Требует доступный порт (по умолчанию 8080)

### Доступные команды бота
- `/start` - Начать работу с ботом
- `/help` - Показать справку
- `/echo <текст>` - Эхо-ответ

### Особенности

#### Сохранение состояния
- Бот автоматически сохраняет offset в файл
- При перезапуске продолжает с последнего обработанного обновления
- Файл: `bot_offset_{bot_id}.txt`

#### Webhook сервер
- Встроенный Flask сервер
- Автоматическая настройка webhook в эмуляторе
- Graceful shutdown с удалением webhook
- Health check endpoint: `/health`

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
- Для webhook режима проверьте логи Flask сервера

## Поддерживаемые функции

- ✅ Отправка и получение текстовых сообщений
- ✅ Команды бота
- ✅ Polling режим (обычный и long polling)
- ✅ Webhook режим с встроенным сервером
- ✅ Сохранение состояния между запусками
- ✅ Эмуляция статусов сообщений
- ✅ Мультипользовательские чаты
- ✅ Групповые чаты
- ✅ Интерактивный выбор режима работы

## Ограничения

- Медиафайлы (фото, видео, документы) пока не поддерживаются
- Inline кнопки и клавиатуры в разработке
- Некоторые продвинутые функции Telegram API могут не работать
