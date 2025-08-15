import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

const useStore = create(
  devtools(
    (set, get) => ({
      // Состояние
      currentUser: null,
      currentChat: null,
      chats: [],
      messages: {},
      users: [],
      bots: [],
      debugEvents: [],
      statistics: {
        messages_count: 0,
        response_time: 0,
        errors_count: 0,
        users_count: 0,
        chats_count: 0
      },
      isLoading: false,
      error: null,
      isConnected: false,

      // Действия
      setCurrentUser: (user) => set({ currentUser: user }),
      
      setCurrentChat: (chat) => set({ currentChat: chat }),
      
      setChats: (chats) => set({ chats }),
      
      addChat: (chat) => set((state) => ({
        chats: [...state.chats, chat]
      })),
      
      updateChat: (chatId, updates) => set((state) => ({
        chats: state.chats.map(chat => 
          chat.id === chatId ? { ...chat, ...updates } : chat
        )
      })),
      
      setMessages: (chatId, messages) => set((state) => ({
        messages: { ...state.messages, [chatId]: messages }
      })),
      
      addMessage: (chatId, message) => set((state) => ({
        messages: {
          ...state.messages,
          [chatId]: [...(state.messages[chatId] || []), message]
        }
      })),
      
      updateMessage: (chatId, messageId, updates) => set((state) => ({
        messages: {
          ...state.messages,
          [chatId]: (state.messages[chatId] || []).map(message =>
            message.id === messageId ? { ...message, ...updates } : message
          )
        }
      })),
      
      setUsers: (users) => set({ users }),
      
      addUser: (user) => set((state) => ({
        users: [...state.users, user]
      })),
      
      updateUser: (userId, updates) => set((state) => ({
        users: state.users.map(user => 
          user.id === userId ? { ...user, ...updates } : user
        )
      })),
      
      setBots: (bots) => set({ bots }),
      
      addBot: (bot) => set((state) => ({
        bots: [...state.bots, bot]
      })),
      
      updateBot: (botId, updates) => set((state) => ({
        bots: state.bots.map(bot => 
          bot.id === botId ? { ...bot, ...updates } : bot
        )
      })),
      
      addDebugEvent: (event) => set((state) => ({
        debugEvents: [event, ...state.debugEvents.slice(0, 99)] // Ограничиваем 100 событиями
      })),
      
      clearDebugEvents: () => set({ debugEvents: [] }),
      
      setStatistics: (statistics) => set({ statistics }),
      
      updateStatistics: (updates) => set((state) => ({
        statistics: { ...state.statistics, ...updates }
      })),
      
      setLoading: (isLoading) => set({ isLoading }),
      
      setError: (error) => set({ error }),
      
      setConnected: (isConnected) => set({ isConnected }),
      
      // Селекторы
      getChatById: (chatId) => {
        const state = get();
        return state.chats.find(chat => chat.id === chatId);
      },
      
      getMessagesByChatId: (chatId) => {
        const state = get();
        return state.messages[chatId] || [];
      },
      
      getUserById: (userId) => {
        const state = get();
        return state.users.find(user => user.id === userId);
      },
      
      getBotById: (botId) => {
        const state = get();
        return state.bots.find(bot => bot.id === botId);
      },
      
      // Сброс состояния
      reset: () => set({
        currentUser: null,
        currentChat: null,
        chats: [],
        messages: {},
        users: [],
        bots: [],
        debugEvents: [],
        statistics: {
          messages_count: 0,
          response_time: 0,
          errors_count: 0,
          users_count: 0,
          chats_count: 0
        },
        isLoading: false,
        error: null,
        isConnected: false
      })
    }),
    {
      name: 'telegram-emulator-store'
    }
  )
);

export default useStore;
