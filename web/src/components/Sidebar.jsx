import React, { useState } from 'react';
import { 
  Search, 
  MessageCircle, 
  Users, 
  Settings, 
  Bug, 
  Wifi, 
  WifiOff,
  Plus,
  User
} from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import clsx from 'clsx';

const Sidebar = ({ 
  chats, 
  currentChat, 
  currentUser, 
  isConnected, 
  onChatSelect, 
  onToggleDebug 
}) => {
  const [searchQuery, setSearchQuery] = useState('');

  const filteredChats = chats.filter(chat => 
    chat.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    chat.username?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const getChatAvatar = (chat) => {
    if (chat.type === 'private' && chat.members?.length > 0) {
      const member = chat.members.find(m => m.id !== currentUser?.id) || chat.members[0];
      return member.first_name?.charAt(0).toUpperCase() || '?';
    }
    return chat.title?.charAt(0).toUpperCase() || '?';
  };

  const getChatTitle = (chat) => {
    if (chat.type === 'private' && chat.members?.length > 0) {
      const member = chat.members.find(m => m.id !== currentUser?.id) || chat.members[0];
      return `${member.first_name} ${member.last_name || ''}`.trim();
    }
    return chat.title;
  };

  const getLastMessageText = (chat) => {
    if (!chat.last_message) return 'Нет сообщений';
    
    const text = chat.last_message.text;
    return text.length > 30 ? `${text.substring(0, 30)}...` : text;
  };

  const formatTime = (timestamp) => {
    try {
      const date = new Date(timestamp);
      const now = new Date();
      const diffInHours = (now - date) / (1000 * 60 * 60);
      
      if (diffInHours < 24) {
        return format(date, 'HH:mm', { locale: ru });
      } else if (diffInHours < 168) { // 7 days
        return format(date, 'EEE', { locale: ru });
      } else {
        return format(date, 'dd.MM.yy', { locale: ru });
      }
    } catch (error) {
      return '';
    }
  };

  return (
    <div className="w-80 bg-telegram-sidebar border-r border-telegram-border flex flex-col">
      {/* Заголовок */}
      <div className="p-4 border-b border-telegram-border bg-telegram-bg">
        <h1 className="text-lg font-medium text-telegram-text mb-3">
          Telegram Emulator
        </h1>
        
        {/* Поиск */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-telegram-secondary w-4 h-4" />
          <input
            type="text"
            placeholder="Поиск чатов..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-2 bg-telegram-sidebar border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
          />
        </div>
      </div>

      {/* Список чатов */}
      <div className="flex-1 overflow-y-auto scrollbar-hide">
        {filteredChats.length === 0 ? (
          <div className="p-8 text-center">
            <MessageCircle className="w-12 h-12 text-telegram-secondary mx-auto mb-3" />
            <h3 className="text-telegram-text font-medium mb-1">Нет чатов</h3>
            <p className="text-telegram-text-secondary text-sm">
              {searchQuery ? 'По вашему запросу ничего не найдено' : 'Создайте чат для начала общения'}
            </p>
          </div>
        ) : (
          <div>
            {filteredChats.map((chat) => (
              <div
                key={chat.id}
                onClick={() => onChatSelect(chat)}
                className={clsx(
                  'chat-item',
                  currentChat?.id === chat.id && 'active'
                )}
              >
                {/* Аватар */}
                <div className="w-12 h-12 rounded-full bg-telegram-primary flex items-center justify-center text-white font-medium mr-3 flex-shrink-0">
                  {getChatAvatar(chat)}
                </div>

                {/* Информация о чате */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between mb-1">
                    <h3 className="text-telegram-text font-medium truncate">
                      {getChatTitle(chat)}
                    </h3>
                    {chat.last_message && (
                      <span className="text-xs text-telegram-text-secondary flex-shrink-0 ml-2">
                        {formatTime(chat.last_message.timestamp)}
                      </span>
                    )}
                  </div>
                  
                  <div className="flex items-center justify-between">
                    <p className="text-sm text-telegram-text-secondary truncate">
                      {getLastMessageText(chat)}
                    </p>
                    {chat.unread_count > 0 && (
                      <span className="bg-telegram-primary text-white text-xs px-2 py-1 rounded-full min-w-[20px] text-center flex-shrink-0 ml-2">
                        {chat.unread_count > 99 ? '99+' : chat.unread_count}
                      </span>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Нижняя панель */}
      <div className="p-4 border-t border-telegram-border bg-telegram-bg">
        {/* Статус соединения */}
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center">
            {isConnected ? (
              <Wifi className="w-4 h-4 text-green-500 mr-2" />
            ) : (
              <WifiOff className="w-4 h-4 text-red-500 mr-2" />
            )}
            <span className={clsx(
              'text-sm',
              isConnected ? 'text-green-500' : 'text-red-500'
            )}>
              {isConnected ? 'Подключено' : 'Отключено'}
            </span>
          </div>
        </div>

        {/* Кнопки действий */}
        <div className="flex space-x-2">
          <button
            onClick={onToggleDebug}
            className="flex-1 flex items-center justify-center px-3 py-2 bg-telegram-sidebar border border-telegram-border rounded-lg text-telegram-text hover:bg-telegram-primary hover:border-telegram-primary transition-colors"
          >
            <Bug className="w-4 h-4 mr-2" />
            Отладка
          </button>
          
          <button className="flex-1 flex items-center justify-center px-3 py-2 bg-telegram-sidebar border border-telegram-border rounded-lg text-telegram-text hover:bg-telegram-primary hover:border-telegram-primary transition-colors">
            <Settings className="w-4 h-4 mr-2" />
            Настройки
          </button>
        </div>

        {/* Информация о пользователе */}
        {currentUser && (
          <div className="mt-3 p-3 bg-telegram-sidebar rounded-lg">
            <div className="flex items-center">
              <div className="w-8 h-8 rounded-full bg-telegram-primary flex items-center justify-center text-white font-medium mr-2">
                {currentUser.first_name?.charAt(0).toUpperCase() || 'U'}
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-telegram-text truncate">
                  {currentUser.first_name} {currentUser.last_name || ''}
                </p>
                <p className="text-xs text-telegram-text-secondary truncate">
                  @{currentUser.username}
                </p>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Sidebar;
