import React, { useEffect, useState } from 'react';
import Sidebar from './components/Sidebar';
import ChatWindow from './components/ChatWindow';
import DebugPanel from './components/DebugPanel';
import CreateUserModal from './components/CreateUserModal';
import CreateChatModal from './components/CreateChatModal';
import BotManager from './components/BotManager';
import SettingsPanel from './components/SettingsPanel';
import ChatMembersModal from './components/ChatMembersModal';
import useStore from './store';
import apiService from './services/api';
import wsService from './services/websocket';
import { format } from 'date-fns';
import { ru, enUS } from 'date-fns/locale';
import { getCurrentLanguage, setLanguage, t } from './locales';

function App() {
  const getCurrentLocale = () => getCurrentLanguage() === 'en' ? enUS : ru;
  
  const {
    currentUser,
    currentChat,
    chats,
    messages,
    users,
    debugEvents,
    isLoading,
    error,
    isConnected,
    isReconnecting,
    setCurrentUser,
    setChats,
    setUsers,
    setMessages,
    addMessage,
    updateMessage,
    updateMessageStatus,
    updateChat,
    addDebugEvent,
    setLoading,
    setError,
    setConnected,
    setReconnecting
  } = useStore();

  const [showDebugPanel, setShowDebugPanel] = useState(false);
  const [showCreateUserModal, setShowCreateUserModal] = useState(false);
  const [showCreateChatModal, setShowCreateChatModal] = useState(false);
  const [showBotManager, setShowBotManager] = useState(false);
  const [showSettingsPanel, setShowSettingsPanel] = useState(false);
  const [showChatMembersModal, setShowChatMembersModal] = useState(false);

  const [isInitialized, setIsInitialized] = useState(false);
  const [isWebSocketSetup, setIsWebSocketSetup] = useState(false);

  // Единый useEffect для инициализации
  useEffect(() => {
    const init = async () => {
      if (isInitialized) return;
      
      try {
        // Инициализируем тему и язык
        initializeThemeAndLanguage();
        
        await initializeApp();
        setIsInitialized(true);
      } catch (error) {
        console.error('Failed to initialize app:', error);
      }
    };

    init();
  }, []); // Пустой массив зависимостей

  // WebSocket подключение при инициализации и смене пользователя
  useEffect(() => {
    if (!isInitialized || !currentUser) return;

    const setup = async () => {
      try {
        // Отключаем предыдущее соединение если есть
        if (wsService.connected) {
          wsService.disconnect();
        }
        
        await setupWebSocket();
        setIsWebSocketSetup(true);
      } catch (error) {
        console.error('Failed to setup WebSocket:', error);
      }
    };

    setup();
  }, [isInitialized, currentUser?.id]); // Добавляем currentUser.id в зависимости

  // Обработчики WebSocket событий
  useEffect(() => {
    // Регистрируем обработчики независимо от состояния подключения
    // чтобы они работали при переподключении

    // Подписываемся на события
    const handleMessage = (data) => {
      console.log('🔍 Received message data:', {
        id: data.id,
        text: data.text,
        from: data.from,
        from_id: data.from_id,
        chat_id: data.chat_id,
        status: data.status
      });
      // Проверяем, является ли сообщение от текущего пользователя
      const isOwnMessage = data.from?.id === currentUser?.id;
      
      // Проверяем, есть ли уже временное сообщение с таким же текстом и отправителем
      const currentState = useStore.getState();
      const existingMessages = currentState.messages[data.chat_id] || [];
      
      // Ищем временное сообщение по тексту и отправителю
      let tempMessageIndex = existingMessages.findIndex(msg => 
        msg.id.startsWith('temp-') && 
        msg.text === data.text && 
        (msg.from?.id === data.from?.id || msg.from_id === data.from?.id)
      );
      
      // Если не найдено, попробуем найти по тексту среди сообщений от текущего пользователя
      if (tempMessageIndex === -1 && isOwnMessage) {
        tempMessageIndex = existingMessages.findIndex(msg => 
          msg.id.startsWith('temp-') && 
          msg.text === data.text
        );
      }
      
      // Если все еще не найдено, попробуем найти самое последнее временное сообщение от текущего пользователя
      if (tempMessageIndex === -1 && isOwnMessage) {
        const tempMessages = existingMessages.filter(msg => 
          msg.id.startsWith('temp-') && 
          (msg.from?.id === currentUser?.id || msg.from_id === currentUser?.id)
        );
        if (tempMessages.length > 0) {
          // Берем самое последнее временное сообщение
          const lastTempMessage = tempMessages[tempMessages.length - 1];
          tempMessageIndex = existingMessages.findIndex(msg => msg.id === lastTempMessage.id);
        }
      }
      
      console.log('🔍 Looking for temp message:', {
        searchText: data.text,
        searchFromId: data.from?.id,
        tempMessages: existingMessages.filter(m => m.id.startsWith('temp-')).map(m => ({
          id: m.id,
          text: m.text,
          fromId: m.from?.id || m.from_id
        })),
        foundIndex: tempMessageIndex
      });

      if (tempMessageIndex !== -1) {
        // Заменяем временное сообщение на реальное
        const updatedMessages = [...existingMessages];
        updatedMessages[tempMessageIndex] = {
          ...data,
          is_outgoing: true, // Помечаем как исходящее
          status: data.status || 'sent' // Устанавливаем статус из полученного сообщения
        };
        setMessages(data.chat_id, updatedMessages);
        
        console.log(`✅ Temporary message replaced: ${existingMessages[tempMessageIndex].id} -> ${data.id} with status: ${data.status}`);
      } else if (!isOwnMessage) {
        // Добавляем новое сообщение только если оно не от текущего пользователя
        addMessage(data.chat_id, {
          ...data,
          is_outgoing: false // Помечаем как входящее
        });
        
        addDebugEvent({
          id: `message-new-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
          timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
          type: 'message',
          description: `Новое сообщение от ${data.from?.username} в чате ${data.chat_id}`
        });
      } else if (isOwnMessage) {
        // Логируем случаи, когда собственное сообщение не заменило временное
        console.log(`❌ Own message received but no temp message found:`, {
          messageId: data.id,
          text: data.text,
          fromId: data.from?.id,
          currentUserId: currentUser?.id,
          existingTempMessages: existingMessages.filter(m => m.id.startsWith('temp-')).map(m => ({
            id: m.id,
            text: m.text,
            fromId: m.from?.id || m.from_id
          }))
        });
      }
      // Игнорируем собственные сообщения без логирования - это нормальное поведение
    };

    const handleChatUpdate = (data) => {
      updateChat(data.id, data);
      // Не логируем обновления чатов - это нормальное поведение
    };

    const handleUserUpdate = (data) => {
      // Не логируем обновления пользователей - это нормальное поведение
    };

    const handleDebugEvent = (data) => {
      addDebugEvent({
        id: `debug-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: data.type,
        description: data.description,
        data: data.data
      });
    };

    const handleMessageStatusUpdate = (data) => {
      updateMessageStatus(data.message_id, data.status);
      // Логируем для отладки проблем со статусами
      console.log(`Status update: ${data.message_id} -> ${data.status}`);
    };



    const handleDisconnect = () => {
      setConnected(false);
      setReconnecting(false);
      addDebugEvent({
        id: `disconnect-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: 'WebSocket соединение разорвано'
      });
    };

    const handleReconnecting = (data) => {
      setReconnecting(true);
      // Не логируем стандартное переподключение в отладку - это нормальное поведение
      console.log(`WebSocket reconnecting (${data.attempt}/${data.maxAttempts})`);
    };

    const handleReconnectError = (data) => {
      // Логируем только неудачные попытки переподключения (не первую)
      addDebugEvent({
        id: `reconnect-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: `${t('reconnectError', getCurrentLanguage())} (${t('attempt', getCurrentLanguage())} ${data.attempt}): ${data.error.message || t('unknownError', getCurrentLanguage())}`
      });
    };

    const handleReconnectFailed = () => {
      setReconnecting(false);
      addDebugEvent({
        id: `reconnect-failed-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: t('maxReconnectAttemptsExceeded', getCurrentLanguage())
      });
    };

    const handleConnect = () => {
      setConnected(true);
      setReconnecting(false);
      addDebugEvent({
        id: `connect-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'info',
        description: t('websocketConnected', getCurrentLanguage())
      });
    };

    const handleConnectError = (data) => {
      setConnected(false);
      setReconnecting(false);
      addDebugEvent({
        id: `connect-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: `${t('websocketConnectionError', getCurrentLanguage())}: ${data.error?.message || t('unknownError', getCurrentLanguage())}`
      });
    };

    // Регистрируем обработчики
    wsService.on('connect', handleConnect);
    wsService.on('connect_error', handleConnectError);
    wsService.on('message', handleMessage);
    wsService.on('message_status_update', handleMessageStatusUpdate);
    wsService.on('chat_update', handleChatUpdate);
    wsService.on('user_update', handleUserUpdate);
    wsService.on('debug_event', handleDebugEvent);

    wsService.on('disconnect', handleDisconnect);
    wsService.on('reconnecting', handleReconnecting);
    wsService.on('reconnect_error', handleReconnectError);
    wsService.on('reconnect_failed', handleReconnectFailed);

    // Очистка обработчиков
    return () => {
      wsService.off('connect', handleConnect);
      wsService.off('connect_error', handleConnectError);
      wsService.off('message', handleMessage);
      wsService.off('message_status_update', handleMessageStatusUpdate);
      wsService.off('chat_update', handleChatUpdate);
      wsService.off('user_update', handleUserUpdate);
      wsService.off('debug_event', handleDebugEvent);

      wsService.off('disconnect', handleDisconnect);
      wsService.off('reconnecting', handleReconnecting);
      wsService.off('reconnect_error', handleReconnectError);
      wsService.off('reconnect_failed', handleReconnectFailed);
    };
      }, [currentUser?.id]); // Убираем isConnected из зависимостей, чтобы обработчики работали при переподключении

  const initializeThemeAndLanguage = () => {
    try {
      // Загружаем настройки из localStorage
      const savedSettings = localStorage.getItem('telegram-emulator-settings');
      if (savedSettings) {
        const settings = JSON.parse(savedSettings);
        
        // Применяем тему
        if (settings.theme === 'dark') {
          document.documentElement.classList.add('dark');
        } else {
          document.documentElement.classList.remove('dark');
        }
        
        // Применяем язык
        if (settings.language) {
          setLanguage(settings.language);
        }
      }
    } catch (error) {
      console.error('Failed to initialize theme and language:', error);
    }
  };

  const initializeApp = async () => {
    try {
      setLoading(true);
      addDebugEvent({
        id: `init-start-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'info',
        description: t('appInitialization', getCurrentLanguage())
      });

      // Загружаем пользователей
      const usersResponse = await apiService.getUsers();
      setUsers(usersResponse.users || []);
      
      if (usersResponse.users && usersResponse.users.length > 0) {
        setCurrentUser(usersResponse.users[0]);
        addDebugEvent({
          id: `user-select-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
          timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
          type: 'info',
          description: `Выбран пользователь: ${usersResponse.users[0].username}`
        });
      }

      // Загружаем чаты
      const chatsResponse = await apiService.getChats(usersResponse.users[0]?.id);
      setChats(chatsResponse.chats || []);
      
      addDebugEvent({
        id: `chats-loaded-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'info',
        description: `Загружено чатов: ${chatsResponse.chats?.length || 0}`
      });

      // Загружаем сообщения для каждого чата
      if (chatsResponse.chats) {
        for (const chat of chatsResponse.chats) {
          try {
            const messagesResponse = await apiService.getChatMessages(chat.id);
            const sortedMessages = (messagesResponse.messages || []).sort((a, b) => 
              new Date(a.timestamp) - new Date(b.timestamp)
            );
            setMessages(chat.id, sortedMessages);
          } catch (error) {
            console.error(`Failed to load messages for chat ${chat.id}:`, error);
          }
        }
      }

      addDebugEvent({
        id: `app-ready-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'info',
        description: t('appReady', getCurrentLanguage())
      });

    } catch (error) {
      console.error('Failed to initialize app:', error);
      setError(error.message);
      addDebugEvent({
        id: `init-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: `${t('initializationError', getCurrentLanguage())}: ${error.message}`
      });
    } finally {
      setLoading(false);
    }
  };

  const setupWebSocket = async () => {
    try {
      // Подключаемся к WebSocket с user_id текущего пользователя
      const userId = currentUser?.id || 'anonymous';
      await wsService.connect('ws://localhost:3001/ws', userId);
      setConnected(true);
      
      addDebugEvent({
        id: `ws-connect-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'info',
        description: `WebSocket соединение установлено для пользователя ${currentUser?.username || 'anonymous'}`
      });

    } catch (error) {
      console.error('Failed to setup WebSocket:', error);
      setConnected(false);
      addDebugEvent({
        id: `ws-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: `Ошибка WebSocket: ${error.message}`
      });
    }
  };

  const handleSendMessage = async (text) => {
    if (!currentChat || !currentUser || !text.trim()) return;

    try {
      // Создаем временное сообщение для оптимистичного обновления UI
      const tempMessage = {
        id: `temp-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        chat_id: currentChat.id,
        from: currentUser,
        from_id: currentUser.id,
        text: text.trim(),
        type: 'text',
        status: 'sending',
        timestamp: new Date(),
        is_outgoing: true
      };

      // Добавляем сообщение в локальное состояние сразу
      addMessage(currentChat.id, tempMessage);
      
      // Очищаем старые временные сообщения (оставляем только последние 5)
      const currentState = useStore.getState();
      const currentMessages = currentState.messages[currentChat.id] || [];
      const tempMessages = currentMessages.filter(msg => msg.id.startsWith('temp-'));
      if (tempMessages.length > 5) {
        const messagesToRemove = tempMessages.slice(0, tempMessages.length - 5);
        const updatedMessages = currentMessages.filter(msg => !messagesToRemove.some(rm => rm.id === msg.id));
        setMessages(currentChat.id, updatedMessages);
        console.log(`🧹 Cleaned up ${messagesToRemove.length} old temp messages`);
      }

      // Fallback: если через 3 секунды сообщение не заменилось, обновляем статус вручную
      setTimeout(() => {
        const currentState = useStore.getState();
        const currentMessages = currentState.messages[currentChat.id] || [];
        const tempMsgIndex = currentMessages.findIndex(msg => msg.id === tempMessage.id);
        
        if (tempMsgIndex !== -1 && currentMessages[tempMsgIndex].status === 'sending') {
          const updatedMessages = [...currentMessages];
          updatedMessages[tempMsgIndex] = { ...currentMessages[tempMsgIndex], status: 'sent' };
          setMessages(currentChat.id, updatedMessages);
          
          console.log(`⚠️ Fallback triggered: status updated to 'sent' for temp message ${tempMessage.id}`);
        }
      }, 3000);

      // Отправляем сообщение через WebSocket
      if (wsService.connected) {
        wsService.sendMessage(currentChat.id, text.trim(), currentUser.id);
      } else {
        addDebugEvent({
          id: `ws-not-connected-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
          timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
          type: 'error',
          description: 'WebSocket не подключен, сообщение не отправлено'
        });
      }
      
    } catch (error) {
      console.error('Failed to send message:', error);
      addDebugEvent({
        id: `send-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: `Ошибка отправки сообщения: ${error.message}`
      });
    }
  };

  const handleDeleteUser = async (userId) => {
    try {
      await apiService.deleteUser(userId);
      // Обновляем список пользователей
      const usersResponse = await apiService.getUsers();
      setUsers(usersResponse.users || []);
      
      // Если удаленный пользователь был текущим, выбираем первого доступного
      if (currentUser?.id === userId) {
        const remainingUsers = usersResponse.users.filter(u => u.id !== userId);
        if (remainingUsers.length > 0) {
          setCurrentUser(remainingUsers[0]);
          const chatsResponse = await apiService.getChats(remainingUsers[0].id);
          setChats(chatsResponse.chats || []);
        } else {
          setCurrentUser(null);
          setChats([]);
        }
      }
      
      addDebugEvent({
        id: `user-deleted-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'warning',
        description: t('userDeleted', getCurrentLanguage())
      });
    } catch (error) {
      console.error('Failed to delete user:', error);
      addDebugEvent({
        id: `user-delete-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: `Ошибка удаления пользователя: ${error.message}`
      });
    }
  };

  const handleDeleteChat = async (chatId) => {
    try {
      await apiService.deleteChat(chatId);
      
      // Обновляем список чатов
      const chatsResponse = await apiService.getChats(currentUser?.id);
      setChats(chatsResponse.chats || []);
      
      // Если удаленный чат был текущим, очищаем выбор
      if (currentChat?.id === chatId) {
        useStore.getState().setCurrentChat(null);
      }
      
      addDebugEvent({
        id: `chat-deleted-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'warning',
        description: t('chatDeleted', getCurrentLanguage())
      });
    } catch (error) {
      console.error('Failed to delete chat:', error);
      addDebugEvent({
        id: `chat-delete-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: `Ошибка удаления чата: ${error.message}`
      });
    }
  };

  const handleReconnect = async () => {
    try {
      addDebugEvent({
        id: `manual-reconnect-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'info',
        description: t('manualWebSocketReconnect', getCurrentLanguage())
      });

      await wsService.forceReconnect('ws://localhost:3001/ws', currentUser?.id);
      
      addDebugEvent({
        id: `manual-reconnect-success-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'info',
        description: 'WebSocket успешно переподключен'
      });
    } catch (error) {
      console.error('Manual reconnection failed:', error);
      addDebugEvent({
        id: `manual-reconnect-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: ru }),
        type: 'error',
        description: `Ошибка ручного переподключения: ${error.message}`
      });
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen bg-telegram-bg">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-telegram-primary mx-auto mb-4"></div>
          <p className="text-telegram-text">{t('loadingTelegramEmulator', getCurrentLanguage())}</p>
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
        users={users}
        isConnected={isConnected}
        isReconnecting={isReconnecting}
        onChatSelect={(chat) => useStore.getState().setCurrentChat(chat)}
        onToggleDebug={() => setShowDebugPanel(!showDebugPanel)}
        onUserSelect={(user) => setCurrentUser(user)}
        onCreateUser={() => setShowCreateUserModal(true)}
        onDeleteUser={handleDeleteUser}
        onCreateChat={() => setShowCreateChatModal(true)}
        onDeleteChat={handleDeleteChat}
        onShowBotManager={() => setShowBotManager(true)}
        onShowSettings={() => setShowSettingsPanel(true)}
        onReconnect={handleReconnect}
      />

      {/* Основная область чата */}
      <div className="flex-1 flex flex-col">
        <ChatWindow
          chat={currentChat}
          messages={currentChat ? messages[currentChat.id] || [] : []}
          currentUser={currentUser}
          onSendMessage={handleSendMessage}
          onShowMembers={() => setShowChatMembersModal(true)}
        />
      </div>

      {/* Панель отладки */}
      {showDebugPanel && (
        <DebugPanel
          events={debugEvents}
          onClose={() => setShowDebugPanel(false)}
        />
      )}

      {/* Модальные окна */}
      <CreateUserModal
        isOpen={showCreateUserModal}
        onClose={() => setShowCreateUserModal(false)}
        onUserCreated={(user) => {
          setCurrentUser(user);
          setShowCreateUserModal(false);
        }}
      />

      <CreateChatModal
        isOpen={showCreateChatModal}
        onClose={() => setShowCreateChatModal(false)}
        onChatCreated={(chat) => {
          useStore.getState().setCurrentChat(chat);
          setShowCreateChatModal(false);
        }}
      />

      <BotManager
        isOpen={showBotManager}
        onClose={() => setShowBotManager(false)}
      />

      <SettingsPanel
        isOpen={showSettingsPanel}
        onClose={() => setShowSettingsPanel(false)}
      />

      <ChatMembersModal
        isOpen={showChatMembersModal}
        onClose={() => setShowChatMembersModal(false)}
        chat={currentChat}
      />
    </div>
  );
}

export default App;
