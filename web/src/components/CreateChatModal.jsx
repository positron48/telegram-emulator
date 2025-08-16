import React, { useState } from 'react';
import { X, Users, MessageCircle, Hash } from 'lucide-react';
import apiService from '../services/api';
import useStore from '../store';
import { t, getCurrentLanguage } from '../locales';

const CreateChatModal = ({ isOpen, onClose, onChatCreated }) => {
  const [formData, setFormData] = useState({
    type: 'private',
    title: '',
    username: '',
    description: '',
    user_ids: []
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const { users, addChat, addDebugEvent } = useStore();

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleUserToggle = (userId) => {
    setFormData(prev => ({
      ...prev,
      user_ids: prev.user_ids.includes(userId)
        ? prev.user_ids.filter(id => id !== userId)
        : [...prev.user_ids, userId]
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    const language = getCurrentLanguage();
    if (!formData.title.trim()) {
      setError(t('chatTitleRequired', language));
      return;
    }

    if (formData.type === 'private' && formData.user_ids.length !== 2) {
      setError(t('privateChatTwoParticipants', language));
      return;
    }

    if (formData.type !== 'private' && formData.user_ids.length < 2) {
      setError(t('groupMinTwoParticipants', language));
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const response = await apiService.createChat(formData);
      
      if (response.chat) {
        addChat(response.chat);
        addDebugEvent({
          id: `chat-created-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
          timestamp: new Date().toLocaleTimeString('ru-RU'),
          type: 'info',
          description: `Создан новый чат: ${response.chat.title}`
        });
        
        onChatCreated?.(response.chat);
        handleClose();
      }
    } catch (error) {
      console.error('Failed to create chat:', error);
      setError(error.message || t('chatCreationError', language));
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    setFormData({
      type: 'private',
      title: '',
      username: '',
      description: '',
      user_ids: []
    });
    setError('');
    setIsLoading(false);
    onClose();
  };

  const getChatTypeIcon = (type) => {
    switch (type) {
      case 'private':
        return <MessageCircle className="w-5 h-5 text-telegram-primary" />;
      case 'group':
        return <Users className="w-5 h-5 text-green-500" />;
      case 'channel':
        return <Hash className="w-5 h-5 text-blue-500" />;
      default:
        return <MessageCircle className="w-5 h-5" />;
    }
  };

  const getChatTypeLabel = (type) => {
    const language = getCurrentLanguage();
    switch (type) {
      case 'private':
        return t('privateChat', language);
      case 'group':
        return t('groupChat', language);
      case 'channel':
        return t('channelChat', language);
      default:
        return t('chat', language);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-telegram-sidebar rounded-lg shadow-xl max-w-lg w-full mx-4 max-h-[90vh] overflow-hidden">
        {/* Заголовок */}
        <div className="flex items-center justify-between p-4 border-b border-telegram-border">
          <h2 className="text-lg font-medium text-telegram-text">
            {t('createChat', getCurrentLanguage())}
          </h2>
          <button
            onClick={handleClose}
            className="p-1 text-telegram-secondary hover:text-telegram-text transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Форма */}
        <form onSubmit={handleSubmit} className="p-4 overflow-y-auto max-h-[calc(90vh-120px)]">
          {/* Тип чата */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              {t('chatType', getCurrentLanguage())}
            </label>
            <div className="grid grid-cols-3 gap-2">
              {['private', 'group', 'channel'].map((type) => (
                <button
                  key={type}
                  type="button"
                  onClick={() => setFormData(prev => ({ ...prev, type }))}
                  className={`flex items-center justify-center p-3 border rounded-lg transition-colors ${
                    formData.type === type
                      ? 'border-telegram-primary bg-telegram-primary/10'
                      : 'border-telegram-border hover:bg-telegram-primary/5'
                  }`}
                >
                  {getChatTypeIcon(type)}
                  <span className="ml-2 text-sm font-medium text-telegram-text">
                    {getChatTypeLabel(type)}
                  </span>
                </button>
              ))}
            </div>
          </div>

          {/* Название */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              {t('title', getCurrentLanguage())} *
            </label>
            <input
              type="text"
              name="title"
              value={formData.title}
              onChange={handleInputChange}
              placeholder={t('chatTitlePlaceholder', getCurrentLanguage())}
              className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
              required
            />
          </div>

          {/* Username (для групп и каналов) */}
          {formData.type !== 'private' && (
            <div className="mb-4">
              <label className="block text-sm font-medium text-telegram-text mb-2">
                Username
              </label>
              <input
                type="text"
                name="username"
                value={formData.username}
                onChange={handleInputChange}
                placeholder={t('usernameOptionalPlaceholder', getCurrentLanguage())}
                className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
              />
            </div>
          )}

          {/* Описание (для групп и каналов) */}
          {formData.type !== 'private' && (
            <div className="mb-4">
                          <label className="block text-sm font-medium text-telegram-text mb-2">
              {t('description', getCurrentLanguage())}
            </label>
              <textarea
                name="description"
                value={formData.description}
                onChange={handleInputChange}
                placeholder={t('chatDescriptionPlaceholder', getCurrentLanguage())}
                rows="3"
                className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary resize-none"
              />
            </div>
          )}

          {/* Выбор участников */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              {t('participants', getCurrentLanguage())} {formData.type === 'private' ? t('selectTwo', getCurrentLanguage()) : t('minimumTwo', getCurrentLanguage())}
            </label>
            <div className="max-h-48 overflow-y-auto border border-telegram-border rounded-lg">
              {users.map((user) => (
                <label
                  key={user.id}
                  className="flex items-center p-3 hover:bg-telegram-primary/5 cursor-pointer border-b border-telegram-border last:border-b-0"
                >
                  <input
                    type="checkbox"
                    checked={formData.user_ids.includes(user.id)}
                    onChange={() => handleUserToggle(user.id)}
                    className="mr-3"
                  />
                  <div className="w-8 h-8 rounded-full bg-telegram-primary flex items-center justify-center text-white font-medium mr-3">
                    {user.first_name?.charAt(0).toUpperCase() || 'U'}
                  </div>
                  <div className="flex-1">
                    <p className="text-sm font-medium text-telegram-text">
                      {user.first_name} {user.last_name || ''}
                    </p>
                    <p className="text-xs text-telegram-text-secondary">
                      @{user.username}
                    </p>
                  </div>
                  {user.is_bot && (
                    <span className="text-xs bg-blue-500 text-white px-2 py-1 rounded">
                      {t('isBot', getCurrentLanguage())}
                    </span>
                  )}
                </label>
              ))}
            </div>
            <p className="text-xs text-telegram-text-secondary mt-1">
              {t('selected', getCurrentLanguage())}: {formData.user_ids.length} {t('participants', getCurrentLanguage())}
            </p>
          </div>

          {/* Ошибка */}
          {error && (
            <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
              <p className="text-red-500 text-sm">{error}</p>
            </div>
          )}

          {/* Кнопки */}
          <div className="flex space-x-3">
            <button
              type="button"
              onClick={handleClose}
              className="flex-1 px-4 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text hover:bg-telegram-primary/10 transition-colors"
            >
              {t('cancel', getCurrentLanguage())}
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className="flex-1 px-4 py-2 bg-telegram-primary text-white rounded-lg hover:bg-telegram-primary/80 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? t('creating', getCurrentLanguage()) : t('create', getCurrentLanguage())}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default CreateChatModal;
