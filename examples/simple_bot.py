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
            return jsonify({"status": "ok"})
        
        webhook_url = f"http://localhost:{port}/webhook"
        print(f"Webhook сервер запущен на {webhook_url}")
        
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
            elif text.lower() == '/help':
                response_text = "Доступные команды:\n/start - Начать\n/help - Помощь\n/echo <текст> - Эхо"
            elif text.lower().startswith('/echo '):
                echo_text = text[6:]  # Убираем '/echo '
                response_text = f"Эхо: {echo_text}"
            else:
                response_text = f"Вы написали: {text}"
            
            # Отправляем ответ
            result = self.send_message(str(chat_id), response_text)
            if result.get('ok'):
                print(f"Ответ отправлен: {response_text}")
            else:
                print(f"Ошибка отправки: {result}")
    
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
                elif text.lower() == '/help':
                    response_text = "Доступные команды:\n/start - Начать\n/help - Помощь\n/echo <текст> - Эхо"
                elif text.lower().startswith('/echo '):
                    echo_text = text[6:]  # Убираем '/echo '
                    response_text = f"Эхо: {echo_text}"
                else:
                    response_text = f"Вы написали: {text}"
                
                # Отправляем ответ
                result = self.send_message(str(chat_id), response_text)
                if result.get('ok'):
                    print(f"Ответ отправлен: {response_text}")
                else:
                    print(f"Ошибка отправки: {result}")
        
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
