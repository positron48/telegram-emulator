// Telegram Emulator Web Interface
class TelegramEmulator {
    constructor() {
        this.currentUser = null;
        this.currentChat = null;
        this.ws = null;
        this.chats = [];
        this.messages = {};
        this.messageCount = 0;
        this.responseTimes = [];
        
        this.init();
    }

    async init() {
        this.addDebugEvent('Инициализация приложения');
        
        // Загружаем пользователей и выбираем первого
        await this.loadUsers();
        
        // Загружаем чаты
        await this.loadChats();
        
        // Инициализируем WebSocket
        this.initWebSocket();
        
        // Настраиваем обработчики событий
        this.setupEventListeners();
        
        this.addDebugEvent('Приложение готово к работе');
    }

    async loadUsers() {
        try {
            const response = await fetch('/api/users');
            const data = await response.json();
            
            if (data.users && data.users.length > 0) {
                this.currentUser = data.users[0]; // Выбираем первого пользователя
                this.addDebugEvent(`Выбран пользователь: ${this.currentUser.username}`);
            }
        } catch (error) {
            console.error('Ошибка загрузки пользователей:', error);
            this.addDebugEvent('Ошибка загрузки пользователей');
        }
    }

    async loadChats() {
        try {
            const response = await fetch('/api/chats');
            const data = await response.json();
            
            this.chats = data.chats || [];
            this.renderChatList();
            this.addDebugEvent(`Загружено чатов: ${this.chats.length}`);
        } catch (error) {
            console.error('Ошибка загрузки чатов:', error);
            this.addDebugEvent('Ошибка загрузки чатов');
        }
    }

    renderChatList() {
        const chatList = document.getElementById('chatList');
        
        if (this.chats.length === 0) {
            chatList.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-comments"></i>
                    <h3>Нет чатов</h3>
                    <p>Создайте чат для начала общения</p>
                </div>
            `;
            return;
        }

        chatList.innerHTML = this.chats.map(chat => `
            <div class="chat-item" data-chat-id="${chat.id}" onclick="app.selectChat('${chat.id}')">
                <div class="chat-avatar">
                    ${this.getChatAvatar(chat)}
                </div>
                <div class="chat-info">
                    <div class="chat-name">${chat.title}</div>
                    <div class="chat-last-message">
                        ${chat.last_message ? chat.last_message.text : 'Нет сообщений'}
                    </div>
                </div>
                <div class="chat-meta">
                    <div class="chat-time">
                        ${chat.last_message ? this.formatTime(chat.last_message.timestamp) : ''}
                    </div>
                    ${chat.unread_count > 0 ? `<div class="unread-badge">${chat.unread_count}</div>` : ''}
                </div>
            </div>
        `).join('');
    }

    getChatAvatar(chat) {
        if (chat.type === 'private' && chat.members.length > 0) {
            const member = chat.members.find(m => m.id !== this.currentUser?.id) || chat.members[0];
            return member.first_name.charAt(0).toUpperCase();
        }
        return chat.title.charAt(0).toUpperCase();
    }

    async selectChat(chatId) {
        this.currentChat = this.chats.find(chat => chat.id === chatId);
        
        if (!this.currentChat) return;

        // Обновляем активный чат в списке
        document.querySelectorAll('.chat-item').forEach(item => {
            item.classList.remove('active');
        });
        document.querySelector(`[data-chat-id="${chatId}"]`).classList.add('active');

        // Обновляем заголовок чата
        this.updateChatHeader();

        // Загружаем сообщения
        await this.loadMessages(chatId);

        this.addDebugEvent(`Выбран чат: ${this.currentChat.title}`);
    }

    updateChatHeader() {
        const header = document.getElementById('chatHeader');
        const avatar = document.getElementById('chatHeaderAvatar');
        const name = document.getElementById('chatHeaderName');
        const status = document.getElementById('chatHeaderStatus');

        header.style.display = 'flex';
        avatar.textContent = this.getChatAvatar(this.currentChat);
        name.textContent = this.currentChat.title;
        status.textContent = this.currentChat.type === 'private' ? 'online' : `${this.currentChat.members.length} участников`;
    }

    async loadMessages(chatId) {
        try {
            const startTime = performance.now();
            const response = await fetch(`/api/chats/${chatId}/messages?limit=50`);
            const data = await response.json();
            
            const endTime = performance.now();
            const responseTime = Math.round(endTime - startTime);
            this.updateResponseTime(responseTime);
            
            this.messages[chatId] = data.messages || [];
            this.renderMessages(chatId);
            
            this.addDebugEvent(`Загружено сообщений: ${this.messages[chatId].length} (${responseTime}ms)`);
        } catch (error) {
            console.error('Ошибка загрузки сообщений:', error);
            this.addDebugEvent('Ошибка загрузки сообщений');
        }
    }

    renderMessages(chatId) {
        const container = document.getElementById('messagesContainer');
        const messages = this.messages[chatId] || [];

        if (messages.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-comment"></i>
                    <h3>Нет сообщений</h3>
                    <p>Начните общение, отправив первое сообщение</p>
                </div>
            `;
            return;
        }

        container.innerHTML = messages.map(message => `
            <div class="message ${message.from_id === this.currentUser?.id ? 'outgoing' : 'incoming'}">
                <div class="message-bubble">
                    <div class="message-text">${this.escapeHtml(message.text)}</div>
                    <div class="message-meta">
                        <span class="message-time">${this.formatTime(message.timestamp)}</span>
                        ${message.from_id === this.currentUser?.id ? 
                            `<span class="message-status ${message.status}"></span>` : ''}
                    </div>
                </div>
            </div>
        `).join('');

        // Прокручиваем к последнему сообщению
        container.scrollTop = container.scrollHeight;
    }

    async sendMessage() {
        const input = document.getElementById('messageInput');
        const text = input.value.trim();
        
        if (!text || !this.currentChat || !this.currentUser) return;

        // Очищаем поле ввода
        input.value = '';
        this.adjustTextareaHeight(input);

        try {
            const startTime = performance.now();
            const response = await fetch(`/api/chats/${this.currentChat.id}/messages`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    from_user_id: this.currentUser.id,
                    text: text,
                    type: 'text'
                })
            });

            const endTime = performance.now();
            const responseTime = Math.round(endTime - startTime);
            this.updateResponseTime(responseTime);

            if (response.ok) {
                const data = await response.json();
                this.addMessage(this.currentChat.id, data.message);
                this.messageCount++;
                this.updateMessageCount();
                this.addDebugEvent(`Сообщение отправлено (${responseTime}ms)`);
            } else {
                this.addDebugEvent('Ошибка отправки сообщения');
            }
        } catch (error) {
            console.error('Ошибка отправки сообщения:', error);
            this.addDebugEvent('Ошибка отправки сообщения');
        }
    }

    addMessage(chatId, message) {
        if (!this.messages[chatId]) {
            this.messages[chatId] = [];
        }
        
        this.messages[chatId].push(message);
        
        if (this.currentChat && this.currentChat.id === chatId) {
            this.renderMessages(chatId);
        }
    }

    initWebSocket() {
        if (!this.currentUser) return;

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws?user_id=${this.currentUser.id}`;
        
        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
            this.addDebugEvent('WebSocket подключен');
            
            // Подписываемся на события
            this.ws.send(JSON.stringify({
                type: 'subscribe',
                data: ['message', 'message_status_update', 'typing', 'chat_read']
            }));
        };

        this.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                this.handleWebSocketMessage(message);
            } catch (error) {
                console.error('Ошибка парсинга WebSocket сообщения:', error);
            }
        };

        this.ws.onclose = () => {
            this.addDebugEvent('WebSocket отключен');
            // Переподключаемся через 5 секунд
            setTimeout(() => this.initWebSocket(), 5000);
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket ошибка:', error);
            this.addDebugEvent('WebSocket ошибка');
        };
    }

    handleWebSocketMessage(message) {
        switch (message.type) {
            case 'message':
                this.handleNewMessage(message.data);
                break;
            case 'message_status_update':
                this.handleMessageStatusUpdate(message.data);
                break;
            case 'typing':
                this.handleTyping(message.data);
                break;
            case 'chat_read':
                this.handleChatRead(message.data);
                break;
            case 'subscribed':
                this.addDebugEvent('Подписка на события подтверждена');
                break;
            case 'pong':
                // Игнорируем pong сообщения
                break;
            default:
                this.addDebugEvent(`Неизвестное WebSocket сообщение: ${message.type}`);
        }
    }

    handleNewMessage(data) {
        this.addMessage(data.chat_id, {
            id: data.id,
            chat_id: data.chat_id,
            from_id: data.from.id,
            text: data.text,
            type: data.type,
            status: data.status,
            timestamp: data.timestamp
        });
        
        this.messageCount++;
        this.updateMessageCount();
        this.addDebugEvent(`Новое сообщение от ${data.from.username}`);
    }

    handleMessageStatusUpdate(data) {
        // Обновляем статус сообщения в UI
        const messageElement = document.querySelector(`[data-message-id="${data.message_id}"]`);
        if (messageElement) {
            const statusElement = messageElement.querySelector('.message-status');
            if (statusElement) {
                statusElement.className = `message-status ${data.status}`;
            }
        }
        
        this.addDebugEvent(`Статус сообщения обновлен: ${data.status}`);
    }

    handleTyping(data) {
        if (data.user_id !== this.currentUser?.id) {
            this.showTypingIndicator();
        }
    }

    handleChatRead(data) {
        // Обновляем счетчик непрочитанных сообщений
        if (this.currentChat && this.currentChat.id === data.chat_id) {
            this.loadChats(); // Перезагружаем список чатов
        }
        
        this.addDebugEvent(`Чат прочитан пользователем ${data.user_id}`);
    }

    showTypingIndicator() {
        const indicator = document.getElementById('typingIndicator');
        indicator.style.display = 'block';
        
        // Скрываем через 3 секунды
        setTimeout(() => {
            indicator.style.display = 'none';
        }, 3000);
    }

    setupEventListeners() {
        // Отправка сообщения по Enter
        const messageInput = document.getElementById('messageInput');
        messageInput.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage();
            }
        });

        // Автоматическое изменение высоты текстового поля
        messageInput.addEventListener('input', () => {
            this.adjustTextareaHeight(messageInput);
        });

        // Кнопка отправки
        const sendButton = document.getElementById('sendButton');
        sendButton.addEventListener('click', () => {
            this.sendMessage();
        });

        // Поиск чатов
        const searchInput = document.getElementById('searchInput');
        searchInput.addEventListener('input', (e) => {
            this.filterChats(e.target.value);
        });
    }

    adjustTextareaHeight(textarea) {
        textarea.style.height = 'auto';
        textarea.style.height = Math.min(textarea.scrollHeight, 100) + 'px';
    }

    filterChats(query) {
        const chatItems = document.querySelectorAll('.chat-item');
        const searchTerm = query.toLowerCase();

        chatItems.forEach(item => {
            const chatName = item.querySelector('.chat-name').textContent.toLowerCase();
            if (chatName.includes(searchTerm)) {
                item.style.display = 'flex';
            } else {
                item.style.display = 'none';
            }
        });
    }

    formatTime(timestamp) {
        const date = new Date(timestamp);
        const now = new Date();
        const diff = now - date;
        
        if (diff < 24 * 60 * 60 * 1000) {
            // Сегодня - показываем только время
            return date.toLocaleTimeString('ru-RU', { 
                hour: '2-digit', 
                minute: '2-digit' 
            });
        } else if (diff < 7 * 24 * 60 * 60 * 1000) {
            // На этой неделе - показываем день недели
            return date.toLocaleDateString('ru-RU', { weekday: 'short' });
        } else {
            // Старые сообщения - показываем дату
            return date.toLocaleDateString('ru-RU', { 
                day: '2-digit', 
                month: '2-digit' 
            });
        }
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    addDebugEvent(message) {
        const eventsContainer = document.getElementById('debugEvents');
        const now = new Date();
        const time = now.toLocaleTimeString('ru-RU');
        
        const eventElement = document.createElement('div');
        eventElement.className = 'debug-event';
        eventElement.innerHTML = `
            <div class="debug-event-time">${time}</div>
            <div class="debug-event-message">${message}</div>
        `;
        
        eventsContainer.insertBefore(eventElement, eventsContainer.firstChild);
        
        // Ограничиваем количество событий
        if (eventsContainer.children.length > 20) {
            eventsContainer.removeChild(eventsContainer.lastChild);
        }
    }

    updateMessageCount() {
        document.getElementById('messageCount').textContent = this.messageCount;
    }

    updateResponseTime(time) {
        this.responseTimes.push(time);
        
        // Оставляем только последние 10 измерений
        if (this.responseTimes.length > 10) {
            this.responseTimes.shift();
        }
        
        // Вычисляем среднее время
        const avgTime = Math.round(
            this.responseTimes.reduce((sum, t) => sum + t, 0) / this.responseTimes.length
        );
        
        document.getElementById('responseTime').textContent = `${avgTime}ms`;
    }
}

// Инициализация приложения
let app;
document.addEventListener('DOMContentLoaded', () => {
    app = new TelegramEmulator();
});
