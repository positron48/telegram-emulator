// Локализация приложения
export const locales = {
  ru: {
    // Общие
    settings: 'Настройки',
    save: 'Сохранить',
    cancel: 'Отмена',
    close: 'Закрыть',
    delete: 'Удалить',
    edit: 'Редактировать',
    create: 'Создать',
    loading: 'Загрузка...',
    error: 'Ошибка',
    success: 'Успешно',
    
    // Настройки
    theme: 'Тема оформления',
    lightTheme: 'Светлая',
    darkTheme: 'Темная',
    systemTheme: 'Системная',
    language: 'Язык интерфейса',
    russian: 'Русский',
    english: 'English',
    additional: 'Дополнительно',
    autoScroll: 'Автоматическая прокрутка к новым сообщениям',
    debugMode: 'Режим отладки',
    showTimestamps: 'Показывать время сообщений',
    websocketSettings: 'WebSocket настройки',
    reconnectAttempts: 'Количество попыток переподключения',
    reconnectDelay: 'Задержка переподключения (мс)',
    messageSettings: 'Настройки сообщений',
    messageHistoryLimit: 'Лимит истории сообщений',
    exportImport: 'Экспорт/Импорт данных',
    export: 'Экспорт',
    import: 'Импорт',
    exportImportDescription: 'Экспорт/импорт включает пользователей, чаты, сообщения и настройки',
    
    // Чаты
    chats: 'Чаты',
    newChat: 'Новый чат',
    privateChat: 'Приватный чат',
    groupChat: 'Группа',
    channelChat: 'Канал',
    noMessages: 'Нет сообщений',
    noChats: 'Нет чатов',
    searchChats: 'Поиск чатов...',
    noSearchResults: 'По вашему запросу ничего не найдено',
    createChatToStart: 'Создайте чат для начала общения',
    sendMessage: 'Отправить сообщение',
    messagePlaceholder: 'Введите сообщение...',
    selectChat: 'Выберите чат',
    selectChatFromList: 'Выберите чат из списка слева для начала общения',
    startConversation: 'Начните общение, отправив первое сообщение',
    
    // Пользователи
    users: 'Пользователи',
    newUser: 'Новый пользователь',
    username: 'Имя пользователя',
    firstName: 'Имя',
    lastName: 'Фамилия',
    isBot: 'Бот',
    selectUser: 'Выберите пользователя',
    
    // Боты
    bots: 'Боты',
    newBot: 'Новый бот',
    botName: 'Имя бота',
    botToken: 'Токен бота',
    webhookUrl: 'URL вебхука',
    isActive: 'Активен',
    
    // Статусы
    connected: 'Подключено',
    disconnected: 'Отключено',
    reconnecting: 'Переподключение...',
    reconnect: 'Переподключиться',
    
    // Действия
    debug: 'Отладка',
    settings: 'Настройки',
    deleteChat: 'Удалить чат',
    deleteUser: 'Удалить пользователя',
    deleteBot: 'Удалить бота',
    confirmDelete: 'Вы уверены, что хотите удалить',
    
    // Время
    today: 'Сегодня',
    yesterday: 'Вчера',
    justNow: 'Только что',
    
    // Ошибки
    failedToLoad: 'Не удалось загрузить',
    failedToSave: 'Не удалось сохранить',
    failedToDelete: 'Не удалось удалить',
    networkError: 'Ошибка сети',
    unknownError: 'Неизвестная ошибка'
  },
  
  en: {
    // General
    settings: 'Settings',
    save: 'Save',
    cancel: 'Cancel',
    close: 'Close',
    delete: 'Delete',
    edit: 'Edit',
    create: 'Create',
    loading: 'Loading...',
    error: 'Error',
    success: 'Success',
    
    // Settings
    theme: 'Theme',
    lightTheme: 'Light',
    darkTheme: 'Dark',
    systemTheme: 'System',
    language: 'Language',
    russian: 'Русский',
    english: 'English',
    additional: 'Additional',
    autoScroll: 'Auto-scroll to new messages',
    debugMode: 'Debug mode',
    showTimestamps: 'Show message timestamps',
    websocketSettings: 'WebSocket Settings',
    reconnectAttempts: 'Reconnection attempts',
    reconnectDelay: 'Reconnection delay (ms)',
    messageSettings: 'Message Settings',
    messageHistoryLimit: 'Message history limit',
    exportImport: 'Export/Import Data',
    export: 'Export',
    import: 'Import',
    exportImportDescription: 'Export/import includes users, chats, messages and settings',
    
    // Chats
    chats: 'Chats',
    newChat: 'New Chat',
    privateChat: 'Private Chat',
    groupChat: 'Group',
    channelChat: 'Channel',
    noMessages: 'No messages',
    noChats: 'No chats',
    searchChats: 'Search chats...',
    noSearchResults: 'No results found for your query',
    createChatToStart: 'Create a chat to start communicating',
    sendMessage: 'Send message',
    messagePlaceholder: 'Type a message...',
    selectChat: 'Select chat',
    selectChatFromList: 'Select a chat from the list on the left to start communicating',
    startConversation: 'Start a conversation by sending the first message',
    
    // Users
    users: 'Users',
    newUser: 'New User',
    username: 'Username',
    firstName: 'First Name',
    lastName: 'Last Name',
    isBot: 'Bot',
    selectUser: 'Select User',
    
    // Bots
    bots: 'Bots',
    newBot: 'New Bot',
    botName: 'Bot Name',
    botToken: 'Bot Token',
    webhookUrl: 'Webhook URL',
    isActive: 'Active',
    
    // Status
    connected: 'Connected',
    disconnected: 'Disconnected',
    reconnecting: 'Reconnecting...',
    reconnect: 'Reconnect',
    
    // Actions
    debug: 'Debug',
    settings: 'Settings',
    deleteChat: 'Delete Chat',
    deleteUser: 'Delete User',
    deleteBot: 'Delete Bot',
    confirmDelete: 'Are you sure you want to delete',
    
    // Time
    today: 'Today',
    yesterday: 'Yesterday',
    justNow: 'Just now',
    
    // Errors
    failedToLoad: 'Failed to load',
    failedToSave: 'Failed to save',
    failedToDelete: 'Failed to delete',
    networkError: 'Network error',
    unknownError: 'Unknown error'
  }
};

// Функция для получения перевода
export const t = (key, language = 'ru') => {
  const locale = locales[language] || locales.ru;
  return locale[key] || key;
};

// Функция для изменения языка
export const setLanguage = (language) => {
  document.documentElement.lang = language;
  localStorage.setItem('telegram-emulator-language', language);
};

// Функция для получения текущего языка
export const getCurrentLanguage = () => {
  return localStorage.getItem('telegram-emulator-language') || 'ru';
};
