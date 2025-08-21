#!/usr/bin/env python3
"""
Simple Telegram Emulator bot example
Demonstrates basic usage of the Telegram Bot API using python-telegram-bot library
"""

import asyncio
import logging
import json
import os
import time
from typing import Dict, Any, List
from flask import Flask, request, jsonify

from telegram import Update, InlineKeyboardButton, InlineKeyboardMarkup, ReplyKeyboardMarkup, ReplyKeyboardRemove
from telegram.ext import Application, CommandHandler, MessageHandler, CallbackQueryHandler, filters, ContextTypes

# Configuration
BOT_TOKEN = "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"  # Bot token in emulator
EMULATOR_URL = "http://localhost:3001"  # Emulator URL
OFFSET_FILE = f"bot_offset_{BOT_TOKEN.split(':')[0]}.txt"

# Setup logging
logging.basicConfig(
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    level=logging.INFO
)
logger = logging.getLogger(__name__)

class TelegramEmulatorBot:
    def __init__(self, token: str, base_url: str = "http://localhost:3001"):
        self.token = token
        self.base_url = base_url
        self.api_url = f"{base_url}/bot{token}"
        self.offset_file = f"bot_offset_{token.split(':')[0]}.txt"
        self.offset = self.load_offset()
        self.application = None
        
    def load_offset(self) -> int:
        """Loads the saved offset from a file"""
        try:
            if os.path.exists(self.offset_file):
                with open(self.offset_file, 'r') as f:
                    return int(f.read().strip())
        except Exception as e:
            print(f"Error loading offset: {e}")
        return 0
    
    def save_offset(self, offset: int) -> None:
        """Saves the offset to a file"""
        try:
            with open(self.offset_file, 'w') as f:
                f.write(str(offset))
        except Exception as e:
            print(f"Error saving offset: {e}")
    
    async def get_me(self) -> Dict[str, Any]:
        """Fetches bot information"""
        try:
            bot_info = await self.application.bot.get_me()
            return {"ok": True, "result": bot_info.to_dict()}
        except Exception as e:
            return {"ok": False, "error_code": 500, "description": str(e)}
    
    async def handle_callback_query(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Handles a callback query from inline buttons"""
        query = update.callback_query
        callback_query_id = query.id
        callback_data = query.data
        message = query.message
        chat_id = message.chat_id
        message_id = message.message_id
        user = query.from_user
        
        print(f"Handling callback_query: {callback_data} from user {user.first_name}")
        print(f"Full callback_query: {query.to_dict()}")
        print(f"DEBUG: chat_id = {chat_id}")
        print(f"DEBUG: message = {message.to_dict()}")
        
        # Ensure callback_data is not empty
        if not callback_data:
            print(f"❌ callback_data is empty or None: {callback_data}")
            if callback_query_id:
                await query.answer("❌ Error: empty callback_data")
            return
        
        # Handle different callback_data types
        if callback_data == 'search':
            print(f"🔍 Handling callback_data 'search'")
            # Show a notification
            await query.answer("🔍 Search in progress...", show_alert=True)
            
            # Send a new message with results
            print(f"Checking chat_id: {chat_id}")
            if chat_id:
                print(f"✅ chat_id found, sending message to chat {chat_id}")
                response_text = "🔍 **Search results:**\n\n✅ Found: 1 result\n⏱️ Search time: 0.1 sec\n📄 Type: text document\n\n_Search completed successfully!_"
                await context.bot.send_message(chat_id=chat_id, text=response_text, parse_mode='Markdown')
                print(f"Send message result: OK")
            else:
                print("❌ chat_id not found, cannot send a message")
            
        elif callback_data == 'notes':
            # Show a notification
            await query.answer("📝 Loading notes...")
            
            # Send a new message
            if chat_id:
                response_text = "📝 **Your notes:**\n\n📌 Note 1: Shopping\n   _Milk, bread, eggs_\n\n📌 Note 2: Meetings\n   _Tomorrow at 15:00_\n\n📌 Note 3: Ideas\n   _New project_\n\n💡 Total notes: 3"
                await context.bot.send_message(chat_id=chat_id, text=response_text, parse_mode='Markdown')
            else:
                print("❌ chat_id not found, cannot send a message")
            
        elif callback_data == 'contacts':
            # Show a notification
            await query.answer("📞 Loading contacts...")
            
            # Send a new message
            if chat_id:
                response_text = "📞 **Support contacts:**\n\n📱 Phone: +7 (999) 123-45-67\n📧 Email: support@example.com\n🤖 Telegram: @support_bot\n\n⏰ Working hours: 24/7\n\n_Feel free to reach out anytime!_"
                await context.bot.send_message(chat_id=chat_id, text=response_text, parse_mode='Markdown')
            else:
                print("❌ chat_id not found, cannot send a message")
            
        else:
            # Unknown callback_data
            await query.answer(f"❓ Unknown command: {callback_data}")
            print(f"Unknown callback_data: {callback_data}")
            
            # Send an error message
            if chat_id:
                response_text = f"❓ **Unknown command:**\n\n🔍 Received: `{callback_data}`\n⚠️ This command is not handled\n\n💡 Try other buttons or the `/help` command"
                await context.bot.send_message(chat_id=chat_id, text=response_text, parse_mode='Markdown')
            else:
                print("❌ chat_id not found, cannot send an error message")
    
    async def start_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Handle /start command"""
        chat_id = update.effective_chat.id
        user = update.effective_user
        text = update.message.text
        
        print(f"Received a message from {user.first_name}: {text}")
        
        response_text = f"Hi! I'm a bot in the Telegram emulator. Your ID: {user.id}"
        
        # Create a keyboard with buttons
        keyboard = [
            [{"text": "ℹ️ Info"}, {"text": "🔧 Settings"}],
            [{"text": "📊 Statistics"}, {"text": "❓ Help"}],
            [{"text": "🎮 Games"}, {"text": "📱 Profile"}]
        ]
        reply_markup = ReplyKeyboardMarkup(keyboard, resize_keyboard=True, one_time_keyboard=False)
        
        await update.message.reply_text(response_text, reply_markup=reply_markup)
        print(f"Reply sent: {response_text}")
    
    async def help_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Handle /help command"""
        chat_id = update.effective_chat.id
        user = update.effective_user
        text = update.message.text
        
        print(f"Received a message from {user.first_name}: {text}")
        
        response_text = "Available commands:\n/start - Start with keyboard\n/help - Help\n/echo <text> - Echo\n/keyboard - Show keyboard\n/inline - Show inline keyboard"
        await update.message.reply_text(response_text)
        print(f"Reply sent: {response_text}")
    
    async def keyboard_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Handle /keyboard command"""
        chat_id = update.effective_chat.id
        user = update.effective_user
        text = update.message.text
        
        print(f"Received a message from {user.first_name}: {text}")
        
        response_text = "Here is a regular keyboard:"
        keyboard = [
            [{"text": "Button 1"}, {"text": "Button 2"}],
            [{"text": "Button 3"}]
        ]
        reply_markup = ReplyKeyboardMarkup(keyboard, resize_keyboard=True)
        
        await update.message.reply_text(response_text, reply_markup=reply_markup)
        print(f"Reply sent: {response_text}")
    
    async def inline_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Handle /inline command"""
        chat_id = update.effective_chat.id
        user = update.effective_user
        text = update.message.text
        
        print(f"Received a message from {user.first_name}: {text}")
        
        response_text = "Here is an inline keyboard:"
        inline_keyboard = [
            [InlineKeyboardButton("🔍 Search", callback_data="search"), InlineKeyboardButton("📝 Notes", callback_data="notes")],
            [InlineKeyboardButton("🌐 Website", url="https://example.com"), InlineKeyboardButton("📞 Contacts", callback_data="contacts")]
        ]
        reply_markup = InlineKeyboardMarkup(inline_keyboard)
        
        await update.message.reply_text(response_text, reply_markup=reply_markup)
        print(f"Reply sent: {response_text}")
    
    async def echo_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Handle /echo command"""
        chat_id = update.effective_chat.id
        user = update.effective_user
        text = update.message.text
        
        print(f"Received a message from {user.first_name}: {text}")
        
        if text.startswith('/echo '):
            echo_text = text[6:]  # Strip '/echo '
            response_text = f"Echo: {echo_text}"
            await update.message.reply_text(response_text)
            print(f"Reply sent: {response_text}")
    
    async def handle_text_message(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Handle regular text messages"""
        chat_id = update.effective_chat.id
        user = update.effective_user
        text = update.message.text
        
        print(f"Received a message from {user.first_name}: {text}")
        
        response_text = f"You wrote: {text}"
        await update.message.reply_text(response_text)
        print(f"Reply sent: {response_text}")
    
    def setup_handlers(self) -> None:
        """Setup command and message handlers"""
        # Commands
        self.application.add_handler(CommandHandler("start", self.start_command))
        self.application.add_handler(CommandHandler("help", self.help_command))
        self.application.add_handler(CommandHandler("keyboard", self.keyboard_command))
        self.application.add_handler(CommandHandler("inline", self.inline_command))
        self.application.add_handler(CommandHandler("echo", self.echo_command))
        
        # Callback queries
        self.application.add_handler(CallbackQueryHandler(self.handle_callback_query))
        
        # Text messages (must be last to catch all non-command text)
        self.application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, self.handle_text_message))
    
    async def run_polling(self, long_polling: bool = True) -> None:
        """Runs the bot in polling mode"""
        mode = "long polling (30s)" if long_polling else "polling"
        print(f"Bot started in {mode} mode...")
        print(f"API URL: {self.api_url}")
        
        # Create application with custom base URL
        self.application = Application.builder().token(self.token).base_url(f"{self.base_url}/bot/").build()
        
        # Setup handlers
        self.setup_handlers()
        
        # Fetch bot info
        me = await self.get_me()
        if me.get('ok'):
            bot_info = me['result']
            print(f"Bot: {bot_info['first_name']} (@{bot_info['username']})")
        else:
            print(f"Failed to get bot info: {me}")
            return
        
        # Start the bot
        await self.application.initialize()
        await self.application.start()
        
        # Configure polling
        if long_polling:
            await self.application.updater.start_polling(timeout=30, allowed_updates=Update.ALL_TYPES)
        else:
            await self.application.updater.start_polling(timeout=0, allowed_updates=Update.ALL_TYPES)
        
        print("🔄 Bot is running and waiting for messages...")
        
        # Keep running
        try:
            await asyncio.Event().wait()  # Infinite wait
        except KeyboardInterrupt:
            print("\nBot stopped")
        finally:
            await self.application.updater.stop()
            await self.application.stop()
            await self.application.shutdown()
    
    async def run_webhook_server(self, port: int = 8080) -> None:
        """Runs a webhook server"""
        print(f"Webhook server started at http://localhost:{port}/webhook")
        print(f"Current offset: {self.get_current_offset()}")
        
        # Create application
        self.application = Application.builder().token(self.token).base_url(f"{self.base_url}/bot/").build()
        
        # Setup handlers
        self.setup_handlers()
        
        # Create Flask app
        app = Flask(__name__)
        
        @app.route('/webhook', methods=['POST'])
        async def webhook():
            try:
                data = request.get_json()
                if data:
                    # Process the update
                    await self.process_webhook_update(data)
                return jsonify({"ok": True})
            except Exception as e:
                print(f"Webhook processing error: {e}")
                return jsonify({"ok": False, "error": str(e)}), 500
        
        @app.route('/health', methods=['GET'])
        def health():
            return jsonify({
                "status": "ok",
                "current_offset": self.get_current_offset(),
                "offset_file": self.offset_file
            })
        
        webhook_url = f"http://localhost:{port}/webhook"
        
        # Set webhook in the emulator
        try:
            await self.application.initialize()
            await self.application.start()
            await self.application.bot.set_webhook(url=webhook_url)
            print(f"Webhook set: {webhook_url}")
        except Exception as e:
            print(f"Webhook setup error: {e}")
        
        try:
            # Note: This is a simplified webhook implementation
            # In a real scenario, you'd need to handle the webhook properly
            print("Webhook server would run here...")
            print("Press Ctrl+C to stop")
            await asyncio.Event().wait()
        except KeyboardInterrupt:
            print("\nWebhook server stopped")
            # Remove webhook on shutdown
            try:
                await self.application.bot.delete_webhook()
            except:
                pass
        finally:
            await self.application.stop()
            await self.application.shutdown()
    
    async def process_webhook_update(self, update: Dict[str, Any]) -> None:
        """Processes an update from webhook"""
        print(f"Received webhook update: {update}")
        
        # Update offset when a webhook update is received
        update_id = update.get('update_id', 0)
        if update_id > 0:
            self.offset = update_id + 1
            self.save_offset(self.offset)
            print(f"DEBUG: webhook offset updated to {self.offset} (update_id: {update_id})")
        
        # Handle messages
        if 'message' in update:
            message = update['message']
            chat_id = message['chat']['id']
            text = message.get('text', '')
            user = message.get('from', {})
            
            print(f"Received a message from {user.get('first_name', 'Unknown')}: {text}")
            
            # Simple bot logic
            if text.lower() == '/start':
                response_text = f"Hi! I'm a bot in the Telegram emulator. Your ID: {user.get('id')}"
                
                # Create a keyboard with buttons
                keyboard = {
                    "keyboard": [
                        [{"text": "ℹ️ Info"}, {"text": "🔧 Settings"}],
                        [{"text": "📊 Statistics"}, {"text": "❓ Help"}],
                        [{"text": "🎮 Games"}, {"text": "📱 Profile"}]
                    ],
                    "resize_keyboard": True,
                    "one_time_keyboard": False
                }
                
                # Send message with keyboard
                await self.application.bot.send_message(
                    chat_id=chat_id, 
                    text=response_text,
                    reply_markup=ReplyKeyboardMarkup.from_dict(keyboard)
                )
                
            elif text.lower() == '/help':
                response_text = "Available commands:\n/start - Start with keyboard\n/help - Help\n/echo <text> - Echo\n/keyboard - Show keyboard\n/inline - Show inline keyboard"
                await self.application.bot.send_message(chat_id=chat_id, text=response_text)
                
            elif text.lower() == '/keyboard':
                response_text = "Here is a regular keyboard:"
                keyboard = {
                    "keyboard": [
                        [{"text": "Button 1"}, {"text": "Button 2"}],
                        [{"text": "Button 3"}]
                    ],
                    "resize_keyboard": True
                }
                await self.application.bot.send_message(
                    chat_id=chat_id, 
                    text=response_text,
                    reply_markup=ReplyKeyboardMarkup.from_dict(keyboard)
                )
                
            elif text.lower() == '/inline':
                response_text = "Here is an inline keyboard:"
                inline_keyboard = {
                    "inline_keyboard": [
                        [{"text": "🔍 Search", "callback_data": "search"}, {"text": "📝 Notes", "callback_data": "notes"}],
                        [{"text": "🌐 Website", "url": "https://example.com"}, {"text": "📞 Contacts", "callback_data": "contacts"}]
                    ]
                }
                await self.application.bot.send_message(
                    chat_id=chat_id, 
                    text=response_text,
                    reply_markup=InlineKeyboardMarkup.from_dict(inline_keyboard)
                )
                
            elif text.lower().startswith('/echo '):
                echo_text = text[6:]  # Strip '/echo '
                response_text = f"Echo: {echo_text}"
                await self.application.bot.send_message(chat_id=chat_id, text=response_text)
                
            else:
                response_text = f"You wrote: {text}"
                await self.application.bot.send_message(chat_id=chat_id, text=response_text)
            
            print(f"Reply sent: {response_text}")
        
        # Handle other update types (callback_query, etc.)
        elif 'callback_query' in update:
            callback_query = update['callback_query']
            print(f"Received callback_query: {callback_query}")
            
            # Handle callback_query
            await self.handle_callback_query(Update.de_json(update, self.application.bot), context=None)
        elif 'edited_message' in update:
            edited_message = update['edited_message']
            print(f"Received edited_message: {edited_message}")
            # You can add edited_message handling here
        else:
            print(f"Unknown update type: {list(update.keys())}")
    
    def get_current_offset(self) -> int:
        """Gets the current offset"""
        return self.offset

async def main():
    """Main function"""
    # Bot token (replace with a real token from the emulator)
    TOKEN = "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"
    
    # Create bot instance
    bot = TelegramEmulatorBot(TOKEN)
    
    print("🤖 Telegram Emulator Bot")
    print("=" * 40)
    print("Choose bot mode:")
    print("1. Polling (regular)")
    print("2. Long Polling (30s)")
    print("3. Webhook")
    print("4. Exit")
    
    while True:
        try:
            choice = input("\nEnter mode number (1-4): ").strip()
            
            if choice == "1":
                print("\n🚀 Starting in Polling mode...")
                await bot.run_polling(long_polling=False)
                break
            elif choice == "2":
                print("\n🚀 Starting in Long Polling mode...")
                await bot.run_polling(long_polling=True)
                break
            elif choice == "3":
                print("\n🚀 Starting in Webhook mode...")
                port = input("Enter port for the webhook server (default 8080): ").strip()
                if not port:
                    port = 8080
                else:
                    port = int(port)
                await bot.run_webhook_server(port=port)
                break
            elif choice == "4":
                print("👋 Goodbye!")
                break
            else:
                print("❌ Invalid choice. Enter a number from 1 to 4.")
                
        except KeyboardInterrupt:
            print("\n👋 Goodbye!")
            break
        except ValueError:
            print("❌ Invalid port format. Using default port 8080.")
            await bot.run_webhook_server(port=8080)
            break
        except Exception as e:
            print(f"❌ Error: {e}")

if __name__ == "__main__":
    try:
        # Create new event loop
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        loop.run_until_complete(main())
    except KeyboardInterrupt:
        logger.info("🛑 Received stop signal")
    except Exception as e:
        logger.error(f"❌ Critical error: {e}")
    finally:
        try:
            loop.close()
        except:
            pass
