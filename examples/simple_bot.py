#!/usr/bin/env python3
"""
Простой пример бота для Telegram Emulator
Демонстрирует базовую работу с Telegram Bot API
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
        """Загружает сохраненный offset из файла"""
        try:
            if os.path.exists(self.offset_file):
                with open(self.offset_file, 'r') as f:
                    return int(f.read().strip())
        except Exception as e:
            print(f"Ошибка загрузки offset: {e}")
        return 0
    
    def save_offset(self, offset: int) -> None:
        """Сохраняет offset в файл"""
        try:
            with open(self.offset_file, 'w') as f:
                f.write(str(offset))
        except Exception as e:
            print(f"Ошибка сохранения offset: {e}")
        
    def get_me(self) -> Dict[str, Any]:
        """Получает информацию о боте"""
        response = requests.get(f"{self.api_url}/getMe")
        return response.json()
    
    def get_updates(self, timeout: int = 30) -> Dict[str, Any]:
        """Получает обновления от эмулятора"""
        params = {
            'offset': self.offset,
            'timeout': timeout,
            'limit': 100
        }
        print(f"DEBUG: запрашиваем обновления с offset={self.offset}")
        response = requests.get(f"{self.api_url}/getUpdates", params=params)
        return response.json()
    
    def send_message(self, chat_id: str, text: str, parse_mode: str = None) -> Dict[str, Any]:
        """Отправляет сообщение"""
        data = {
            'chat_id': chat_id,
            'text': text
        }
        if parse_mode:
            data['parse_mode'] = parse_mode
            
        response = requests.post(f"{self.api_url}/sendMessage", json=data)
        return response.json()
    
    def send_message_with_keyboard(self, chat_id: str, text: str, keyboard: Dict[str, Any], parse_mode: str = None) -> Dict[str, Any]:
        """Отправляет сообщение с обычной клавиатурой"""
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
        """Отправляет сообщение с inline клавиатурой"""
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
        """Отвечает на callback query (показывает уведомление)"""
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
        """Редактирует текст сообщения"""
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
        """Обрабатывает callback query от inline кнопок"""
        callback_query_id = callback_query.get('id')
        callback_data = callback_query.get('data')
        message = callback_query.get('message', {})
        
        # Извлекаем chat_id из message
        chat_id = message.get('chat_id')  # Прямо из message, а не из message.chat.id
        message_id = message.get('id')    # message_id это id сообщения
        user = callback_query.get('from', {})
        
        print(f"Обрабатываем callback_query: {callback_data} от пользователя {user.get('first_name', 'Unknown')}")
        print(f"Полный callback_query: {callback_query}")
        print(f"DEBUG: chat_id = {chat_id}")
        print(f"DEBUG: message = {message}")
        print(f"DEBUG: message.get('chat_id') = {message.get('chat_id')}")
        
        # Проверяем, что callback_data не пустой
        if not callback_data:
            print(f"❌ callback_data пустой или None: {callback_data}")
            if callback_query_id:
                result = self.answer_callback_query(callback_query_id, "❌ Ошибка: пустой callback_data")
                print(f"Ответ на пустой callback_query: {result}")
            return
        
        # Обрабатываем разные типы callback_data
        if callback_data == 'search':
            print(f"🔍 Обрабатываю callback_data 'search'")
            # Показываем уведомление
            result = self.answer_callback_query(callback_query_id, "🔍 Поиск выполняется...", show_alert=True)
            print(f"Ответ на callback_query: {result}")
            
            # Отправляем новое сообщение с результатами
            print(f"Проверяю chat_id: {chat_id}")
            if chat_id:
                print(f"✅ chat_id найден, отправляю сообщение в чат {chat_id}")
                response = self.send_message(str(chat_id), "🔍 **Результаты поиска:**\n\n✅ Найдено: 1 результат\n⏱️ Время поиска: 0.1 сек\n📄 Тип: текстовый документ\n\n_Поиск выполнен успешно!_")
                print(f"Результат отправки сообщения: {response}")
            else:
                print("❌ chat_id не найден, не могу отправить сообщение")
            
        elif callback_data == 'notes':
            # Показываем уведомление
            result = self.answer_callback_query(callback_query_id, "📝 Заметки загружаются...")
            print(f"Ответ на callback_query: {result}")
            
            # Отправляем новое сообщение
            if chat_id:
                self.send_message(str(chat_id), "📝 **Ваши заметки:**\n\n📌 Заметка 1: Покупки\n   _Молоко, хлеб, яйца_\n\n📌 Заметка 2: Встречи\n   _Завтра в 15:00_\n\n📌 Заметка 3: Идеи\n   _Новый проект_\n\n💡 Всего заметок: 3")
            else:
                print("❌ chat_id не найден, не могу отправить сообщение")
            
        elif callback_data == 'contacts':
            # Показываем уведомление
            result = self.answer_callback_query(callback_query_id, "📞 Контакты загружаются...")
            print(f"Ответ на callback_query: {result}")
            
            # Отправляем новое сообщение
            if chat_id:
                self.send_message(str(chat_id), "📞 **Контакты поддержки:**\n\n📱 Телефон: +7 (999) 123-45-67\n📧 Email: support@example.com\n🤖 Telegram: @support_bot\n\n⏰ Время работы: 24/7\n\n_Обращайтесь в любое время!_")
            else:
                print("❌ chat_id не найден, не могу отправить сообщение")
            
        else:
            # Неизвестный callback_data
            result = self.answer_callback_query(callback_query_id, f"❓ Неизвестная команда: {callback_data}")
            print(f"Неизвестный callback_data: {callback_data}")
            
            # Отправляем сообщение об ошибке
            if chat_id:
                self.send_message(str(chat_id), f"❓ **Неизвестная команда:**\n\n🔍 Получено: `{callback_data}`\n⚠️ Эта команда не обрабатывается\n\n💡 Попробуйте другие кнопки или команду `/help`")
            else:
                print("❌ chat_id не найден, не могу отправить сообщение об ошибке")

    def set_webhook(self, url: str) -> Dict[str, Any]:
        """Устанавливает webhook"""
        data = {'url': url}
        response = requests.post(f"{self.api_url}/setWebhook", json=data)
        return response.json()
    
    def delete_webhook(self) -> Dict[str, Any]:
        """Удаляет webhook"""
        response = requests.get(f"{self.api_url}/deleteWebhook")
        return response.json()
    
    def get_webhook_info(self) -> Dict[str, Any]:
        """Получает информацию о webhook"""
        response = requests.get(f"{self.api_url}/getWebhookInfo")
        return response.json()
    
    def get_current_offset(self) -> int:
        """Получает текущий offset"""
        return self.offset
    
    def run_webhook_server(self, port: int = 8080) -> None:
        """Запускает webhook сервер"""
        from flask import Flask, request, jsonify
        
        app = Flask(__name__)
        
        @app.route('/webhook', methods=['POST'])
        def webhook():
            try:
                data = request.get_json()
                if data:
                    # Обрабатываем обновление
                    self.process_webhook_update(data)
                return jsonify({"ok": True})
            except Exception as e:
                print(f"Ошибка обработки webhook: {e}")
                return jsonify({"ok": False, "error": str(e)}), 500
        
        @app.route('/health', methods=['GET'])
        def health():
            return jsonify({
                "status": "ok",
                "current_offset": self.get_current_offset(),
                "offset_file": self.offset_file
            })
        
        webhook_url = f"http://localhost:{port}/webhook"
        print(f"Webhook сервер запущен на {webhook_url}")
        print(f"Текущий offset: {self.get_current_offset()}")
        
        # Устанавливаем webhook в эмуляторе
        result = self.set_webhook(webhook_url)
        if result.get('ok'):
            print(f"Webhook установлен: {webhook_url}")
        else:
            print(f"Ошибка установки webhook: {result}")
        
        try:
            app.run(host='0.0.0.0', port=port, debug=False)
        except KeyboardInterrupt:
            print("\nWebhook сервер остановлен")
            # Удаляем webhook при остановке
            self.delete_webhook()
    
    def process_webhook_update(self, update: Dict[str, Any]) -> None:
        """Обрабатывает обновление из webhook"""
        print(f"Получено webhook обновление: {update}")
        
        # Обновляем offset при получении webhook обновления
        update_id = update.get('update_id', 0)
        if update_id > 0:
            self.offset = update_id + 1
            self.save_offset(self.offset)
            print(f"DEBUG: webhook offset обновлен до {self.offset} (update_id: {update_id})")
        
        # Обрабатываем сообщения
        if 'message' in update:
            message = update['message']
            chat_id = message['chat']['id']
            text = message.get('text', '')
            user = message.get('from', {})
            
            print(f"Получено сообщение от {user.get('first_name', 'Unknown')}: {text}")
            
            # Простая логика бота
            if text.lower() == '/start':
                response_text = f"Привет! Я бот в эмуляторе Telegram. Ваш ID: {user.get('id')}"
                
                # Создаем клавиатуру с кнопками
                keyboard = {
                    "keyboard": [
                        [{"text": "ℹ️ Информация"}, {"text": "🔧 Настройки"}],
                        [{"text": "📊 Статистика"}, {"text": "❓ Помощь"}],
                        [{"text": "🎮 Игры"}, {"text": "📱 Профиль"}]
                    ],
                    "resize_keyboard": True,
                    "one_time_keyboard": False
                }
                
                result = self.send_message_with_keyboard(str(chat_id), response_text, keyboard)
            elif text.lower() == '/help':
                response_text = "Доступные команды:\n/start - Начать с клавиатурой\n/help - Помощь\n/echo <текст> - Эхо\n/keyboard - Показать клавиатуру\n/inline - Показать inline клавиатуру"
                result = self.send_message(str(chat_id), response_text)
            elif text.lower() == '/keyboard':
                response_text = "Вот обычная клавиатура:"
                keyboard = {
                    "keyboard": [
                        [{"text": "Кнопка 1"}, {"text": "Кнопка 2"}],
                        [{"text": "Кнопка 3"}]
                    ],
                    "resize_keyboard": True
                }
                result = self.send_message_with_keyboard(str(chat_id), response_text, keyboard)
            elif text.lower() == '/inline':
                response_text = "Вот inline клавиатура:"
                inline_keyboard = {
                    "inline_keyboard": [
                        [{"text": "🔍 Поиск", "callback_data": "search"}, {"text": "📝 Заметки", "callback_data": "notes"}],
                        [{"text": "🌐 Сайт", "url": "https://example.com"}, {"text": "📞 Контакты", "callback_data": "contacts"}]
                    ]
                }
                result = self.send_message_with_inline_keyboard(str(chat_id), response_text, inline_keyboard)
            elif text.lower().startswith('/echo '):
                echo_text = text[6:]  # Убираем '/echo '
                response_text = f"Эхо: {echo_text}"
                result = self.send_message(str(chat_id), response_text)
            else:
                response_text = f"Вы написали: {text}"
                result = self.send_message(str(chat_id), response_text)
            
            if result.get('ok'):
                print(f"Ответ отправлен: {response_text}")
            else:
                print(f"Ошибка отправки: {result}")
        
        # Обрабатываем другие типы обновлений (callback_query, etc.)
        elif 'callback_query' in update:
            callback_query = update['callback_query']
            print(f"Получен callback_query: {callback_query}")
            
            # Обрабатываем callback_query
            self.handle_callback_query(callback_query)
        elif 'edited_message' in update:
            edited_message = update['edited_message']
            print(f"Получено edited_message: {edited_message}")
            # Здесь можно добавить обработку edited_message
        else:
            print(f"Неизвестный тип обновления: {list(update.keys())}")
    
    def process_updates(self, updates: list) -> None:
        """Обрабатывает полученные обновления"""
        max_update_id = 0
        
        for update in updates:
            update_id = update.get('update_id', 0)
            max_update_id = max(max_update_id, update_id)
            
            # Обрабатываем сообщения
            if 'message' in update:
                message = update['message']
                chat_id = message['chat']['id']
                text = message.get('text', '')
                user = message.get('from', {})
                
                print(f"Получено сообщение от {user.get('first_name', 'Unknown')}: {text}")
                
                # Простая логика бота
                if text.lower() == '/start':
                    response_text = f"Привет! Я бот в эмуляторе Telegram. Ваш ID: {user.get('id')}"
                    
                    # Создаем клавиатуру с кнопками
                    keyboard = {
                        "keyboard": [
                            [{"text": "ℹ️ Информация"}, {"text": "🔧 Настройки"}],
                            [{"text": "📊 Статистика"}, {"text": "❓ Помощь"}],
                            [{"text": "🎮 Игры"}, {"text": "📱 Профиль"}]
                        ],
                        "resize_keyboard": True,
                        "one_time_keyboard": False
                    }
                    
                    result = self.send_message_with_keyboard(str(chat_id), response_text, keyboard)
                elif text.lower() == '/help':
                    response_text = "Доступные команды:\n/start - Начать с клавиатурой\n/help - Помощь\n/echo <текст> - Эхо\n/keyboard - Показать клавиатуру\n/inline - Показать inline клавиатуру"
                    result = self.send_message(str(chat_id), response_text)
                elif text.lower() == '/keyboard':
                    response_text = "Вот обычная клавиатура:"
                    keyboard = {
                        "keyboard": [
                            [{"text": "Кнопка 1"}, {"text": "Кнопка 2"}],
                            [{"text": "Кнопка 3"}]
                        ],
                        "resize_keyboard": True
                    }
                    result = self.send_message_with_keyboard(str(chat_id), response_text, keyboard)
                elif text.lower() == '/inline':
                    response_text = "Вот inline клавиатура:"
                    inline_keyboard = {
                        "inline_keyboard": [
                            [{"text": "🔍 Поиск", "callback_data": "search"}, {"text": "📝 Заметки", "callback_data": "notes"}],
                            [{"text": "🌐 Сайт", "url": "https://example.com"}, {"text": "📞 Контакты", "callback_data": "contacts"}]
                        ]
                    }
                    result = self.send_message_with_inline_keyboard(str(chat_id), response_text, inline_keyboard)
                elif text.lower().startswith('/echo '):
                    echo_text = text[6:]  # Убираем '/echo '
                    response_text = f"Эхо: {echo_text}"
                    result = self.send_message(str(chat_id), response_text)
                else:
                    response_text = f"Вы написали: {text}"
                    result = self.send_message(str(chat_id), response_text)
                if result.get('ok'):
                    print(f"Ответ отправлен: {response_text}")
                else:
                    print(f"Ошибка отправки: {result}")
            
            # Обрабатываем callback query
            elif 'callback_query' in update:
                callback_query = update['callback_query']
                print(f"Получен callback_query: {callback_query}")
                
                # Обрабатываем callback_query
                self.handle_callback_query(callback_query)
        
        # Обновляем offset до последнего обработанного update_id + 1
        # Это правильная логика Telegram Bot API
        if max_update_id > 0:
            self.offset = max_update_id + 1
            self.save_offset(self.offset)
            print(f"DEBUG: offset обновлен до {self.offset} (последний обработанный update_id + 1)")
    
    def run_polling(self, long_polling: bool = True) -> None:
        """Запускает бота в режиме polling"""
        mode = "long polling (30s)" if long_polling else "polling"
        print(f"Бот запущен в режиме {mode}...")
        print(f"API URL: {self.api_url}")
        
        # Получаем информацию о боте
        me = self.get_me()
        if me.get('ok'):
            bot_info = me['result']
            print(f"Бот: {bot_info['first_name']} (@{bot_info['username']})")
        else:
            print(f"Ошибка получения информации о боте: {me}")
            return
        
        while True:
            try:
                # Получаем обновления
                timeout = 30 if long_polling else 0
                updates_response = self.get_updates(timeout=timeout)
                
                if updates_response.get('ok'):
                    updates = updates_response.get('result', [])
                    if updates:
                        print(f"Получено {len(updates)} обновлений")
                        self.process_updates(updates)
                    else:
                        print("Новых обновлений нет")
                else:
                    print(f"Ошибка получения обновлений: {updates_response}")
                
                # Пауза только для обычного polling
                if not long_polling:
                    time.sleep(1)
                
            except KeyboardInterrupt:
                print("\nБот остановлен")
                break
            except Exception as e:
                print(f"Ошибка: {e}")
                time.sleep(5)  # Пауза при ошибке

def main():
    # Токен бота (замените на реальный токен из эмулятора)
    TOKEN = "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"
    
    # Создаем экземпляр бота
    bot = TelegramEmulatorBot(TOKEN)
    
    print("🤖 Telegram Emulator Bot")
    print("=" * 40)
    print("Выберите режим работы бота:")
    print("1. Polling (обычный)")
    print("2. Long Polling (30s)")
    print("3. Webhook")
    print("4. Выход")
    
    while True:
        try:
            choice = input("\nВведите номер режима (1-4): ").strip()
            
            if choice == "1":
                print("\n🚀 Запуск в режиме Polling...")
                bot.run_polling(long_polling=False)
                break
            elif choice == "2":
                print("\n🚀 Запуск в режиме Long Polling...")
                bot.run_polling(long_polling=True)
                break
            elif choice == "3":
                print("\n🚀 Запуск в режиме Webhook...")
                port = input("Введите порт для webhook сервера (по умолчанию 8080): ").strip()
                if not port:
                    port = 8080
                else:
                    port = int(port)
                bot.run_webhook_server(port=port)
                break
            elif choice == "4":
                print("👋 До свидания!")
                break
            else:
                print("❌ Неверный выбор. Введите число от 1 до 4.")
                
        except KeyboardInterrupt:
            print("\n👋 До свидания!")
            break
        except ValueError:
            print("❌ Неверный формат порта. Используется порт по умолчанию 8080.")
            bot.run_webhook_server(port=8080)
            break
        except Exception as e:
            print(f"❌ Ошибка: {e}")

if __name__ == "__main__":
    main()
