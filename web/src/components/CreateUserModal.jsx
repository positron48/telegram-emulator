import React, { useState } from 'react';
import { X, User, Bot } from 'lucide-react';
import apiService from '../services/api';
import useStore from '../store';
import { t, getCurrentLanguage } from '../locales';

const CreateUserModal = ({ isOpen, onClose, onUserCreated }) => {
  const [formData, setFormData] = useState({
    username: '',
    first_name: '',
    last_name: '',
    is_bot: false
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const { addUser, addDebugEvent } = useStore();

  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.username.trim() || !formData.first_name.trim()) {
      setError(t('usernameAndFirstNameRequired', getCurrentLanguage()));
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const response = await apiService.createUser(formData);
      
      if (response.user) {
        addUser(response.user);
        addDebugEvent({
          id: `user-created-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
          timestamp: new Date().toLocaleTimeString(getCurrentLanguage() === 'ru' ? 'ru-RU' : 'en-US'),
          type: 'info',
          description: `${t('userCreated', getCurrentLanguage())}: ${response.user.username}`
        });
        
        onUserCreated?.(response.user);
        handleClose();
      }
    } catch (error) {
      console.error('Failed to create user:', error);
              setError(error.message || t('userCreationError', getCurrentLanguage()));
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    setFormData({
      username: '',
      first_name: '',
      last_name: '',
      is_bot: false
    });
    setError('');
    setIsLoading(false);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-telegram-sidebar rounded-lg shadow-xl max-w-md w-full mx-4">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-telegram-border">
          <h2 className="text-lg font-medium text-telegram-text">
            {t('createUser', getCurrentLanguage())}
          </h2>
          <button
            onClick={handleClose}
            className="p-1 text-telegram-secondary hover:text-telegram-text transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-4">
          {/* User Type */}
          <div className="mb-4">
            <label className="flex items-center cursor-pointer">
              <input
                type="checkbox"
                name="is_bot"
                checked={formData.is_bot}
                onChange={handleInputChange}
                className="sr-only"
              />
              <div className="flex items-center p-3 border border-telegram-border rounded-lg hover:bg-telegram-primary/10 transition-colors">
                {formData.is_bot ? (
                  <Bot className="w-5 h-5 text-blue-500 mr-2" />
                ) : (
                  <User className="w-5 h-5 text-telegram-primary mr-2" />
                )}
                <span className="text-telegram-text font-medium">
                  {formData.is_bot ? t('isBot', getCurrentLanguage()) : t('user', getCurrentLanguage())}
                </span>
              </div>
            </label>
          </div>

          {/* Username */}
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

          {/* First Name */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              {t('firstName', getCurrentLanguage())} *
            </label>
            <input
              type="text"
              name="first_name"
              value={formData.first_name}
              onChange={handleInputChange}
              placeholder={t('firstNamePlaceholder', getCurrentLanguage())}
              className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
              required
            />
          </div>

          {/* Last Name */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-telegram-text mb-2">
              {t('lastName', getCurrentLanguage())}
            </label>
            <input
              type="text"
              name="last_name"
              value={formData.last_name}
              onChange={handleInputChange}
              placeholder={t('lastNamePlaceholder', getCurrentLanguage())}
              className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text placeholder-telegram-secondary focus:outline-none focus:border-telegram-primary"
            />
          </div>

          {/* Error */}
          {error && (
            <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
              <p className="text-red-500 text-sm">{error}</p>
            </div>
          )}

          {/* Buttons */}
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

export default CreateUserModal;
