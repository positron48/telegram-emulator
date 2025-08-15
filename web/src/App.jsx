import React, { useEffect, useState } from 'react';
import Sidebar from './components/Sidebar';
import ChatWindow from './components/ChatWindow';
import DebugPanel from './components/DebugPanel';
import useStore from './store';
import apiService from './services/api';
import wsService from './services/websocket';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

function App() {
  const {
    currentUser,
    currentChat,
    chats,
    messages,
    users,
    debugEvents,
    statistics,
    isLoading,
    error,
    isConnected,
    setCurrentUser,
    setChats,
    setUsers,
    setMessages,
    addMessage,
    updateChat,
    addDebugEvent,
    setStatistics,
    setLoading,
    setError,
    setConnected
  } = useStore();

  const [showDebugPanel, setShowDebugPanel] = useState(false);

  useEffect(() => {
    initializeApp();
  }, []);

  useEffect(() => {
    setupWebSocket();
  }, []);

  const initializeApp = async () => {
    try {
      setLoading(true);
      addDebugEvent({
        id: Date.now().toString(),
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'info',
        description: 'Инициализация приложения'
      });

      // Загружаем пользователей
      const usersResponse = await apiService.getUsers();
      setUsers(usersResponse.users || []);
      
      if (usersResponse.users && usersResponse.users.length > 0) {
        setCurrentUser(usersResponse.users[0]);
        addDebugEvent({
          id: (Date.now() + 1).toString(),
          timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
          type: 'info',
          description: `Выбран пользователь: ${usersResponse.users[0].username}`
        });
      }

      // Загружаем чаты
      const chatsResponse = await apiService.getChats();
      setChats(chatsResponse.chats || []);
      
      addDebugEvent({
        id: (Date.now() + 2).toString(),
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'info',
        description: `Загружено чатов: ${chatsResponse.chats?.length || 0}`
      });

      // Загружаем сообщения для каждого чата
      if (chatsResponse.chats) {
        for (const chat of chatsResponse.chats) {
          try {
            const messagesResponse = await apiService.getChatMessages(chat.id);
            setMessages(chat.id, messagesResponse.messages || []);
          } catch (error) {
            console.error(`Failed to load messages for chat ${chat.id}:`, error);
          }
        }
      }

      addDebugEvent({
        id: (Date.now() + 3).toString(),
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'info',
        description: 'Приложение готово к работе'
      });

    } catch (error) {
      console.error('Failed to initialize app:', error);
      setError(error.message);
      addDebugEvent({
        id: (Date.now() + 4).toString(),
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: `Ошибка инициализации: ${error.message}`
      });
    } finally {
      setLoading(false);
    }
  };

  const setupWebSocket = async () => {
    try {
      await wsService.connect();
      setConnected(true);
      
      addDebugEvent({
        id: (Date.now() + 5).toString(),
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'info',
        description: 'WebSocket соединение установлено'
      });

      // Подписываемся на события
      wsService.on('message', (data) => {
        addMessage(data.chat_id, data);
        addDebugEvent({
          id: Date.now().toString(),
          timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
          type: 'message',
          description: `Новое сообщение в чате ${data.chat_id}`
        });
      });

      wsService.on('chat_update', (data) => {
        updateChat(data.id, data);
        addDebugEvent({
          id: Date.now().toString(),
          timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
          type: 'info',
          description: `Обновление чата: ${data.title}`
        });
      });

      wsService.on('user_update', (data) => {
        addDebugEvent({
          id: Date.now().toString(),
          timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
          type: 'info',
          description: `Обновление пользователя: ${data.username}`
        });
      });

      wsService.on('debug_event', (data) => {
        addDebugEvent({
          id: Date.now().toString(),
          timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
          type: data.type,
          description: data.description,
          data: data.data
        });
      });

      wsService.on('statistics_update', (data) => {
        setStatistics(data);
      });

      wsService.on('disconnect', () => {
        setConnected(false);
        addDebugEvent({
          id: Date.now().toString(),
          timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
          type: 'error',
          description: 'WebSocket соединение разорвано'
        });
      });

    } catch (error) {
      console.error('Failed to setup WebSocket:', error);
      setConnected(false);
      addDebugEvent({
        id: Date.now().toString(),
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: `Ошибка WebSocket: ${error.message}`
      });
    }
  };

  const handleSendMessage = async (text) => {
    if (!currentChat || !currentUser || !text.trim()) return;

    try {
      const messageData = {
        text: text.trim(),
        from_user_id: currentUser.id,
        type: 'text'
      };

      const response = await apiService.sendMessage(currentChat.id, messageData);
      
      if (response.message) {
        addMessage(currentChat.id, response.message);
        
        addDebugEvent({
          id: Date.now().toString(),
          timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
          type: 'message',
          description: `Отправлено сообщение: "${text.trim()}"`
        });
      }
    } catch (error) {
      console.error('Failed to send message:', error);
      addDebugEvent({
        id: Date.now().toString(),
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: `Ошибка отправки сообщения: ${error.message}`
      });
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen bg-telegram-bg">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-telegram-primary mx-auto mb-4"></div>
          <p className="text-telegram-text">Загрузка Telegram Emulator...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-screen bg-telegram-bg">
        <div className="text-center">
          <div className="text-red-500 text-6xl mb-4">⚠️</div>
          <h1 className="text-2xl font-bold text-telegram-text mb-2">Ошибка загрузки</h1>
          <p className="text-telegram-text-secondary mb-4">{error}</p>
          <button 
            onClick={initializeApp}
            className="bg-telegram-primary text-white px-4 py-2 rounded-lg hover:bg-telegram-primary/80 transition-colors"
          >
            Попробовать снова
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-screen bg-telegram-bg">
      {/* Боковая панель */}
      <Sidebar 
        chats={chats}
        currentChat={currentChat}
        currentUser={currentUser}
        isConnected={isConnected}
        onChatSelect={(chat) => useStore.getState().setCurrentChat(chat)}
        onToggleDebug={() => setShowDebugPanel(!showDebugPanel)}
      />

      {/* Основная область чата */}
      <div className="flex-1 flex flex-col">
        <ChatWindow
          chat={currentChat}
          messages={currentChat ? messages[currentChat.id] || [] : []}
          currentUser={currentUser}
          onSendMessage={handleSendMessage}
        />
      </div>

      {/* Панель отладки */}
      {showDebugPanel && (
        <DebugPanel
          events={debugEvents}
          statistics={statistics}
          onClose={() => setShowDebugPanel(false)}
        />
      )}
    </div>
  );
}

export default App;
