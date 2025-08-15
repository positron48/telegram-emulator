# Примеры использования Telegram Emulator Web Interface

## 🚀 Быстрый старт

### 1. Установка и запуск

```bash
# Клонирование репозитория
git clone <repository-url>
cd telegram-emulator

# Установка зависимостей
make install-deps
make install-frontend-deps

# Запуск backend и frontend
make dev
```

### 2. Открытие в браузере

Перейдите по адресу: http://localhost:3000

## 📱 Основные функции

### Поиск чатов

Введите текст в поле поиска в боковой панели для фильтрации чатов по названию или имени пользователя.

### Отправка сообщений

1. Выберите чат из списка слева
2. Введите сообщение в поле ввода внизу
3. Нажмите Enter или кнопку отправки

### Просмотр статусов сообщений

- ⏳ **Отправляется** - анимированный индикатор
- ✓ **Отправлено** - серая галочка
- ✓✓ **Доставлено** - две серые галочки  
- ✓✓ **Прочитано** - две синие галочки

### Панель отладки

Нажмите кнопку "Отладка" в нижней части боковой панели для открытия панели отладки с:

- **События** - список всех событий в реальном времени
- **Статистика** - метрики работы приложения

## 🔧 API Примеры

### Создание пользователя

```javascript
import apiService from './services/api';

const newUser = await apiService.createUser({
  username: 'test_user',
  first_name: 'Тестовый',
  last_name: 'Пользователь',
  is_bot: false
});
```

### Создание чата

```javascript
const newChat = await apiService.createChat({
  type: 'private',
  title: 'Тестовый чат',
  user_ids: ['user1_id', 'user2_id']
});
```

### Отправка сообщения

```javascript
const message = await apiService.sendMessage(chatId, {
  text: 'Привет, мир!',
  from_user_id: currentUser.id,
  type: 'text'
});
```

### Получение сообщений чата

```javascript
const messages = await apiService.getChatMessages(chatId, 50, 0);
```

## 🔌 WebSocket Примеры

### Подключение к WebSocket

```javascript
import wsService from './services/websocket';

// Подключение
await wsService.connect();

// Подписка на события
wsService.on('message', (data) => {
  console.log('Новое сообщение:', data);
});

wsService.on('chat_update', (data) => {
  console.log('Обновление чата:', data);
});
```

### Отправка сообщения через WebSocket

```javascript
wsService.sendMessage(chatId, 'Текст сообщения', userId);
```

### Обновление статуса пользователя

```javascript
wsService.updateUserStatus(userId, true); // online
wsService.updateUserStatus(userId, false); // offline
```

## 📊 Управление состоянием

### Использование Zustand Store

```javascript
import useStore from './store';

function MyComponent() {
  const {
    currentUser,
    currentChat,
    chats,
    messages,
    debugEvents,
    statistics
  } = useStore();

  // Действия
  const { setCurrentChat, addMessage, addDebugEvent } = useStore();

  return (
    <div>
      <h1>Привет, {currentUser?.first_name}!</h1>
      <p>Активных чатов: {chats.length}</p>
    </div>
  );
}
```

### Добавление пользовательских действий

```javascript
// В store/index.js
const useStore = create((set, get) => ({
  // ... существующее состояние

  // Пользовательские действия
  sendMessageWithRetry: async (chatId, text, maxRetries = 3) => {
    const { addMessage, addDebugEvent } = get();
    
    for (let i = 0; i < maxRetries; i++) {
      try {
        const response = await apiService.sendMessage(chatId, { text });
        addMessage(chatId, response.message);
        addDebugEvent({
          id: Date.now().toString(),
          timestamp: new Date().toISOString(),
          type: 'message',
          description: `Сообщение отправлено (попытка ${i + 1})`
        });
        return response;
      } catch (error) {
        addDebugEvent({
          id: Date.now().toString(),
          timestamp: new Date().toISOString(),
          type: 'error',
          description: `Ошибка отправки (попытка ${i + 1}): ${error.message}`
        });
        
        if (i === maxRetries - 1) throw error;
        await new Promise(resolve => setTimeout(resolve, 1000 * (i + 1)));
      }
    }
  }
}));
```

## 🎨 Кастомизация

### Изменение цветовой схемы

```css
/* В tailwind.config.js */
theme: {
  extend: {
    colors: {
      telegram: {
        bg: '#your-color',
        sidebar: '#your-color',
        primary: '#your-color',
        // ...
      }
    }
  }
}
```

### Добавление новых типов сообщений

```javascript
// В MessageBubble.jsx
const getMessageContent = () => {
  switch (message.type) {
    case 'text':
      return <div>{message.text}</div>;
    case 'sticker':
      return (
        <div className="flex items-center space-x-2">
          <div className="w-8 h-8 bg-telegram-secondary rounded flex items-center justify-center">
            😀
          </div>
          <span>Стикер</span>
        </div>
      );
    // Добавьте новые типы здесь
    default:
      return <div>{message.text}</div>;
  }
};
```

### Создание пользовательских компонентов

```javascript
// components/CustomMessage.jsx
import React from 'react';

const CustomMessage = ({ message, isOwn }) => {
  return (
    <div className={`custom-message ${isOwn ? 'own' : 'other'}`}>
      <div className="message-content">
        {message.text}
      </div>
      <div className="message-time">
        {formatTime(message.timestamp)}
      </div>
    </div>
  );
};

export default CustomMessage;
```

## 🧪 Тестирование

### Тестирование компонентов

```javascript
// tests/MessageBubble.test.jsx
import { render, screen } from '@testing-library/react';
import MessageBubble from '../components/MessageBubble';

test('отображает текст сообщения', () => {
  const message = {
    id: '1',
    text: 'Тестовое сообщение',
    timestamp: new Date().toISOString(),
    from: { first_name: 'Тест' }
  };

  render(<MessageBubble message={message} isOwn={false} />);
  
  expect(screen.getByText('Тестовое сообщение')).toBeInTheDocument();
});
```

### Тестирование API

```javascript
// tests/api.test.js
import apiService from '../services/api';

test('получение пользователей', async () => {
  const users = await apiService.getUsers();
  expect(Array.isArray(users.users)).toBe(true);
});
```

## 🚀 Продакшн деплой

### Сборка для продакшена

```bash
npm run build
```

### Настройка nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        root /path/to/telegram-emulator/web/dist;
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:3001;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /ws {
        proxy_pass http://localhost:3001;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

### Docker деплой

```dockerfile
# Dockerfile для frontend
FROM node:18-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 📚 Дополнительные ресурсы

- [React Documentation](https://react.dev/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Zustand](https://github.com/pmndrs/zustand)
- [Socket.io Client](https://socket.io/docs/v4/client-api/)
- [Vite](https://vitejs.dev/)
