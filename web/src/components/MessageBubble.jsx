import React from 'react';
import { format } from 'date-fns';
import { ru, enUS } from 'date-fns/locale';
import { Check, CheckCheck } from 'lucide-react';
import clsx from 'clsx';
import { t, getCurrentLanguage } from '../locales';
import { parseTelegramText } from '../utils/textParser.jsx';

const MessageBubble = ({ message, isOwn, currentUser, onSendMessage, onCallbackQuery }) => {
  const formatTime = (timestamp) => {
    try {
      const language = getCurrentLanguage();
      const locale = language === 'en' ? enUS : ru;
      const date = new Date(timestamp);
      return format(date, 'HH:mm', { locale });
    } catch (error) {
      return '';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'sending':
        return <div className="w-4 h-4 border-2 border-telegram-secondary border-t-transparent rounded-full animate-spin" />;
      case 'sent':
        return <Check className="w-4 h-4 text-telegram-secondary" />;
      case 'delivered':
        return <CheckCheck className="w-4 h-4 text-telegram-secondary" />;
      case 'read':
        return <CheckCheck className="w-4 h-4 text-blue-500" />;
      default:
        return null;
    }
  };

  const getSenderName = () => {
    const language = getCurrentLanguage();
    if (isOwn) return t('you', language);
    return message.from?.first_name || message.from?.username || t('unknown', language);
  };

  const getMessageContent = () => {
    const language = getCurrentLanguage();
    switch (message.type) {
      case 'text':
        return (
          <div className="whitespace-pre-wrap break-words">
            {renderTextWithCommands(message.text)}
          </div>
        );
      case 'file':
        return (
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-telegram-secondary rounded flex items-center justify-center">
              📎
            </div>
            <span>{t('file', language)}</span>
          </div>
        );
      case 'voice':
        return (
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-telegram-secondary rounded flex items-center justify-center">
              🎤
            </div>
            <span>{t('voiceMessage', language)}</span>
          </div>
        );
      case 'photo':
        return (
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-telegram-secondary rounded flex items-center justify-center">
              📷
            </div>
            <span>{t('photo', language)}</span>
          </div>
        );
      default:
        return <div>{renderTextWithCommands(message.text)}</div>;
    }
  };

  const renderTextWithCommands = (text) => {
    if (!text) return null;
    
    // Сначала обрабатываем URL-адреса
    const urlRegex = /(https?:\/\/[^\s]+)/g;
    let parts = text.split(urlRegex);
    
    // Обрабатываем каждую часть на предмет команд и форматирования
    parts = parts.map((part, index) => {
      // Если это URL, делаем его кликабельным
      if (urlRegex.test(part)) {
        return (
          <a
            key={`url-${index}`}
            href={part}
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-500 hover:text-blue-600 dark:text-blue-400 dark:hover:text-blue-300 underline"
          >
            {part}
          </a>
        );
      }
      
      // Если это не URL, ищем команды
      // Команды должны быть в начале строки или после пробела
      const commandRegex = /(^|\s)(\/[a-zA-Z0-9_]+)/g;
      const commandParts = part.split(commandRegex);
      
      return commandParts.map((commandPart, commandIndex) => {
        if (commandRegex.test(commandPart)) {
          // Это пробел + команда, разделяем
          const space = commandPart.charAt(0) === ' ' ? ' ' : '';
          const command = commandPart.substring(space.length);
          
          return (
            <React.Fragment key={`command-${index}-${commandIndex}`}>
              {space}
              <button
                style={{
                  color: 'rgb(59, 130, 246)',
                  cursor: 'pointer',
                  fontWeight: '500',
                  transition: 'color 0.2s'
                }}
                className="hover:text-blue-600 dark:hover:text-[rgb(156,197,255)] inline-flex items-center"
                onClick={() => {
                  if (onSendMessage) {
                    onSendMessage(command);
                  }
                  console.log('Command clicked:', command);
                }}
                title={`Click to send: ${command}`}
              >
                <span 
                  style={{
                    color: 'rgb(59, 130, 246)',
                    fontWeight: '500'
                  }}
                  className="dark:text-[rgb(121,171,252)]"
                >
                  {command}
                </span>
              </button>
            </React.Fragment>
          );
        }
        
        // Применяем Telegram форматирование к обычному тексту
        return parseTelegramText(commandPart);
      });
    });
    
    return parts;
  };

  const renderInlineKeyboard = () => {
    if (!message.reply_markup || !message.reply_markup.inline_keyboard) return null;

    return (
      <div className="mt-2 space-y-1">
        {message.reply_markup.inline_keyboard.map((row, rowIndex) => (
          <div key={rowIndex} className="flex space-x-1">
            {row.map((button, buttonIndex) => (
              <button
                key={buttonIndex}
                className="px-3 py-2 bg-blue-500 text-white rounded-lg text-sm font-medium hover:bg-blue-600 transition-colors"
                onClick={() => {
                  // Для inline клавиатуры отправляем callback_query
                  if (onCallbackQuery) {
                    onCallbackQuery(button);
                  }
                  console.log('Inline button clicked:', button);
                }}
              >
                {button.text}
              </button>
            ))}
          </div>
        ))}
      </div>
    );
  };

  return (
    <div className={clsx(
      'flex',
      isOwn ? 'justify-end' : 'justify-start'
    )}>
      <div className={clsx(
        'max-w-xs lg:max-w-md',
        isOwn ? 'order-2' : 'order-1'
      )}>
        {/* Sender name (only for other messages in group chats) */}
        {!isOwn && message.chat?.type === 'group' && (
          <div className="text-xs text-telegram-text-secondary mb-1 ml-1">
            {getSenderName()}
          </div>
        )}

        {/* Message bubble */}
        <div className={clsx(
          'message-bubble',
          isOwn ? 'outgoing' : 'incoming'
        )}>
          {getMessageContent()}
          {renderInlineKeyboard()}
        </div>

        {/* Time and status */}
        <div className={clsx(
          'flex items-center mt-1 space-x-1',
          isOwn ? 'justify-end' : 'justify-start'
        )}>
          <span className="text-xs text-telegram-text-secondary">
            {formatTime(message.timestamp)}
          </span>
          
          {isOwn && (
            <div className="flex items-center">
              {getStatusIcon(message.status)}
            </div>
          )}
        </div>
      </div>

      {/* Avatar (only for other messages) */}
      {!isOwn && (
        <div className={clsx(
          'w-8 h-8 rounded-full bg-gradient-to-br from-telegram-primary to-blue-600 flex items-center justify-center text-white text-sm font-medium ml-2 order-2 flex-shrink-0 shadow-sm',
          'self-end mb-1'
        )}>
          {message.from?.first_name?.charAt(0).toUpperCase() || '?'}
        </div>
      )}
    </div>
  );
};

export default MessageBubble;
