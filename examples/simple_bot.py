#!/usr/bin/env python3
"""
Simple Telegram Emulator bot example
Demonstrates basic usage of the Telegram Bot API
"""

import requests
import time
import json
import os
from typing import Dict, Any

class TelegramEmulatorBot:
    def __init__(self, token: str, base_url: str = "http://localhost:3001"):
        self.token = token
        self.base_url = base_url
        self.api_url = f"{base_url}/bot{token}"
        self.offset_file = f"bot_offset_{token.split(':')[0]}.txt"
        self.offset = self.load_offset()
        
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
        
    def get_me(self) -> Dict[str, Any]:
        """Fetches bot information"""
        response = requests.get(f"{self.api_url}/getMe")
        return response.json()
    
    def get_updates(self, timeout: int = 30) -> Dict[str, Any]:
        """Retrieves updates from the emulator"""
        params = {
            'offset': self.offset,
            'timeout': timeout,
            'limit': 100
        }
        print(f"DEBUG: requesting updates with offset={self.offset}")
        response = requests.get(f"{self.api_url}/getUpdates", params=params)
        return response.json()
    
    def send_message(self, chat_id: str, text: str, parse_mode: str = None) -> Dict[str, Any]:
        """Sends a message"""
        data = {
            'chat_id': chat_id,
            'text': text
        }
        if parse_mode:
            data['parse_mode'] = parse_mode
            
        response = requests.post(f"{self.api_url}/sendMessage", json=data)
        return response.json()
    
    def send_message_with_keyboard(self, chat_id: str, text: str, keyboard: Dict[str, Any], parse_mode: str = None) -> Dict[str, Any]:
        """Sends a message with a regular keyboard"""
        data = {
            'chat_id': chat_id,
            'text': text,
            'reply_markup': keyboard
        }
        if parse_mode:
            data['parse_mode'] = parse_mode
            
        response = requests.post(f"{self.api_url}/sendMessage", json=data)
        return response.json()
    
    def send_message_with_inline_keyboard(self, chat_id: str, text: str, inline_keyboard: Dict[str, Any], parse_mode: str = None) -> Dict[str, Any]:
        """Sends a message with an inline keyboard"""
        data = {
            'chat_id': chat_id,
            'text': text,
            'reply_markup': inline_keyboard
        }
        if parse_mode:
            data['parse_mode'] = parse_mode
            
        response = requests.post(f"{self.api_url}/sendMessage", json=data)
        return response.json()
    
    def answer_callback_query(self, callback_query_id: str, text: str = None, show_alert: bool = False) -> Dict[str, Any]:
        """Answers a callback query (shows a notification)"""
        data = {
            'callback_query_id': callback_query_id
        }
        if text:
            data['text'] = text
        if show_alert:
            data['show_alert'] = show_alert
            
        response = requests.post(f"{self.api_url}/answerCallbackQuery", json=data)
        return response.json()
    
    def edit_message_text(self, chat_id: str, message_id: str, text: str, reply_markup: Dict[str, Any] = None) -> Dict[str, Any]:
        """Edits message text"""
        data = {
            'chat_id': chat_id,
            'message_id': message_id,
            'text': text
        }
        if reply_markup:
            data['reply_markup'] = reply_markup
            
        response = requests.post(f"{self.api_url}/editMessageText", json=data)
        return response.json()
    
    def handle_callback_query(self, callback_query: Dict[str, Any]) -> None:
        """Handles a callback query from inline buttons"""
        callback_query_id = callback_query.get('id')
        callback_data = callback_query.get('data')
        message = callback_query.get('message', {})
        
        # Extract chat_id from message
        chat_id = message.get('chat_id')  # Directly from message, not message.chat.id
        message_id = message.get('id')    # message_id is the message id
        user = callback_query.get('from', {})
        
        print(f"Handling callback_query: {callback_data} from user {user.get('first_name', 'Unknown')}")
        print(f"Full callback_query: {callback_query}")
        print(f"DEBUG: chat_id = {chat_id}")
        print(f"DEBUG: message = {message}")
        print(f"DEBUG: message.get('chat_id') = {message.get('chat_id')}")
        
        # Ensure callback_data is not empty
        if not callback_data:
            print(f"❌ callback_data is empty or None: {callback_data}")
            if callback_query_id:
                result = self.answer_callback_query(callback_query_id, "❌ Error: empty callback_data")
                print(f"Answer to empty callback_query: {result}")
            return
        
        # Handle different callback_data types
        if callback_data == 'search':
            print(f"🔍 Handling callback_data 'search'")
            # Show a notification
            result = self.answer_callback_query(callback_query_id, "🔍 Search in progress...", show_alert=True)
            print(f"Callback query answer: {result}")
            
            # Send a new message with results
            print(f"Checking chat_id: {chat_id}")
            if chat_id:
                print(f"✅ chat_id found, sending message to chat {chat_id}")
                response = self.send_message(str(chat_id), "🔍 **Search results:**\n\n✅ Found: 1 result\n⏱️ Search time: 0.1 sec\n📄 Type: text document\n\n_Search completed successfully!_")
                print(f"Send message result: {response}")
            else:
                print("❌ chat_id not found, cannot send a message")
            
        elif callback_data == 'notes':
            # Show a notification
            result = self.answer_callback_query(callback_query_id, "📝 Loading notes...")
            print(f"Callback query answer: {result}")
            
            # Send a new message
            if chat_id:
                self.send_message(str(chat_id), "📝 **Your notes:**\n\n📌 Note 1: Shopping\n   _Milk, bread, eggs_\n\n📌 Note 2: Meetings\n   _Tomorrow at 15:00_\n\n📌 Note 3: Ideas\n   _New project_\n\n💡 Total notes: 3")
            else:
                print("❌ chat_id not found, cannot send a message")
            
        elif callback_data == 'contacts':
            # Show a notification
            result = self.answer_callback_query(callback_query_id, "📞 Loading contacts...")
            print(f"Callback query answer: {result}")
            
            # Send a new message
            if chat_id:
                self.send_message(str(chat_id), "📞 **Support contacts:**\n\n📱 Phone: +7 (999) 123-45-67\n📧 Email: support@example.com\n🤖 Telegram: @support_bot\n\n⏰ Working hours: 24/7\n\n_Feel free to reach out anytime!_")
            else:
                print("❌ chat_id not found, cannot send a message")
            
        else:
            # Unknown callback_data
            result = self.answer_callback_query(callback_query_id, f"❓ Unknown command: {callback_data}")
            print(f"Unknown callback_data: {callback_data}")
            
            # Send an error message
            if chat_id:
                self.send_message(str(chat_id), f"❓ **Unknown command:**\n\n🔍 Received: `{callback_data}`\n⚠️ This command is not handled\n\n💡 Try other buttons or the `/help` command")
            else:
                print("❌ chat_id not found, cannot send an error message")

    def set_webhook(self, url: str) -> Dict[str, Any]:
        """Sets a webhook"""
        data = {'url': url}
        response = requests.post(f"{self.api_url}/setWebhook", json=data)
        return response.json()
    
    def delete_webhook(self) -> Dict[str, Any]:
        """Deletes the webhook"""
        response = requests.get(f"{self.api_url}/deleteWebhook")
        return response.json()
    
    def get_webhook_info(self) -> Dict[str, Any]:
        """Gets webhook info"""
        response = requests.get(f"{self.api_url}/getWebhookInfo")
        return response.json()
    
    def get_current_offset(self) -> int:
        """Gets the current offset"""
        return self.offset
    
    def run_webhook_server(self, port: int = 8080) -> None:
        """Runs a webhook server"""
        from flask import Flask, request, jsonify
        
        app = Flask(__name__)
        
        @app.route('/webhook', methods=['POST'])
        def webhook():
            try:
                data = request.get_json()
                if data:
                    # Process the update
                    self.process_webhook_update(data)
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
        print(f"Webhook server started at {webhook_url}")
        print(f"Current offset: {self.get_current_offset()}")
        
        # Set webhook in the emulator
        result = self.set_webhook(webhook_url)
        if result.get('ok'):
            print(f"Webhook set: {webhook_url}")
        else:
            print(f"Webhook setup error: {result}")
        
        try:
            app.run(host='0.0.0.0', port=port, debug=False)
        except KeyboardInterrupt:
            print("\nWebhook server stopped")
            # Remove webhook on shutdown
            self.delete_webhook()
    
    def process_webhook_update(self, update: Dict[str, Any]) -> None:
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
                
                result = self.send_message_with_keyboard(str(chat_id), response_text, keyboard)
            elif text.lower() == '/help':
                response_text = "Available commands:\n/start - Start with keyboard\n/help - Help\n/echo <text> - Echo\n/keyboard - Show keyboard\n/inline - Show inline keyboard"
                result = self.send_message(str(chat_id), response_text)
            elif text.lower() == '/keyboard':
                response_text = "Here is a regular keyboard:"
                keyboard = {
                    "keyboard": [
                        [{"text": "Button 1"}, {"text": "Button 2"}],
                        [{"text": "Button 3"}]
                    ],
                    "resize_keyboard": True
                }
                result = self.send_message_with_keyboard(str(chat_id), response_text, keyboard)
            elif text.lower() == '/inline':
                response_text = "Here is an inline keyboard:"
                inline_keyboard = {
                    "inline_keyboard": [
                        [{"text": "🔍 Search", "callback_data": "search"}, {"text": "📝 Notes", "callback_data": "notes"}],
                        [{"text": "🌐 Website", "url": "https://example.com"}, {"text": "📞 Contacts", "callback_data": "contacts"}]
                    ]
                }
                result = self.send_message_with_inline_keyboard(str(chat_id), response_text, inline_keyboard)
            elif text.lower().startswith('/echo '):
                echo_text = text[6:]  # Strip '/echo '
                response_text = f"Echo: {echo_text}"
                result = self.send_message(str(chat_id), response_text)
            else:
                response_text = f"You wrote: {text}"
                result = self.send_message(str(chat_id), response_text)
            
            if result.get('ok'):
                print(f"Reply sent: {response_text}")
            else:
                print(f"Send error: {result}")
        
        # Handle other update types (callback_query, etc.)
        elif 'callback_query' in update:
            callback_query = update['callback_query']
            print(f"Received callback_query: {callback_query}")
            
            # Handle callback_query
            self.handle_callback_query(callback_query)
        elif 'edited_message' in update:
            edited_message = update['edited_message']
            print(f"Received edited_message: {edited_message}")
            # You can add edited_message handling here
        else:
            print(f"Unknown update type: {list(update.keys())}")
    
    def process_updates(self, updates: list) -> None:
        """Processes received updates"""
        max_update_id = 0
        
        for update in updates:
            update_id = update.get('update_id', 0)
            max_update_id = max(max_update_id, update_id)
            
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
                    
                    result = self.send_message_with_keyboard(str(chat_id), response_text, keyboard)
                elif text.lower() == '/help':
                    response_text = "Available commands:\n/start - Start with keyboard\n/help - Help\n/echo <text> - Echo\n/keyboard - Show keyboard\n/inline - Show inline keyboard"
                    result = self.send_message(str(chat_id), response_text)
                elif text.lower() == '/keyboard':
                    response_text = "Here is a regular keyboard:"
                    keyboard = {
                        "keyboard": [
                            [{"text": "Button 1"}, {"text": "Button 2"}],
                            [{"text": "Button 3"}]
                        ],
                        "resize_keyboard": True
                    }
                    result = self.send_message_with_keyboard(str(chat_id), response_text, keyboard)
                elif text.lower() == '/inline':
                    response_text = "Here is an inline keyboard:"
                    inline_keyboard = {
                        "inline_keyboard": [
                            [{"text": "🔍 Search", "callback_data": "search"}, {"text": "📝 Notes", "callback_data": "notes"}],
                            [{"text": "🌐 Website", "url": "https://example.com"}, {"text": "📞 Contacts", "callback_data": "contacts"}]
                        ]
                    }
                    result = self.send_message_with_inline_keyboard(str(chat_id), response_text, inline_keyboard)
                elif text.lower().startswith('/echo '):
                    echo_text = text[6:]  # Strip '/echo '
                    response_text = f"Echo: {echo_text}"
                    result = self.send_message(str(chat_id), response_text)
                else:
                    response_text = f"You wrote: {text}"
                    result = self.send_message(str(chat_id), response_text)
                if result.get('ok'):
                    print(f"Reply sent: {response_text}")
                else:
                    print(f"Send error: {result}")
            
            # Handle callback query
            elif 'callback_query' in update:
                callback_query = update['callback_query']
                print(f"Received callback_query: {callback_query}")
                
                # Handle callback_query
                self.handle_callback_query(callback_query)
        
        # Update offset to the last processed update_id + 1
        # This is the correct Telegram Bot API logic
        if max_update_id > 0:
            self.offset = max_update_id + 1
            self.save_offset(self.offset)
            print(f"DEBUG: offset updated to {self.offset} (last processed update_id + 1)")
    
    def run_polling(self, long_polling: bool = True) -> None:
        """Runs the bot in polling mode"""
        mode = "long polling (30s)" if long_polling else "polling"
        print(f"Bot started in {mode} mode...")
        print(f"API URL: {self.api_url}")
        
        # Fetch bot info
        me = self.get_me()
        if me.get('ok'):
            bot_info = me['result']
            print(f"Bot: {bot_info['first_name']} (@{bot_info['username']})")
        else:
            print(f"Failed to get bot info: {me}")
            return
        
        while True:
            try:
                # Get updates
                timeout = 30 if long_polling else 0
                updates_response = self.get_updates(timeout=timeout)
                
                if updates_response.get('ok'):
                    updates = updates_response.get('result', [])
                    if updates:
                        print(f"Received {len(updates)} updates")
                        self.process_updates(updates)
                    else:
                        print("No new updates")
                else:
                    print(f"Failed to get updates: {updates_response}")
                
                # Sleep only for regular polling
                if not long_polling:
                    time.sleep(1)
                
            except KeyboardInterrupt:
                print("\nBot stopped")
                break
            except Exception as e:
                print(f"Error: {e}")
                time.sleep(5)  # Pause on error

def main():
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
                bot.run_polling(long_polling=False)
                break
            elif choice == "2":
                print("\n🚀 Starting in Long Polling mode...")
                bot.run_polling(long_polling=True)
                break
            elif choice == "3":
                print("\n🚀 Starting in Webhook mode...")
                port = input("Enter port for the webhook server (default 8080): ").strip()
                if not port:
                    port = 8080
                else:
                    port = int(port)
                bot.run_webhook_server(port=port)
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
            bot.run_webhook_server(port=8080)
            break
        except Exception as e:
            print(f"❌ Error: {e}")

if __name__ == "__main__":
    main()
