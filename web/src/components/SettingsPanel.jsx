import React, { useState, useEffect } from 'react';
import { X, Save, Moon, Sun } from 'lucide-react';
import useStore from '../store';
import { t, setLanguage, getCurrentLanguage } from '../locales';

const SettingsPanel = ({ isOpen, onClose }) => {
  const [settings, setSettings] = useState({
    theme: 'light',
    language: 'ru'
  });

  const { addDebugEvent } = useStore();

  // Загружаем настройки из localStorage при открытии
  useEffect(() => {
    if (isOpen) {
      const savedSettings = localStorage.getItem('telegram-emulator-settings');
      if (savedSettings) {
        try {
          const parsed = JSON.parse(savedSettings);
          setSettings(prev => ({ ...prev, ...parsed }));
        } catch (error) {
          console.error('Failed to parse saved settings:', error);
        }
      }
      
      // Загружаем язык из localStorage
      const currentLanguage = getCurrentLanguage();
      setSettings(prev => ({ ...prev, language: currentLanguage }));
    }
  }, [isOpen]);

  const handleSettingChange = (key, value) => {
    setSettings(prev => ({
      ...prev,
      [key]: value
    }));

    addDebugEvent({
      id: `setting-change-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      timestamp: new Date().toLocaleTimeString('ru-RU'),
      type: 'info',
      description: `Настройка изменена: ${key} = ${value}`
    });
  };

  const handleSaveSettings = () => {
    try {
      localStorage.setItem('telegram-emulator-settings', JSON.stringify(settings));
      
      addDebugEvent({
        id: `settings-saved-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: new Date().toLocaleTimeString('ru-RU'),
        type: 'success',
        description: 'Настройки сохранены'
      });

      // Применяем настройки
      applySettings(settings);
      
      onClose();
    } catch (error) {
      console.error('Failed to save settings:', error);
      addDebugEvent({
        id: `settings-error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: new Date().toLocaleTimeString('ru-RU'),
        type: 'error',
        description: `Ошибка сохранения настроек: ${error.message}`
      });
    }
  };

  const applySettings = (newSettings) => {
    // Применяем тему
    if (newSettings.theme === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }

    // Применяем язык
    setLanguage(newSettings.language);
  };



  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-telegram-sidebar rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[80vh] overflow-hidden">
        {/* Заголовок */}
        <div className="flex items-center justify-between p-4 border-b border-telegram-border">
          <h2 className="text-lg font-medium text-telegram-text">
            {t('settings', settings.language)}
          </h2>
          <button
            onClick={onClose}
            className="p-1 text-telegram-secondary hover:text-telegram-text transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Содержимое */}
        <div className="p-6 overflow-y-auto max-h-[calc(80vh-120px)]">
          <div className="space-y-6">
            {/* Тема */}
            <div>
              <h3 className="text-md font-medium text-telegram-text mb-3">{t('theme', settings.language)}</h3>
              <div className="grid grid-cols-2 gap-3">
                {[
                  { value: 'light', label: t('lightTheme', settings.language), icon: Sun },
                  { value: 'dark', label: t('darkTheme', settings.language), icon: Moon }
                ].map((theme) => {
                  const Icon = theme.icon;
                  return (
                    <button
                      key={theme.value}
                      onClick={() => handleSettingChange('theme', theme.value)}
                      className={`flex flex-col items-center p-4 border rounded-lg transition-colors ${
                        settings.theme === theme.value
                          ? 'border-telegram-primary bg-telegram-primary/10'
                          : 'border-telegram-border hover:bg-telegram-primary/5'
                      }`}
                    >
                      <Icon className="w-6 h-6 mb-2 text-telegram-primary" />
                      <span className="text-sm text-telegram-text">{theme.label}</span>
                    </button>
                  );
                })}
              </div>
            </div>

            {/* Язык */}
            <div>
              <label className="block text-sm font-medium text-telegram-text mb-2">
                {t('language', settings.language)}
              </label>
              <select
                value={settings.language}
                onChange={(e) => handleSettingChange('language', e.target.value)}
                className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text focus:outline-none focus:border-telegram-primary"
              >
                <option value="ru">{t('russian', settings.language)}</option>
                <option value="en">{t('english', settings.language)}</option>
              </select>
            </div>
          </div>
        </div>

        {/* Нижняя панель */}
        <div className="flex items-center justify-between p-4 border-t border-telegram-border">
          <div className="text-sm text-telegram-text-secondary">
            Telegram Emulator v1.0.0
          </div>
          <div className="flex space-x-3">
            <button
              onClick={onClose}
              className="px-4 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text hover:bg-telegram-primary/10 transition-colors"
            >
              {t('cancel', settings.language)}
            </button>
            <button
              onClick={handleSaveSettings}
              className="flex items-center px-4 py-2 bg-telegram-primary text-white rounded-lg hover:bg-telegram-primary/80 transition-colors"
            >
              <Save className="w-4 h-4 mr-2" />
              {t('save', settings.language)}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SettingsPanel;
