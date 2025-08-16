import React, { useState } from 'react';
import { X, Save, RefreshCw, Download, Upload, Moon, Sun, Monitor } from 'lucide-react';
import useStore from '../store';

const SettingsPanel = ({ isOpen, onClose }) => {
  const [settings, setSettings] = useState({
    theme: 'system',
    language: 'ru',
    notifications: true,
    sound: true,
    autoScroll: true,
    debugMode: false
  });

  const { addDebugEvent } = useStore();

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

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-telegram-sidebar rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[80vh] overflow-hidden">
        {/* Заголовок */}
        <div className="flex items-center justify-between p-4 border-b border-telegram-border">
          <h2 className="text-lg font-medium text-telegram-text">
            Настройки
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
              <h3 className="text-md font-medium text-telegram-text mb-3">Тема оформления</h3>
              <div className="grid grid-cols-3 gap-3">
                {[
                  { value: 'light', label: 'Светлая', icon: Sun },
                  { value: 'dark', label: 'Темная', icon: Moon },
                  { value: 'system', label: 'Системная', icon: Monitor }
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
                Язык интерфейса
              </label>
              <select
                value={settings.language}
                onChange={(e) => handleSettingChange('language', e.target.value)}
                className="w-full px-3 py-2 bg-telegram-bg border border-telegram-border rounded-lg text-telegram-text focus:outline-none focus:border-telegram-primary"
              >
                <option value="ru">Русский</option>
                <option value="en">English</option>
              </select>
            </div>

            {/* Уведомления */}
            <div>
              <h3 className="text-md font-medium text-telegram-text mb-3">Уведомления</h3>
              <div className="space-y-3">
                <label className="flex items-center">
                  <input
                    type="checkbox"
                    checked={settings.notifications}
                    onChange={(e) => handleSettingChange('notifications', e.target.checked)}
                    className="mr-3"
                  />
                  <span className="text-sm text-telegram-text">Включить уведомления</span>
                </label>
                <label className="flex items-center">
                  <input
                    type="checkbox"
                    checked={settings.sound}
                    onChange={(e) => handleSettingChange('sound', e.target.checked)}
                    className="mr-3"
                  />
                  <span className="text-sm text-telegram-text">Звуковые уведомления</span>
                </label>
              </div>
            </div>

            {/* Дополнительно */}
            <div>
              <h3 className="text-md font-medium text-telegram-text mb-3">Дополнительно</h3>
              <div className="space-y-3">
                <label className="flex items-center">
                  <input
                    type="checkbox"
                    checked={settings.autoScroll}
                    onChange={(e) => handleSettingChange('autoScroll', e.target.checked)}
                    className="mr-3"
                  />
                  <span className="text-sm text-telegram-text">Автоматическая прокрутка к новым сообщениям</span>
                </label>
                <label className="flex items-center">
                  <input
                    type="checkbox"
                    checked={settings.debugMode}
                    onChange={(e) => handleSettingChange('debugMode', e.target.checked)}
                    className="mr-3"
                  />
                  <span className="text-sm text-telegram-text">Режим отладки</span>
                </label>
              </div>
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
              Отмена
            </button>
            <button
              onClick={onClose}
              className="flex items-center px-4 py-2 bg-telegram-primary text-white rounded-lg hover:bg-telegram-primary/80 transition-colors"
            >
              <Save className="w-4 h-4 mr-2" />
              Сохранить
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SettingsPanel;
