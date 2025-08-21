#!/usr/bin/env python3
"""
Простой бот для тестирования эмулятора Telegram Bot API
Использует официальную библиотеку python-telegram-bot
"""

import asyncio
import logging
from typing import Dict, Any, List

from telegram import Update, InlineKeyboardButton, InlineKeyboardMarkup, ReplyKeyboardMarkup, ReplyKeyboardRemove, ForceReply
from telegram.ext import Application, CommandHandler, MessageHandler, CallbackQueryHandler, filters, ContextTypes

# Конфигурация для эмулятора
BOT_TOKEN = "test-token-123"  # Токен бота в эмуляторе
EMULATOR_URL = "http://localhost:3001"  # URL эмулятора

# Настройка логирования
logging.basicConfig(
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    level=logging.INFO
)
logger = logging.getLogger(__name__)

class EmulatorBot:
    def __init__(self):
        self.application = None
    
    async def start_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Обработка команды /start"""
        user = update.effective_user
        keyboard = [
            [
                InlineKeyboardButton("✅ Да", callback_data="confirm:yes"),
                InlineKeyboardButton("❌ Нет", callback_data="confirm:no")
            ],
            [
                InlineKeyboardButton("🌐 Сайт", url="https://example.com")
            ]
        ]
        reply_markup = InlineKeyboardMarkup(keyboard)
        
        text = f"Привет, {user.first_name}! Я бот для тестирования эмулятора Telegram API.\n\nВыберите опцию:"
        
        await update.message.reply_text(text, reply_markup=reply_markup)
        logger.info(f"Отправлено приветственное сообщение пользователю {user.first_name}")
    
    async def help_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Обработка команды /help"""
        help_text = """
🤖 Доступные команды:

/start - Запуск бота
/help - Показать эту справку
/keyboard - Показать обычную клавиатуру
/remove_keyboard - Убрать клавиатуру
/force_reply - Принудительный ответ
/entities - Демонстрация entities

Поддерживаемые возможности:
• Inline клавиатуры
• Обычные клавиатуры
• Callback queries
• Message entities (команды, упоминания, хештеги)
• Force reply
        """
        
        await update.message.reply_text(help_text)
        logger.info("Отправлена справка")
    
    async def keyboard_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Обработка команды /keyboard"""
        keyboard = [
            ["📱 Главное меню", "⚙️ Настройки"],
            ["❓ Помощь", "📞 Контакты"],
            ["🔙 Назад"]
        ]
        reply_markup = ReplyKeyboardMarkup(keyboard, resize_keyboard=True, one_time_keyboard=False)
        
        await update.message.reply_text("Выберите опцию из клавиатуры:", reply_markup=reply_markup)
        logger.info("Отправлена обычная клавиатура")
    
    async def remove_keyboard_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Обработка команды /remove_keyboard"""
        reply_markup = ReplyKeyboardRemove()
        await update.message.reply_text("Клавиатура удалена!", reply_markup=reply_markup)
        logger.info("Клавиатура удалена")
    
    async def force_reply_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Обработка команды /force_reply"""
        reply_markup = ForceReply()
        await update.message.reply_text("Пожалуйста, введите ваш ответ:", reply_markup=reply_markup)
        logger.info("Отправлен force reply")
    
    async def entities_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Обработка команды /entities"""
        entities_text = """
Демонстрация Message Entities:

1. Команды: /start /help /settings
2. Упоминания: @username @test_user
3. Хештеги: #telegram #bot #api
4. URL: https://core.telegram.org/bots/api

Это сообщение должно содержать entities в ответе API.
        """
        
        await update.message.reply_text(entities_text)
        logger.info("Отправлена демонстрация entities")
    
    async def button_callback(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Обработка callback query"""
        query = update.callback_query
        await query.answer()  # Отвечаем на callback query
        
        user = query.from_user
        data = query.data
        
        logger.info(f"Получен callback query от {user.first_name}: {data}")
        
        # Обрабатываем callback data
        if data == "confirm:yes":
            await query.edit_message_text("✅ Вы выбрали 'Да'!")
        elif data == "confirm:no":
            await query.edit_message_text("❌ Вы выбрали 'Нет'!")
    
    async def handle_text_message(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Обработка текстового сообщения"""
        text = update.message.text
        user = update.effective_user
        
        # Проверяем, является ли сообщение командой
        if text.startswith('/'):
            response = f"Вы отправили команду: {text}\nЭто сообщение должно содержать entity типа 'bot_command'."
        else:
            # Проверяем наличие упоминаний и хештегов
            has_mention = '@' in text
            has_hashtag = '#' in text
            
            response = f"Вы написали: {text}\n\n"
            
            if has_mention:
                response += "✅ Обнаружены упоминания (@username)\n"
            if has_hashtag:
                response += "✅ Обнаружены хештеги (#hashtag)\n"
            
            if not has_mention and not has_hashtag:
                response += "📝 Обычное текстовое сообщение"
        
        await update.message.reply_text(response)
        logger.info(f"Ответ отправлен пользователю {user.first_name}")
    
    async def error_handler(self, update: object, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Обработка ошибок"""
        logger.error(f"Ошибка при обработке обновления: {context.error}")
    
    def setup_handlers(self) -> None:
        """Настройка обработчиков команд и сообщений"""
        # Команды
        self.application.add_handler(CommandHandler("start", self.start_command))
        self.application.add_handler(CommandHandler("help", self.help_command))
        self.application.add_handler(CommandHandler("keyboard", self.keyboard_command))
        self.application.add_handler(CommandHandler("remove_keyboard", self.remove_keyboard_command))
        self.application.add_handler(CommandHandler("force_reply", self.force_reply_command))
        self.application.add_handler(CommandHandler("entities", self.entities_command))
        
        # Callback queries
        self.application.add_handler(CallbackQueryHandler(self.button_callback))
        
        # Текстовые сообщения
        self.application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, self.handle_text_message))
        
        # Обработчик ошибок
        self.application.add_error_handler(self.error_handler)
    
    async def run(self) -> None:
        """Запуск бота"""
        logger.info("🤖 Запуск бота для эмулятора Telegram API")
        logger.info("=" * 50)
        logger.info(f"📡 URL эмулятора: {EMULATOR_URL}")
        logger.info(f"🔑 Токен бота: {BOT_TOKEN}")
        logger.info("=" * 50)
        
        # Создаем приложение с кастомным base_url
        self.application = Application.builder().token(BOT_TOKEN).base_url(f"{EMULATOR_URL}/bot").build()
        
        # Настраиваем обработчики
        self.setup_handlers()
        
        # Получаем информацию о боте
        try:
            bot_info = await self.application.bot.get_me()
            logger.info(f"🤖 Бот: {bot_info.first_name} (@{bot_info.username})")
        except Exception as e:
            logger.error(f"❌ Не удалось получить информацию о боте: {e}")
            return
        
        logger.info("🔄 Бот запущен и ожидает сообщения...")
        
        # Запускаем бота
        await self.application.initialize()
        await self.application.start()
        await self.application.updater.start_polling(allowed_updates=Update.ALL_TYPES)
        
        # Ждем завершения
        try:
            await asyncio.Event().wait()  # Бесконечное ожидание
        except KeyboardInterrupt:
            pass
        finally:
            await self.application.updater.stop()
            await self.application.stop()
            await self.application.shutdown()

async def main():
    """Главная функция"""
    bot = EmulatorBot()
    await bot.run()

if __name__ == "__main__":
    try:
        # Создаем новый event loop
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        loop.run_until_complete(main())
    except KeyboardInterrupt:
        logger.info("🛑 Получен сигнал остановки")
    except Exception as e:
        logger.error(f"❌ Критическая ошибка: {e}")
    finally:
        try:
            loop.close()
        except:
            pass
