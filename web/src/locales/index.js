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
    createChat: 'Создать чат',
    chatType: 'Тип чата',
    title: 'Название',
    description: 'Описание',
    participants: 'Участники',
    selectTwo: '(выберите 2)',
    minimumTwo: '(минимум 2)',
    selected: 'Выбрано',
    chatTitleRequired: 'Название чата обязательно',
    privateChatTwoParticipants: 'Приватный чат должен содержать ровно 2 участника',
    groupMinTwoParticipants: 'Группа должна содержать минимум 2 участника',
    chatCreationError: 'Ошибка создания чата',
    
    // Пользователи
    users: 'Пользователи',
    newUser: 'Новый пользователь',
    user: 'Пользователь',
    username: 'Имя пользователя',
    firstName: 'Имя',
    lastName: 'Фамилия',
    isBot: 'Бот',
    selectUser: 'Выберите пользователя',
    createUser: 'Создать пользователя',
    createNewUser: 'Создать нового пользователя',
    confirmDeleteUser: 'Вы уверены, что хотите удалить этого пользователя?',
    usernameAndFirstNameRequired: 'Имя пользователя и имя обязательны',
    userCreationError: 'Ошибка создания пользователя',
    creating: 'Создание...',
    you: 'Вы',
    unknown: 'Неизвестный',
    file: 'Файл',
    voiceMessage: 'Голосовое сообщение',
    photo: 'Фото',
    
    // Боты
    bots: 'Боты',
    newBot: 'Новый бот',
    botName: 'Имя бота',
    botToken: 'Токен бота',
    webhookUrl: 'URL вебхука',
    isActive: 'Активен',
    botManagement: 'Управление ботами',
    createBot: 'Создать бота',
    noBots: 'Нет ботов',
    createFirstBot: 'Создайте первого бота для начала работы',
    loadingBots: 'Загрузка ботов...',
    botsLoadError: 'Ошибка загрузки ботов',
    confirmDeleteBot: 'Вы уверены, что хотите удалить этого бота?',
    botToggleError: 'Ошибка изменения статуса бота',
    token: 'Токен',
    copyToken: 'Копировать токен',
    copyUrl: 'Копировать URL',
    active: 'Активен',
    inactive: 'Неактивен',
    created: 'Создан',
    activate: 'Активировать',
    deactivate: 'Деактивировать',
    nameAndUsernameRequired: 'Название и username обязательны',
    botCreationError: 'Ошибка создания бота',
    events: 'События',
    statistics: 'Статистика',
    
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
    createChat: 'Create Chat',
    chatType: 'Chat Type',
    title: 'Title',
    description: 'Description',
    participants: 'Participants',
    selectTwo: '(select 2)',
    minimumTwo: '(minimum 2)',
    selected: 'Selected',
    chatTitleRequired: 'Chat title is required',
    privateChatTwoParticipants: 'Private chat must contain exactly 2 participants',
    groupMinTwoParticipants: 'Group must contain at least 2 participants',
    chatCreationError: 'Error creating chat',
    
    // Users
    users: 'Users',
    newUser: 'New User',
    user: 'User',
    username: 'Username',
    firstName: 'First Name',
    lastName: 'Last Name',
    isBot: 'Bot',
    selectUser: 'Select User',
    createUser: 'Create User',
    createNewUser: 'Create New User',
    confirmDeleteUser: 'Are you sure you want to delete this user?',
    usernameAndFirstNameRequired: 'Username and first name are required',
    userCreationError: 'Error creating user',
    creating: 'Creating...',
    you: 'You',
    unknown: 'Unknown',
    file: 'File',
    voiceMessage: 'Voice message',
    photo: 'Photo',
    
    // Bots
    bots: 'Bots',
    newBot: 'New Bot',
    botName: 'Bot Name',
    botToken: 'Bot Token',
    webhookUrl: 'Webhook URL',
    isActive: 'Active',
    botManagement: 'Bot Management',
    createBot: 'Create Bot',
    noBots: 'No bots',
    createFirstBot: 'Create your first bot to get started',
    loadingBots: 'Loading bots...',
    botsLoadError: 'Error loading bots',
    confirmDeleteBot: 'Are you sure you want to delete this bot?',
    botToggleError: 'Error changing bot status',
    token: 'Token',
    copyToken: 'Copy token',
    copyUrl: 'Copy URL',
    active: 'Active',
    inactive: 'Inactive',
    created: 'Created',
    activate: 'Activate',
    deactivate: 'Deactivate',
    nameAndUsernameRequired: 'Name and username are required',
    botCreationError: 'Error creating bot',
    events: 'Events',
    statistics: 'Statistics',
    
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
