# Telegram Emulator - Web Interface

Современный веб-интерфейс для Telegram Emulator, построенный на React 18 с использованием современных технологий.

## 🚀 Технологический стек

- **React 18** - Основной фреймворк
- **Vite** - Сборщик и dev-сервер
- **Tailwind CSS** - Стилизация
- **Zustand** - Управление состоянием
- **Socket.io-client** - WebSocket соединения
- **Lucide React** - Иконки
- **date-fns** - Работа с датами

## 📦 Установка и запуск

### Предварительные требования

- Node.js 18+
- npm или yarn

### Установка зависимостей

```bash
npm install
```

### Запуск в режиме разработки

```bash
npm run dev
```

Приложение будет доступно по адресу: http://localhost:3000

### Сборка для продакшена

```bash
npm run build
```

### Предварительный просмотр сборки

```bash
npm run preview
```

## 🏗️ Архитектура

```
src/
├── components/          # React компоненты
│   ├── Sidebar.jsx     # Боковая панель с чатами
│   ├── ChatWindow.jsx  # Окно чата
│   ├── MessageBubble.jsx # Пузырьки сообщений
│   └── DebugPanel.jsx  # Панель отладки
├── services/           # Сервисы для работы с API
│   ├── api.js         # REST API клиент
│   └── websocket.js   # WebSocket клиент
├── store/             # Управление состоянием
│   └── index.js       # Zustand store
├── types/             # Типы данных
│   └── index.js       # JSDoc типы
├── App.jsx            # Главный компонент
├── main.jsx           # Точка входа
└── index.css          # Глобальные стили
```

## 🎨 Дизайн

Интерфейс построен в стиле Telegram с использованием кастомной цветовой палитры:

- **Основной фон**: `#17212b`
- **Боковая панель**: `#242f3d`
- **Границы**: `#0e1621`
- **Акцентный цвет**: `#2b5278`
- **Вторичный текст**: `#7d8e98`

## 🔧 Конфигурация

### Vite

Конфигурация Vite настроена для проксирования API запросов к backend серверу:

```javascript
// vite.config.js
export default defineConfig({
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:3001',
        changeOrigin: true,
      },
      '/ws': {
        target: 'ws://localhost:3001',
        ws: true,
      }
    }
  }
})
```

### Tailwind CSS

Кастомная конфигурация с Telegram-подобными цветами и анимациями.

## 📱 Компоненты

### Sidebar

Боковая панель содержит:
- Поиск чатов
- Список чатов с аватарами и последними сообщениями
- Статус соединения
- Информацию о текущем пользователе
- Кнопки для отладки и настроек

### ChatWindow

Основное окно чата включает:
- Заголовок с информацией о собеседнике
- Область сообщений с автоскроллом
- Поле ввода с поддержкой Enter для отправки
- Индикатор печати

### MessageBubble

Пузырьки сообщений поддерживают:
- Разные стили для своих и чужих сообщений
- Статусы доставки (отправлено, доставлено, прочитано)
- Разные типы сообщений (текст, файлы, голосовые, фото)
- Время отправки

### DebugPanel

Панель отладки с:
- Списком событий в реальном времени
- Статистикой работы приложения
- Экспортом логов
- Системной информацией

## 🔌 API Интеграция

### REST API

Все API вызовы централизованы в `services/api.js`:

```javascript
// Примеры использования
const users = await apiService.getUsers();
const chats = await apiService.getChats();
const messages = await apiService.getChatMessages(chatId);
```

### WebSocket

WebSocket соединение управляется через `services/websocket.js`:

```javascript
// Подключение
await wsService.connect();

// Подписка на события
wsService.on('message', handleNewMessage);
wsService.on('chat_update', handleChatUpdate);
```

## 📊 Управление состоянием

Состояние приложения управляется через Zustand store (`store/index.js`):

```javascript
const {
  currentUser,
  currentChat,
  chats,
  messages,
  debugEvents,
  statistics
} = useStore();
```

## 🎯 Основные функции

- ✅ **Real-time сообщения** через WebSocket
- ✅ **Поиск чатов** с фильтрацией
- ✅ **Статусы сообщений** (отправлено, доставлено, прочитано)
- ✅ **Панель отладки** с событиями и статистикой
- ✅ **Адаптивный дизайн** для разных размеров экрана
- ✅ **Автоскролл** к новым сообщениям
- ✅ **Экспорт логов** в JSON формате
- ✅ **Индикатор печати**
- ✅ **Онлайн/оффлайн статусы**

## 🚀 Разработка

### Структура компонентов

Все компоненты следуют принципам:
- Функциональные компоненты с хуками
- Пропсы для передачи данных
- Локальное состояние для UI логики
- Переиспользуемые компоненты

### Стилизация

Используется Tailwind CSS с кастомными классами:

```css
.message-bubble.outgoing {
  @apply bg-telegram-primary text-white;
}

.message-bubble.incoming {
  @apply bg-telegram-sidebar text-telegram-text;
}
```

### Обработка ошибок

Централизованная обработка ошибок через store:

```javascript
try {
  await apiService.sendMessage(chatId, messageData);
} catch (error) {
  addDebugEvent({
    type: 'error',
    description: `Ошибка отправки: ${error.message}`
  });
}
```

## 📝 Лицензия

MIT License
