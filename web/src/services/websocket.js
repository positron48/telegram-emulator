class WebSocketService {
  constructor() {
    this.socket = null;
    this.isConnected = false;
    this.isReconnecting = false; // Новое состояние для переподключения
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000;
    this.eventHandlers = new Map();
    this.reconnectTimer = null;
    this.currentUserId = null; // Сохраняем текущий userId
  }

  connect(url = 'ws://localhost:3001/ws', userId = null) {
    // Сохраняем userId для переподключения
    if (userId) {
      this.currentUserId = userId;
    }

    // Очищаем старый socket если он есть
    if (this.socket) {
      this.socket.onclose = null; // Убираем обработчик onclose чтобы избежать рекурсии
      this.socket.onerror = null;
      this.socket.onmessage = null;
      this.socket.onopen = null;
      this.socket.close();
      this.socket = null;
    }

    // Добавляем user_id к URL если он передан
    const wsUrl = this.currentUserId ? `${url}?user_id=${this.currentUserId}` : url;

    return new Promise((resolve, reject) => {
      try {
        this.socket = new WebSocket(wsUrl);

        this.socket.onopen = () => {
          console.log('WebSocket connected');
          this.isConnected = true;
          this.isReconnecting = false; // Сбрасываем состояние переподключения
          this.reconnectAttempts = 0;
          this.triggerEvent('connect');
          resolve();
        };

        this.socket.onclose = (event) => {
          console.log('WebSocket disconnected:', event.code, event.reason);
          this.isConnected = false;
          
          // Очищаем socket
          this.socket = null;
          
          // Автоматическое переподключение только если это не было намеренное отключение
          if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            this.isReconnecting = true; // Устанавливаем состояние переподключения
            
            // Триггерим событие переподключения только для первой попытки
            if (this.reconnectAttempts === 1) {
              this.triggerEvent('reconnecting', { attempt: this.reconnectAttempts, maxAttempts: this.maxReconnectAttempts });
            }
            
            this.reconnectTimer = setTimeout(() => {
              console.log(`Starting reconnection attempt ${this.reconnectAttempts}...`);
              this.connect(url, this.currentUserId).catch(error => {
                console.error('Reconnection failed:', error);
                // Триггерим ошибку переподключения только если это не первая попытка
                if (this.reconnectAttempts > 1) {
                  this.triggerEvent('reconnect_error', { error, attempt: this.reconnectAttempts });
                }
              });
            }, this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1)); // Экспоненциальная задержка
          } else {
            // Соединение действительно потеряно
            this.isReconnecting = false;
            if (this.reconnectAttempts >= this.maxReconnectAttempts) {
              console.log('Max reconnection attempts reached');
              this.triggerEvent('reconnect_failed');
            } else {
              // Намеренное отключение
              this.triggerEvent('disconnect', { code: event.code, reason: event.reason });
            }
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
    console.log('Disconnecting WebSocket...');
    
    // Очищаем таймер переподключения
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    
    // Сбрасываем счетчик попыток и состояние переподключения
    this.reconnectAttempts = 0;
    this.isReconnecting = false;
    
    if (this.socket) {
      // Убираем обработчики чтобы избежать рекурсии
      this.socket.onclose = null;
      this.socket.onerror = null;
      this.socket.onmessage = null;
      this.socket.onopen = null;
      
      this.socket.close(1000); // Намеренное отключение
      this.socket = null;
      this.isConnected = false;
      
      console.log('WebSocket disconnected');
    }
  }

  resetReconnectAttempts() {
    console.log('Resetting reconnection attempts');
    this.reconnectAttempts = 0;
    this.isReconnecting = false;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
  }

  forceReconnect(url = 'ws://localhost:3001/ws', userId = null) {
    console.log('Force reconnecting WebSocket...');
    
    // Сбрасываем состояние
    this.resetReconnectAttempts();
    this.isConnected = false;
    this.isReconnecting = false;
    
    // Очищаем старый socket
    if (this.socket) {
      this.socket.onclose = null;
      this.socket.onerror = null;
      this.socket.onmessage = null;
      this.socket.onopen = null;
      this.socket.close();
      this.socket = null;
    }
    
    // Подключаемся заново
    return this.connect(url, userId);
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

  get reconnecting() {
    return this.isReconnecting;
  }

  get socketId() {
    return this.socket ? this.socket.id : null;
  }
}

export default new WebSocketService();
