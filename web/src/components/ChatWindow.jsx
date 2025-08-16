import React, { useState, useRef, useEffect } from 'react';
import { Send, Paperclip, Mic, Smile } from 'lucide-react';
import { format } from 'date-fns';
import { ru, enUS } from 'date-fns/locale';
import clsx from 'clsx';
import MessageBubble from './MessageBubble';
import { t, getCurrentLanguage } from '../locales';

const ChatWindow = ({ chat, messages, currentUser, onSendMessage }) => {
  const [inputText, setInputText] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const messagesEndRef = useRef(null);
  const inputRef = useRef(null);

  // Автоскролл к последнему сообщению
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // Фокус на поле ввода при смене чата
  useEffect(() => {
    if (chat) {
      inputRef.current?.focus();
    }
  }, [chat?.id]);

  const handleSendMessage = () => {
    if (!inputText.trim() || !chat) return;
    
    onSendMessage(inputText);
    setInputText('');
    setIsTyping(false);
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const handleInputChange = (e) => {
    setInputText(e.target.value);
    setIsTyping(e.target.value.length > 0);
  };

  const getChatTitle = () => {
    if (!chat) return '';
    
    if (chat.type === 'private' && chat.members?.length > 0) {
      const member = chat.members.find(m => m.id !== currentUser?.id) || chat.members[0];
      return `${member.first_name} ${member.last_name || ''}`.trim();
    }
    return chat.title;
  };

  const getChatAvatar = () => {
    if (!chat) return '';
    
    if (chat.type === 'private' && chat.members?.length > 0) {
      const member = chat.members.find(m => m.id !== currentUser?.id) || chat.members[0];
      return member.first_name?.charAt(0).toUpperCase() || '?';
    }
    return chat.title?.charAt(0).toUpperCase() || '?';
  };

  const getOnlineStatus = () => {
    if (!chat || chat.type !== 'private') return null;
    
    const member = chat.members?.find(m => m.id !== currentUser?.id);
    if (!member) return null;
    
    return member.is_online ? 'online' : 'offline';
  };

  if (!chat) {
    const language = getCurrentLanguage();
    return (
      <div className="flex-1 flex items-center justify-center bg-telegram-bg">
        <div className="text-center">
          <div className="text-6xl mb-4">💬</div>
          <h2 className="text-xl font-medium text-telegram-text mb-2">
            {t('selectChat', language)}
          </h2>
          <p className="text-telegram-text-secondary">
            {t('selectChatFromList', language)}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 flex flex-col bg-telegram-bg h-full">
      {/* Заголовок чата */}
      <div className="flex items-center p-4 border-b border-telegram-border bg-telegram-sidebar">
        <div className="w-10 h-10 rounded-full bg-gradient-to-br from-telegram-primary to-blue-600 flex items-center justify-center text-white font-medium mr-3 shadow-sm">
          {getChatAvatar()}
        </div>
        
        <div className="flex-1 min-w-0">
          <h2 className="text-telegram-text font-medium truncate">
            {getChatTitle()}
          </h2>
          {getOnlineStatus() && (
            <p className={clsx(
              'text-sm',
              getOnlineStatus() === 'online' ? 'text-green-500' : 'text-telegram-text-secondary'
            )}>
              {getOnlineStatus() === 'online' ? 'online' : 'offline'}
            </p>
          )}
        </div>
      </div>

      {/* Область сообщений */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4" style={{ height: 'calc(100vh - 200px)' }}>
        {messages.length === 0 ? (
          <div className="flex items-center justify-center h-full">
            <div className="text-center">
              <div className="text-4xl mb-3">📱</div>
              <h3 className="text-telegram-text font-medium mb-1">
                {t('noMessages', getCurrentLanguage())}
              </h3>
              <p className="text-telegram-text-secondary text-sm">
                {t('startConversation', getCurrentLanguage())}
              </p>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            {messages.map((message) => (
              <MessageBubble
                key={message.id}
                message={message}
                isOwn={message.from?.id === currentUser?.id || message.from_id === currentUser?.id || message.is_outgoing}
                currentUser={currentUser}
              />
            ))}
            <div ref={messagesEndRef} />
          </div>
        )}
      </div>

      {/* Поле ввода */}
      <div className="p-4 border-t border-telegram-border bg-telegram-sidebar">
        <div className="flex items-end space-x-2">
          {/* Кнопки действий */}
          <div className="flex space-x-1">
            <button className="p-2 text-telegram-secondary hover:text-telegram-text transition-colors">
              <Paperclip className="w-5 h-5" />
            </button>
            <button className="p-2 text-telegram-secondary hover:text-telegram-text transition-colors">
              <Smile className="w-5 h-5" />
            </button>
          </div>

          {/* Поле ввода */}
          <div className="flex-1 relative">
            <textarea
              ref={inputRef}
              value={inputText}
              onChange={handleInputChange}
              onKeyPress={handleKeyPress}
              placeholder={t('messagePlaceholder', getCurrentLanguage())}
              rows="1"
              className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary resize-none max-h-32"
              style={{
                minHeight: '40px',
                maxHeight: '120px'
              }}
            />
          </div>

          {/* Кнопка отправки */}
          <div className="flex space-x-1">
            <button className="p-2 text-telegram-secondary hover:text-telegram-text transition-colors">
              <Mic className="w-5 h-5" />
            </button>
            <button
              onClick={handleSendMessage}
              disabled={!inputText.trim()}
              className={clsx(
                'p-2 rounded-lg transition-colors',
                inputText.trim()
                  ? 'bg-telegram-primary text-white hover:bg-telegram-primary/80'
                  : 'bg-telegram-bg text-telegram-secondary cursor-not-allowed'
              )}
            >
              <Send className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Индикатор печати убран - не нужен для отправителя */}
      </div>
    </div>
  );
};

export default ChatWindow;
