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

  // Single useEffect for initialization
  useEffect(() => {
    const init = async () => {
      if (isInitialized) return;
      
      try {
        // Initialize theme and language
        initializeThemeAndLanguage();
        
        await initializeApp();
        setIsInitialized(true);
      } catch (error) {
        console.error('Failed to initialize app:', error);
      }
    };

    init();
  }, []); // Empty dependency array

  // WebSocket connection on initialization and user change
  useEffect(() => {
    if (!isInitialized || !currentUser) return;

    const setup = async () => {
      try {
        // Disconnect previous connection if exists
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
  }, [isInitialized, currentUser?.id]); // Add currentUser.id to dependencies

  // WebSocket event handlers
  useEffect(() => {
    // Register handlers regardless of connection state
    // so they work on reconnection

    // Subscribe to events
    const handleMessage = (data) => {
      console.log('🔍 Received message data:', {
        id: data.id,
        text: data.text,
        from: data.from,
        from_id: data.from_id,
        chat_id: data.chat_id,
        status: data.status
      });
      // Check if message is from current user
      const isOwnMessage = data.from?.id === currentUser?.id;
      
      // Check if there's already a temporary message with the same text and sender
      const currentState = useStore.getState();
      const existingMessages = currentState.messages[data.chat_id] || [];
      
      // Look for temporary message by text and sender
      let tempMessageIndex = existingMessages.findIndex(msg => 
        msg.id.startsWith('temp-') && 
        msg.text === data.text && 
        (msg.from?.id === data.from?.id || msg.from_id === data.from?.id)
      );
      
      // If not found, try to find by text among messages from current user
      if (tempMessageIndex === -1 && isOwnMessage) {
        tempMessageIndex = existingMessages.findIndex(msg => 
          msg.id.startsWith('temp-') && 
          msg.text === data.text
        );
      }
      
      // If still not found, try to find the latest temporary message from current user
      if (tempMessageIndex === -1 && isOwnMessage) {
        const tempMessages = existingMessages.filter(msg => 
          msg.id.startsWith('temp-') && 
          (msg.from?.id === currentUser?.id || msg.from_id === currentUser?.id)
        );
        if (tempMessages.length > 0) {
          // Take the latest temporary message
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
        // Replace temporary message with real one
        const updatedMessages = [...existingMessages];
        updatedMessages[tempMessageIndex] = {
          ...data,
          is_outgoing: true, // Mark as outgoing
          status: data.status || 'sent' // Set status from received message
        };
        setMessages(data.chat_id, updatedMessages);
        
        console.log(`✅ Temporary message replaced: ${existingMessages[tempMessageIndex].id} -> ${data.id} with status: ${data.status}`);
      } else if (!isOwnMessage) {
        // Add new message only if it's not from current user
        addMessage(data.chat_id, {
          ...data,
          is_outgoing: false // Mark as incoming
        });
        
        addDebugEvent({
          id: `message-new-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
          timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
          type: 'message',
          description: `${t('newMessageFrom', getCurrentLanguage())} ${data.from?.username} ${t('inChat', getCurrentLanguage())} ${data.chat_id}`
        });
      } else if (isOwnMessage) {
        // Log cases when own message didn't replace temporary one
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
      // Ignore own messages without logging - this is normal behavior
    };

    const handleChatUpdate = (data) => {
      updateChat(data.id, data);
      // Don't log chat updates - this is normal behavior
    };

    const handleUserUpdate = (data) => {
      // Don't log user updates - this is normal behavior
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
      // Log for debugging status issues
      console.log(`Status update: ${data.message_id} -> ${data.status}`);
    };



    const handleDisconnect = () => {
      setConnected(false);
      setReconnecting(false);
      addDebugEvent({
        id: `disconnect-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'error',
        description: t('websocketConnectionError', getCurrentLanguage())
      });
    };

    const handleReconnecting = (data) => {
      setReconnecting(true);
      // Don't log standard reconnection to debug - this is normal behavior
      console.log(`WebSocket reconnecting (${data.attempt}/${data.maxAttempts})`);
    };

    const handleReconnectError = (data) => {
      // Log only failed reconnection attempts (not the first one)
      addDebugEvent({
        id: `reconnect-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'error',
        description: `${t('reconnectError', getCurrentLanguage())} (${t('attempt', getCurrentLanguage())} ${data.attempt}): ${data.error.message || t('unknownError', getCurrentLanguage())}`
      });
    };

    const handleReconnectFailed = () => {
      setReconnecting(false);
      addDebugEvent({
        id: `reconnect-failed-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'error',
        description: t('maxReconnectAttemptsExceeded', getCurrentLanguage())
      });
    };

    const handleConnect = () => {
      setConnected(true);
      setReconnecting(false);
      addDebugEvent({
        id: `connect-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'info',
        description: t('websocketConnected', getCurrentLanguage())
      });
    };

    const handleConnectError = (data) => {
      setConnected(false);
      setReconnecting(false);
      addDebugEvent({
        id: `connect-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'error',
        description: `${t('websocketConnectionError', getCurrentLanguage())}: ${data.error?.message || t('unknownError', getCurrentLanguage())}`
      });
    };

    // Register handlers
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

    // Cleanup handlers
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
      }, [currentUser?.id]); // Remove isConnected from dependencies so handlers work on reconnection

  const initializeThemeAndLanguage = () => {
    try {
      // Load settings from localStorage
      const savedSettings = localStorage.getItem('telegram-emulator-settings');
      if (savedSettings) {
        const settings = JSON.parse(savedSettings);
        
        // Apply theme
        if (settings.theme === 'dark') {
          document.documentElement.classList.add('dark');
        } else {
          document.documentElement.classList.remove('dark');
        }
        
        // Apply language
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
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'info',
        description: t('appInitialization', getCurrentLanguage())
      });

      // Load users
      const usersResponse = await apiService.getUsers();
      setUsers(usersResponse.users || []);
      
      if (usersResponse.users && usersResponse.users.length > 0) {
        setCurrentUser(usersResponse.users[0]);
        addDebugEvent({
          id: `user-select-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
          timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
          type: 'info',
          description: `${t('user', getCurrentLanguage())} selected: ${usersResponse.users[0].username}`
        });
      }

      // Load chats
      const chatsResponse = await apiService.getChats(usersResponse.users[0]?.id);
      setChats(chatsResponse.chats || []);
      
      addDebugEvent({
        id: `chats-loaded-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'info',
        description: `${t('chatsLoaded', getCurrentLanguage())}: ${chatsResponse.chats?.length || 0}`
      });

      // Load messages for each chat
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
      // Connect to WebSocket with current user's user_id
      const userId = currentUser?.id || 'anonymous';
      await wsService.connect('ws://localhost:3001/ws', userId);
      setConnected(true);
      
      addDebugEvent({
        id: `ws-connect-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'info',
        description: `${t('websocketConnected', getCurrentLanguage())} for user ${currentUser?.username || 'anonymous'}`
      });

    } catch (error) {
      console.error('Failed to setup WebSocket:', error);
      setConnected(false);
      addDebugEvent({
        id: `ws-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'error',
        description: `${t('websocketConnectionError', getCurrentLanguage())}: ${error.message}`
      });
    }
  };

  const handleCallbackQuery = async (button) => {
    try {
      console.log('Callback query:', button);
      
      // Отправляем callback query через WebSocket
      if (wsService.connected && currentChat) {
        wsService.sendCallbackQuery(button, currentChat.id);
      }
      
      // Логируем событие
      addDebugEvent({
        id: `callback-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'info',
        description: `Callback query: ${button.text} (${button.callback_data || button.url || 'no data'})`
      });
    } catch (error) {
      console.error('Failed to handle callback query:', error);
      addDebugEvent({
        id: `callback-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'error',
        description: `Callback query error: ${error.message}`
      });
    }
  };

  const handleSendMessage = async (text) => {
    if (!currentChat || !currentUser || !text.trim()) return;

    try {
      // Create temporary message for optimistic UI update
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

      // Add message to local state immediately
      addMessage(currentChat.id, tempMessage);
      
      // Clean up old temporary messages (keep only last 5)
      const currentState = useStore.getState();
      const currentMessages = currentState.messages[currentChat.id] || [];
      const tempMessages = currentMessages.filter(msg => msg.id.startsWith('temp-'));
      if (tempMessages.length > 5) {
        const messagesToRemove = tempMessages.slice(0, tempMessages.length - 5);
        const updatedMessages = currentMessages.filter(msg => !messagesToRemove.some(rm => rm.id === msg.id));
        setMessages(currentChat.id, updatedMessages);
        console.log(`🧹 Cleaned up ${messagesToRemove.length} old temp messages`);
      }

      // Fallback: if after 3 seconds message wasn't replaced, update status manually
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

      // Send message via WebSocket
      if (wsService.connected) {
        wsService.sendMessage(currentChat.id, text.trim(), currentUser.id);
      } else {
        addDebugEvent({
          id: `ws-not-connected-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
          timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
          type: 'error',
          description: t('websocketConnectionError', getCurrentLanguage())
        });
      }
      
    } catch (error) {
      console.error('Failed to send message:', error);
      addDebugEvent({
        id: `send-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'error',
        description: `${t('failedToSave', getCurrentLanguage())}: ${error.message}`
      });
    }
  };

  const handleDeleteUser = async (userId) => {
    try {
      await apiService.deleteUser(userId);
      // Update users list
      const usersResponse = await apiService.getUsers();
      setUsers(usersResponse.users || []);
      
      // If deleted user was current, select first available
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
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'warning',
        description: t('userDeleted', getCurrentLanguage())
      });
    } catch (error) {
      console.error('Failed to delete user:', error);
      addDebugEvent({
        id: `user-delete-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'error',
        description: `${t('failedToDelete', getCurrentLanguage())}: ${error.message}`
      });
    }
  };

  const handleDeleteChat = async (chatId) => {
    try {
      await apiService.deleteChat(chatId);
      
      // Update chats list
      const chatsResponse = await apiService.getChats(currentUser?.id);
      setChats(chatsResponse.chats || []);
      
      // If deleted chat was current, clear selection
      if (currentChat?.id === chatId) {
        useStore.getState().setCurrentChat(null);
      }
      
      addDebugEvent({
        id: `chat-deleted-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'warning',
        description: t('chatDeleted', getCurrentLanguage())
      });
    } catch (error) {
      console.error('Failed to delete chat:', error);
      addDebugEvent({
        id: `chat-delete-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'error',
        description: `${t('failedToDelete', getCurrentLanguage())}: ${error.message}`
      });
    }
  };

  const handleReconnect = async () => {
    try {
      addDebugEvent({
        id: `manual-reconnect-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'info',
        description: t('manualWebSocketReconnect', getCurrentLanguage())
      });

      await wsService.forceReconnect('ws://localhost:3001/ws', currentUser?.id);
      
      addDebugEvent({
        id: `manual-reconnect-success-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'info',
        description: t('websocketConnected', getCurrentLanguage())
      });
    } catch (error) {
      console.error('Manual reconnection failed:', error);
      addDebugEvent({
        id: `manual-reconnect-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: format(new Date(), 'HH:mm:ss', { locale: getCurrentLanguage() === 'ru' ? ru : enUS }),
        type: 'error',
        description: `${t('reconnectError', getCurrentLanguage())}: ${error.message}`
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
          <h1 className="text-2xl font-bold text-telegram-text mb-2">{t('loadingError', getCurrentLanguage())}</h1>
          <p className="text-telegram-text-secondary mb-4">{error}</p>
          <button 
            onClick={initializeApp}
            className="bg-telegram-primary text-white px-4 py-2 rounded-lg hover:bg-telegram-primary/80 transition-colors"
          >
            {t('tryAgain', getCurrentLanguage())}
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-screen bg-telegram-bg">
      {/* Sidebar */}
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

      {/* Main chat area */}
      <div className="flex-1 flex flex-col">
        <ChatWindow
          chat={currentChat}
          messages={currentChat ? messages[currentChat.id] || [] : []}
          currentUser={currentUser}
          onSendMessage={handleSendMessage}
          onShowMembers={() => setShowChatMembersModal(true)}
          onCallbackQuery={handleCallbackQuery}
        />
      </div>

      {/* Debug panel */}
      {showDebugPanel && (
        <DebugPanel
          events={debugEvents}
          onClose={() => setShowDebugPanel(false)}
        />
      )}

      {/* Modal windows */}
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
