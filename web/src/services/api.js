const API_BASE_URL = '/api';

class ApiService {
  constructor() {
    this.baseURL = API_BASE_URL;
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      return await response.json();
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Users API
  async getUsers() {
    return this.request('/users');
  }

  async createUser(userData) {
    return this.request('/users', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  }

  async getUserById(userId) {
    return this.request(`/users/${userId}`);
  }

  async updateUser(userId, userData) {
    return this.request(`/users/${userId}`, {
      method: 'PUT',
      body: JSON.stringify(userData),
    });
  }

  async deleteUser(userId) {
    return this.request(`/users/${userId}`, {
      method: 'DELETE',
    });
  }

  async getUserChats(userId) {
    return this.request(`/users/${userId}/chats`);
  }

  // Chats API
  async getChats() {
    return this.request('/chats');
  }

  async createChat(chatData) {
    return this.request('/chats', {
      method: 'POST',
      body: JSON.stringify(chatData),
    });
  }

  async getChatById(chatId) {
    return this.request(`/chats/${chatId}`);
  }

  async updateChat(chatId, chatData) {
    return this.request(`/chats/${chatId}`, {
      method: 'PUT',
      body: JSON.stringify(chatData),
    });
  }

  async deleteChat(chatId) {
    return this.request(`/chats/${chatId}`, {
      method: 'DELETE',
    });
  }

  async getChatMessages(chatId, limit = 50, offset = 0) {
    return this.request(`/chats/${chatId}/messages?limit=${limit}&offset=${offset}`);
  }

  async addChatMember(chatId, userId) {
    return this.request(`/chats/${chatId}/members`, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    });
  }

  async removeChatMember(chatId, userId) {
    return this.request(`/chats/${chatId}/members/${userId}`, {
      method: 'DELETE',
    });
  }

  // Messages API
  async sendMessage(chatId, messageData) {
    return this.request(`/chats/${chatId}/messages`, {
      method: 'POST',
      body: JSON.stringify(messageData),
    });
  }

  async getMessageById(messageId) {
    return this.request(`/messages/${messageId}`);
  }

  async updateMessageStatus(messageId, status) {
    return this.request(`/messages/${messageId}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status }),
    });
  }

  // Bots API
  async getBots() {
    return this.request('/bots');
  }

  async createBot(botData) {
    return this.request('/bots', {
      method: 'POST',
      body: JSON.stringify(botData),
    });
  }

  async getBotById(botId) {
    return this.request(`/bots/${botId}`);
  }

  async updateBot(botId, botData) {
    return this.request(`/bots/${botId}`, {
      method: 'PUT',
      body: JSON.stringify(botData),
    });
  }

  async deleteBot(botId) {
    return this.request(`/bots/${botId}`, {
      method: 'DELETE',
    });
  }

  async sendBotMessage(botId, messageData) {
    return this.request(`/bots/${botId}/sendMessage`, {
      method: 'POST',
      body: JSON.stringify(messageData),
    });
  }
}

export default new ApiService();
