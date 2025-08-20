## Telegram Emulator Bot Examples

English | [Русский](README.md)

This directory contains example bots that demonstrate how to work with the Telegram Emulator API.

### Requirements
- Python 3.7+
- requests
- flask (for webhook mode)

### Install dependencies
```bash
pip install requests flask
```

### Usage

1. Start the Telegram Emulator
   ```bash
   cd /path/to/telegram-emulator
   make dev
   ```

2. Create a bot in the web UI
   - Open `http://localhost:3001`
   - Go to "Bot Management"
   - Create a new bot with token `1234567890:ABCdefGHIjklMNOpqrsTUVwxyz`

3. Run the interactive bot
   ```bash
   cd examples
   python simple_bot.py
   ```

4. Choose a mode
   ```
   🤖 Telegram Emulator Bot
   ========================================
   Choose bot mode:
   1. Polling (regular)
   2. Long Polling (30s)
   3. Webhook
   4. Exit

   Enter mode number (1-4):
   ```

### Modes

#### 1. Polling (regular)
- The bot fetches updates every second
- Suitable for testing and development
- Uses fewer resources

#### 2. Long Polling (30s)
- The bot waits for new updates for up to 30 seconds
- More efficient for production
- Matches Telegram Bot API standards

#### 3. Webhook
- The emulator sends updates to your server
- The bot runs a built-in Flask server
- Maximum performance
- Requires an available port (8080 by default)

### Available bot commands
- `/start` - Start interacting with the bot
- `/help` - Show help
- `/echo <text>` - Echo reply

### Features

#### State persistence
- The bot automatically saves the offset to a file
- On restart it continues from the last processed update
- File: `bot_offset_{bot_id}.txt`

#### Webhook server
- Built-in Flask server
- Automatic webhook setup in the emulator
- Graceful shutdown with webhook removal
- Health check endpoint: `/health`

## API Structure

The Telegram Emulator supports the following core Telegram Bot API methods:

### Get bot info
```http
GET /bot{token}/getMe
```

### Get updates
```http
GET /bot{token}/getUpdates?offset=0&limit=100&timeout=30
```

### Send a message
```http
POST /bot{token}/sendMessage
Content-Type: application/json

{
  "chat_id": "123456789",
  "text": "Hello, world!",
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

## Build your own bot

1. Create a bot in the emulator
   - Use the web UI to create a bot
   - Save the bot token

2. Use any Telegram Bot API library
   - python-telegram-bot
   - aiogram
   - telebot
   - Or build your own client

3. Configure the emulator URL
   - Replace `api.telegram.org` with `localhost:3001`
   - Use your bot's token

## Example with python-telegram-bot

```python
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters

def start(update, context):
    update.message.reply_text('Hello! I\'m a bot in the emulator.')

def echo(update, context):
    update.message.reply_text(update.message.text)

# Configure the bot for the emulator
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

## Debugging

- Check emulator logs in the console
- Use the web UI to monitor messages
- Check the bot status in the "Bot Management" section
- For webhook mode, check Flask server logs

## Supported features

- ✅ Sending and receiving text messages
- ✅ Bot commands
- ✅ Polling mode (regular and long polling)
- ✅ Webhook mode with built-in server
- ✅ State persistence between runs
- ✅ Message status emulation
- ✅ Multi-user chats
- ✅ Group chats
- ✅ Interactive mode selection

## Limitations

- Media files (photos, videos, documents) are not supported yet
- Inline buttons and keyboards are under development
- Some advanced Telegram API features may not work


