import React, { useState } from 'react';
import { X, Bot } from 'lucide-react';
import apiService from '../services/api';
import useStore from '../store';
import { t, getCurrentLanguage } from '../locales';

const CreateBotModal = ({ isOpen, onClose, onBotCreated }) => {
  const [formData, setFormData] = useState({
    name: '',
    username: '',
    token: '',
    webhook_url: ''
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const { addDebugEvent } = useStore();

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.name.trim() || !formData.username.trim()) {
      setError(t('nameAndUsernameRequired', getCurrentLanguage()));
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const response = await apiService.createBot(formData);
      
      if (response.bot) {
        addDebugEvent({
          id: `bot-created-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
          timestamp: new Date().toLocaleTimeString('ru-RU'),
          type: 'info',
          description: `Создан новый бот: ${response.bot.name}`
        });
        
        onBotCreated?.(response.bot);
        handleClose();
      }
    } catch (error) {
      console.error('Failed to create bot:', error);
      setError(error.message || t('botCreationError', getCurrentLanguage()));
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    setFormData({
      name: '',
      username: '',
      token: '',
      webhook_url: ''
    });
    setError('');
    setIsLoading(false);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-telegram-sidebar rounded-lg shadow-xl max-w-md w-full mx-4">
        {/* Заголовок */}
        <div className="flex items-center justify-between p-4 border-b border-telegram-border">
          <h2 className="text-lg font-medium text-telegram-text">
            {t('createBot', getCurrentLanguage())}
          </h2>
          <button
            onClick={handleClose}
            className="p-1 text-telegram-secondary hover:text-telegram-text transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Форма */}
        <form onSubmit={handleSubmit} className="p-4">
          {/* Название */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              Название *
            </label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleInputChange}
              placeholder="Название бота"
              className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
              required
            />
          </div>

          {/* Username */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              Username *
            </label>
            <input
              type="text"
              name="username"
              value={formData.username}
              onChange={handleInputChange}
              placeholder="username"
              className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
              required
            />
          </div>

          {/* Токен */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              Токен
            </label>
            <input
              type="text"
              name="token"
              value={formData.token}
              onChange={handleInputChange}
              placeholder="1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"
              className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
            />
          </div>

          {/* Webhook URL */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              Webhook URL
            </label>
            <input
              type="url"
              name="webhook_url"
              value={formData.webhook_url}
              onChange={handleInputChange}
              placeholder="https://example.com/webhook"
              className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
            />
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
              Отмена
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className="flex-1 px-4 py-2 bg-telegram-primary text-white rounded-lg hover:bg-telegram-primary/80 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? 'Создание...' : 'Создать'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default CreateBotModal;
