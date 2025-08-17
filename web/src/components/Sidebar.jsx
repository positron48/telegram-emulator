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
  User,
  Bot,
  Trash2
} from 'lucide-react';
import UserSelector from './UserSelector';
import { format } from 'date-fns';
import { ru, enUS } from 'date-fns/locale';
import clsx from 'clsx';
import { t, getCurrentLanguage } from '../locales';

const Sidebar = ({ 
  chats, 
  currentChat, 
  currentUser, 
  users, 
  isConnected, 
  isReconnecting, 
  onChatSelect, 
  onToggleDebug, 
  onUserSelect, 
  onCreateUser, 
  onDeleteUser, 
  onCreateChat, 
  onDeleteChat, 
  onShowBotManager, 
  onShowSettings, 
  onReconnect 
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

  const getChatIcon = (chat) => {
    switch (chat.type) {
      case 'private':
        return '👤';
      case 'group':
        return '👥';
      default:
        return '💬';
    }
  };

  const getChatTitle = (chat) => {
    if (chat.type === 'private' && chat.members?.length > 0) {
      const member = chat.members.find(m => m.id !== currentUser?.id) || chat.members[0];
      return `${member.first_name} ${member.last_name || ''}`.trim();
    }
    return chat.title;
  };

  const getLastMessageText = (chat) => {
    const language = getCurrentLanguage();
    if (!chat.last_message) return t('noMessages', language);
    
    const text = chat.last_message.text;
    return text.length > 30 ? `${text.substring(0, 30)}...` : text;
  };

  const formatTime = (timestamp) => {
    try {
      const language = getCurrentLanguage();
      const locale = language === 'en' ? enUS : ru;
      const date = new Date(timestamp);
      const now = new Date();
      const diffInHours = (now - date) / (1000 * 60 * 60);
      
      if (diffInHours < 24) {
        return format(date, 'HH:mm', { locale });
      } else if (diffInHours < 168) { // 7 days
        return format(date, 'EEE', { locale });
      } else {
        return format(date, 'dd.MM.yy', { locale });
      }
    } catch (error) {
      return '';
    }
  };

  return (
    <div className="w-80 bg-telegram-sidebar border-r border-telegram-border flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-telegram-border bg-telegram-bg">
        <h1 className="text-lg font-medium text-telegram-text mb-3">
          Telegram Emulator
        </h1>
        
        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-telegram-secondary w-4 h-4" />
          <input
            type="text"
            placeholder={t('searchChats', getCurrentLanguage())}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-2 bg-telegram-sidebar border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
          />
        </div>
      </div>

      {/* Chats list */}
      <div className="flex-1 overflow-y-auto scrollbar-hide">
        {filteredChats.length === 0 ? (
          <div className="p-8 text-center">
            <MessageCircle className="w-12 h-12 text-telegram-secondary mx-auto mb-3" />
            <h3 className="text-telegram-text font-medium mb-1">{t('noChats', getCurrentLanguage())}</h3>
            <p className="text-telegram-text-secondary text-sm">
              {searchQuery ? t('noSearchResults', getCurrentLanguage()) : t('createChatToStart', getCurrentLanguage())}
            </p>
          </div>
        ) : (
          <div>
            {filteredChats.map((chat) => (
              <div
                key={chat.id}
                className={clsx(
                  'chat-item relative group',
                  currentChat?.id === chat.id && 'active'
                )}
              >
                <div
                  onClick={() => onChatSelect(chat)}
                  className="flex items-center p-2 hover:bg-telegram-primary/10 transition-colors cursor-pointer"
                >
                {/* Avatar */}
                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-telegram-primary to-blue-600 flex items-center justify-center text-white font-medium mr-3 flex-shrink-0 shadow-sm">
                  {chat.type === 'private' ? getChatAvatar(chat) : getChatIcon(chat)}
                </div>

                {/* Chat information */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <h3 className="text-telegram-text font-medium truncate text-sm">
                      {getChatTitle(chat)}
                    </h3>
                    {chat.last_message && (
                      <span className="text-xs text-telegram-text-secondary flex-shrink-0 ml-2">
                        {formatTime(chat.last_message.timestamp)}
                      </span>
                    )}
                  </div>
                  
                  <div className="flex items-center justify-between mt-0.5">
                    <p className="text-xs text-telegram-text-secondary truncate">
                      {getLastMessageText(chat)}
                    </p>
                    {chat.unread_count > 0 && (
                      <span className={clsx(
                        "text-white text-xs px-1.5 py-0.5 rounded-full min-w-[16px] text-center flex-shrink-0 ml-2 shadow-sm font-medium",
                        currentChat?.id === chat.id 
                          ? "bg-red-600 shadow-md" 
                          : "bg-red-500"
                      )}>
                        {chat.unread_count > 99 ? '99+' : chat.unread_count}
                      </span>
                    )}
                  </div>
                </div>

                </div>

                {/* Delete chat button */}
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    const language = getCurrentLanguage();
                    if (confirm(`${t('confirmDelete', language)} "${getChatTitle(chat)}"?`)) {
                      onDeleteChat(chat.id);
                    }
                  }}
                  className="absolute top-2 right-2 p-1 text-red-500 hover:text-red-600 opacity-0 group-hover:opacity-100 transition-opacity"
                  title={t('deleteChat', getCurrentLanguage())}
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Bottom panel */}
      <div className="p-4 border-t border-telegram-border bg-telegram-bg">
        {/* Connection status */}
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center">
            {isConnected ? (
              <Wifi className="w-4 h-4 text-green-500 mr-2" />
            ) : isReconnecting ? (
              <Wifi className="w-4 h-4 text-yellow-500 mr-2 animate-pulse" />
            ) : (
              <WifiOff className="w-4 h-4 text-red-500 mr-2" />
            )}
            <span className={clsx(
              'text-sm',
              isConnected ? 'text-green-500' : isReconnecting ? 'text-yellow-500' : 'text-red-500'
            )}>
              {isConnected ? t('connected', getCurrentLanguage()) : isReconnecting ? t('reconnecting', getCurrentLanguage()) : t('disconnected', getCurrentLanguage())}
            </span>
          </div>
          
          {/* Reconnect button */}
          {!isConnected && !isReconnecting && onReconnect && (
            <button
              onClick={onReconnect}
              className="px-2 py-1 text-xs bg-telegram-primary text-white rounded hover:bg-telegram-primary/80 transition-colors"
              title={t('reconnect', getCurrentLanguage())}
            >
              {t('reconnect', getCurrentLanguage())}
            </button>
          )}
        </div>

        {/* Action buttons */}
        <div className="grid grid-cols-2 gap-2 mb-2">
          <button
            onClick={onToggleDebug}
            className="flex items-center justify-center px-3 py-2 bg-telegram-sidebar border border-telegram-border rounded-lg text-telegram-text hover:bg-telegram-primary hover:border-telegram-primary transition-colors"
          >
            <Bug className="w-4 h-4 mr-2" />
            {t('debug', getCurrentLanguage())}
          </button>
          
          <button
            onClick={onShowBotManager}
            className="flex items-center justify-center px-3 py-2 bg-telegram-sidebar border border-telegram-border rounded-lg text-telegram-text hover:bg-telegram-primary hover:border-telegram-primary transition-colors"
          >
            <Bot className="w-4 h-4 mr-2" />
            {t('bots', getCurrentLanguage())}
          </button>
        </div>

        <div className="grid grid-cols-2 gap-2 mb-3">
          <button
            onClick={onCreateChat}
            className="flex items-center justify-center px-3 py-2 bg-telegram-sidebar border border-telegram-border rounded-lg text-telegram-text hover:bg-telegram-primary hover:border-telegram-primary transition-colors"
          >
            <Plus className="w-4 h-4 mr-2" />
            {t('newChat', getCurrentLanguage())}
          </button>
          
          <button
            onClick={onShowSettings}
            className="flex items-center justify-center px-3 py-2 bg-telegram-sidebar border border-telegram-border rounded-lg text-telegram-text hover:bg-telegram-primary hover:border-telegram-primary transition-colors"
          >
            <Settings className="w-4 h-4 mr-2" />
            {t('settings', getCurrentLanguage())}
          </button>
        </div>

        {/* User selection */}
        <div className="mt-3">
          <UserSelector
            users={users}
            currentUser={currentUser}
            onUserSelect={onUserSelect}
            onCreateUser={onCreateUser}
            onDeleteUser={onDeleteUser}
          />
        </div>
      </div>
    </div>
  );
};

export default Sidebar;
