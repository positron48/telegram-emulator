import React, { useState, useEffect } from 'react';
import { X, Bot, Save, Copy } from 'lucide-react';
import apiService from '../services/api';
import useStore from '../store';
import { t, getCurrentLanguage } from '../locales';

const EditBotModal = ({ isOpen, onClose, bot, onBotUpdated }) => {
  const [formData, setFormData] = useState({
    name: '',
    username: '',
    token: '',
    webhook_url: '',
    is_active: true
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const { addDebugEvent } = useStore();

  useEffect(() => {
    if (bot) {
      setFormData({
        name: bot.name || '',
        username: bot.username || '',
        token: bot.token || '',
        webhook_url: bot.webhook_url || '',
        is_active: bot.is_active !== undefined ? bot.is_active : true
      });
    }
  }, [bot]);

  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
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
      const response = await apiService.updateBot(bot.id, formData);
      
      if (response.bot) {
        addDebugEvent({
          id: `bot-updated-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
          timestamp: new Date().toLocaleTimeString(getCurrentLanguage() === 'ru' ? 'ru-RU' : 'en-US'),
          type: 'info',
          description: `${t('botUpdated', getCurrentLanguage())}: ${response.bot.name}`
        });
        
        onBotUpdated?.(response.bot);
        handleClose();
      }
    } catch (error) {
      console.error('Failed to update bot:', error);
      setError(error.message || t('botUpdateError', getCurrentLanguage()));
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    setError('');
    setIsLoading(false);
    onClose();
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text).then(() => {
      addDebugEvent({
        id: `copy-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: new Date().toLocaleTimeString(getCurrentLanguage() === 'ru' ? 'ru-RU' : 'en-US'),
        type: 'info',
        description: t('copiedToClipboard', getCurrentLanguage())
      });
    });
  };

  if (!isOpen || !bot) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-telegram-sidebar rounded-lg shadow-xl max-w-md w-full mx-4">
        <div className="flex items-center justify-between p-4 border-b border-telegram-border">
          <h2 className="text-lg font-medium text-telegram-text">
            {t('editBot', getCurrentLanguage())}
          </h2>
          <button
            onClick={handleClose}
            className="p-1 text-telegram-secondary hover:text-telegram-text transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-4">
          <div className="mb-4">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              {t('title', getCurrentLanguage())} *
            </label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleInputChange}
              placeholder={t('botNamePlaceholder', getCurrentLanguage())}
              className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
              required
            />
          </div>

          <div className="mb-4">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              {t('username', getCurrentLanguage())} *
            </label>
            <input
              type="text"
              name="username"
              value={formData.username}
              onChange={handleInputChange}
              placeholder={t('usernamePlaceholder', getCurrentLanguage())}
              className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
              required
            />
          </div>

          <div className="mb-4">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              {t('token', getCurrentLanguage())}
            </label>
            <div className="relative">
              <input
                type="text"
                name="token"
                value={formData.token}
                onChange={handleInputChange}
                placeholder={t('tokenPlaceholder', getCurrentLanguage())}
                className="w-full px-3 py-2 pr-10 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
              />
              {formData.token && (
                <button
                  type="button"
                  onClick={() => copyToClipboard(formData.token)}
                  className="absolute right-2 top-1/2 transform -translate-y-1/2 p-1 text-telegram-secondary hover:text-telegram-text transition-colors"
                >
                  <Copy className="w-4 h-4" />
                </button>
              )}
            </div>
          </div>

          <div className="mb-4">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              {t('webhookUrl', getCurrentLanguage())}
            </label>
            <input
              type="url"
              name="webhook_url"
              value={formData.webhook_url}
              onChange={handleInputChange}
              placeholder={t('webhookUrlPlaceholder', getCurrentLanguage())}
              className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
            />
          </div>

          <div className="mb-6">
            <label className="flex items-center">
              <input
                type="checkbox"
                name="is_active"
                checked={formData.is_active}
                onChange={handleInputChange}
                className="mr-2 w-4 h-4 text-telegram-primary bg-telegram-bg border-telegram-border rounded focus:ring-telegram-primary focus:ring-2"
              />
              <span className="text-sm text-telegram-text">{t('botIsActive', getCurrentLanguage())}</span>
            </label>
          </div>

          {error && (
            <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
              <p className="text-red-500 text-sm">{error}</p>
            </div>
          )}

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
              {isLoading ? t('saving', getCurrentLanguage()) : t('save', getCurrentLanguage())}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default EditBotModal;
