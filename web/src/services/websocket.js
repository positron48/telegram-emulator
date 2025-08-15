import { io } from 'socket.io-client';

class WebSocketService {
  constructor() {
    this.socket = null;
    this.isConnected = false;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000;
    this.eventHandlers = new Map();
  }

  connect(url = 'ws://localhost:3001/ws') {
    if (this.socket && this.isConnected) {
      return Promise.resolve();
    }

    return new Promise((resolve, reject) => {
      try {
        this.socket = io(url, {
          transports: ['websocket'],
          timeout: 5000,
          reconnection: true,
          reconnectionAttempts: this.maxReconnectAttempts,
          reconnectionDelay: this.reconnectDelay,
        });

        this.socket.on('connect', () => {
          console.log('WebSocket connected');
          this.isConnected = true;
          this.reconnectAttempts = 0;
          this.emit('subscribe', {
            events: ['message', 'user_update', 'chat_update', 'bot_update']
          });
          resolve();
        });

        this.socket.on('disconnect', (reason) => {
          console.log('WebSocket disconnected:', reason);
          this.isConnected = false;
          this.triggerEvent('disconnect', { reason });
        });

        this.socket.on('connect_error', (error) => {
          console.error('WebSocket connection error:', error);
          this.isConnected = false;
          this.triggerEvent('connect_error', { error });
          reject(error);
        });

        this.socket.on('reconnect', (attemptNumber) => {
          console.log('WebSocket reconnected after', attemptNumber, 'attempts');
          this.isConnected = true;
          this.triggerEvent('reconnect', { attemptNumber });
        });

        this.socket.on('reconnect_error', (error) => {
          console.error('WebSocket reconnection error:', error);
          this.triggerEvent('reconnect_error', { error });
        });

        this.socket.on('reconnect_failed', () => {
          console.error('WebSocket reconnection failed');
          this.triggerEvent('reconnect_failed');
        });

        // Обработка событий эмулятора
        this.socket.on('message', (data) => {
          this.triggerEvent('message', data);
        });

        this.socket.on('user_update', (data) => {
          this.triggerEvent('user_update', data);
        });

        this.socket.on('chat_update', (data) => {
          this.triggerEvent('chat_update', data);
        });

        this.socket.on('bot_update', (data) => {
          this.triggerEvent('bot_update', data);
        });

        this.socket.on('debug_event', (data) => {
          this.triggerEvent('debug_event', data);
        });

        this.socket.on('statistics_update', (data) => {
          this.triggerEvent('statistics_update', data);
        });

      } catch (error) {
        console.error('Failed to create WebSocket connection:', error);
        reject(error);
      }
    });
  }

  disconnect() {
    if (this.socket) {
      this.socket.disconnect();
      this.socket = null;
      this.isConnected = false;
    }
  }

  emit(event, data) {
    if (this.socket && this.isConnected) {
      this.socket.emit(event, data);
    } else {
      console.warn('WebSocket not connected, cannot emit event:', event);
    }
  }

  on(event, handler) {
    if (!this.eventHandlers.has(event)) {
      this.eventHandlers.set(event, []);
    }
    this.eventHandlers.get(event).push(handler);
  }

  off(event, handler) {
    if (this.eventHandlers.has(event)) {
      const handlers = this.eventHandlers.get(event);
      const index = handlers.indexOf(handler);
      if (index > -1) {
        handlers.splice(index, 1);
      }
    }
  }

  triggerEvent(event, data) {
    if (this.eventHandlers.has(event)) {
      this.eventHandlers.get(event).forEach(handler => {
        try {
          handler(data);
        } catch (error) {
          console.error(`Error in event handler for ${event}:`, error);
        }
      });
    }
  }

  // Специфичные методы для эмулятора
  subscribeToEvents(events) {
    this.emit('subscribe', { events });
  }

  unsubscribeFromEvents(events) {
    this.emit('unsubscribe', { events });
  }

  sendMessage(chatId, text, fromUserId) {
    this.emit('send_message', {
      chat_id: chatId,
      text,
      from_user_id: fromUserId
    });
  }

  markMessageAsRead(messageId) {
    this.emit('mark_message_read', { message_id: messageId });
  }

  updateUserStatus(userId, isOnline) {
    this.emit('update_user_status', {
      user_id: userId,
      is_online: isOnline
    });
  }

  // Геттеры
  get connected() {
    return this.isConnected;
  }

  get socketId() {
    return this.socket ? this.socket.id : null;
  }
}

export default new WebSocketService();
