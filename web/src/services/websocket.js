class WebSocketService {
  constructor() {
    this.socket = null;
    this.isConnected = false;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000;
    this.eventHandlers = new Map();
    this.reconnectTimer = null;
    this.currentUserId = null; // Сохраняем текущий userId
  }

  connect(url = 'ws://localhost:3001/ws', userId = null) {
    if (this.socket && this.isConnected) {
      return Promise.resolve();
    }

    // Сохраняем userId для переподключения
    if (userId) {
      this.currentUserId = userId;
    }

    // Добавляем user_id к URL если он передан
    const wsUrl = this.currentUserId ? `${url}?user_id=${this.currentUserId}` : url;

    return new Promise((resolve, reject) => {
      try {
        this.socket = new WebSocket(wsUrl);

        this.socket.onopen = () => {
          console.log('WebSocket connected');
          this.isConnected = true;
          this.reconnectAttempts = 0;
          this.triggerEvent('connect');
          resolve();
        };

        this.socket.onclose = (event) => {
          console.log('WebSocket disconnected:', event.code, event.reason);
          this.isConnected = false;
          this.triggerEvent('disconnect', { code: event.code, reason: event.reason });
          
          // Автоматическое переподключение
          if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
            this.reconnectTimer = setTimeout(() => {
              this.connect(url, this.currentUserId);
            }, this.reconnectDelay * this.reconnectAttempts);
          } else {
            this.triggerEvent('reconnect_failed');
          }
        };

        this.socket.onerror = (error) => {
          console.error('WebSocket connection error:', error);
          this.isConnected = false;
          this.triggerEvent('connect_error', { error });
          reject(error);
        };

        this.socket.onmessage = (event) => {
          try {
            const message = JSON.parse(event.data);
            
            // Обработка событий эмулятора
            switch (message.type) {
              case 'message':
                this.triggerEvent('message', message.data);
                break;
              case 'message_status_update':
                this.triggerEvent('message_status_update', message.data);
                break;
              case 'message_delete':
                this.triggerEvent('message_delete', message.data);
                break;
              case 'chat_read':
                this.triggerEvent('chat_read', message.data);
                break;
              case 'user_update':
                this.triggerEvent('user_update', message.data);
                break;
              case 'chat_update':
                this.triggerEvent('chat_update', message.data);
                break;
              case 'bot_update':
                this.triggerEvent('bot_update', message.data);
                break;
              case 'debug_event':
                this.triggerEvent('debug_event', message.data);
                break;
              case 'statistics_update':
                this.triggerEvent('statistics_update', message.data);
                break;
              default:
                console.log('Unknown message type:', message.type);
            }
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
          }
        };

      } catch (error) {
        console.error('Failed to create WebSocket connection:', error);
        reject(error);
      }
    });
  }

  disconnect() {
    if (this.socket) {
      if (this.reconnectTimer) {
        clearTimeout(this.reconnectTimer);
        this.reconnectTimer = null;
      }
      this.socket.close();
      this.socket = null;
      this.isConnected = false;
    }
  }

  emit(event, data) {
    if (this.socket && this.isConnected) {
      const message = {
        type: event,
        data: data
      };
      this.socket.send(JSON.stringify(message));
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
