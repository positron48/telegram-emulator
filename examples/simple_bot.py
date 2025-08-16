#!/usr/bin/env python3
"""
Простой пример бота для Telegram Emulator
Демонстрирует базовую работу с Telegram Bot API
"""

import requests
import time
import json
from typing import Dict, Any

class TelegramEmulatorBot:
    def __init__(self, token: str, base_url: str = "http://localhost:3001"):
        self.token = token
        self.base_url = base_url
        self.api_url = f"{base_url}/bot{token}"
        self.offset = 0
        
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
    
    def process_updates(self, updates: list) -> None:
        """Обрабатывает полученные обновления"""
        for update in updates:
            update_id = update.get('update_id', 0)
            self.offset = max(self.offset, update_id + 1)
            
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
    
    def run_polling(self) -> None:
        """Запускает бота в режиме polling"""
        print("Бот запущен в режиме polling...")
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
                updates_response = self.get_updates(timeout=30)
                
                if updates_response.get('ok'):
                    updates = updates_response.get('result', [])
                    if updates:
                        print(f"Получено {len(updates)} обновлений")
                        self.process_updates(updates)
                    else:
                        print("Новых обновлений нет")
                else:
                    print(f"Ошибка получения обновлений: {updates_response}")
                
                time.sleep(1)  # Небольшая пауза между запросами
                
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
    
    # Запускаем бота
    bot.run_polling()

if __name__ == "__main__":
    main()
