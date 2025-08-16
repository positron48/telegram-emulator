import React, { useState, useEffect } from 'react';
import { X, Bot, Plus, Settings, Play, Pause, Trash2, Copy } from 'lucide-react';
import apiService from '../services/api';
import useStore from '../store';
import CreateBotModal from './CreateBotModal';

const BotManager = ({ isOpen, onClose }) => {
  const [bots, setBots] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);

  const { addDebugEvent } = useStore();

  useEffect(() => {
    if (isOpen) {
      loadBots();
    }
  }, [isOpen]);

  const loadBots = async () => {
    setIsLoading(true);
    try {
      const response = await apiService.getBots();
      setBots(response.bots || []);
    } catch (error) {
      console.error('Failed to load bots:', error);
      setError('Ошибка загрузки ботов');
    } finally {
      setIsLoading(false);
    }
  };

  const handleToggleBot = async (botId, isActive) => {
    try {
      const response = await apiService.updateBot(botId, { is_active: !isActive });
      if (response.bot) {
        setBots(prev => prev.map(bot => 
          bot.id === botId ? { ...bot, is_active: response.bot.is_active } : bot
        ));
        addDebugEvent({
          id: `bot-toggle-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
          timestamp: new Date().toLocaleTimeString('ru-RU'),
          type: 'info',
          description: `Бот ${response.bot.name} ${response.bot.is_active ? 'активирован' : 'деактивирован'}`
        });
      }
    } catch (error) {
      console.error('Failed to toggle bot:', error);
      setError('Ошибка изменения статуса бота');
    }
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    addDebugEvent({
      id: `copy-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      timestamp: new Date().toLocaleTimeString('ru-RU'),
      type: 'info',
      description: 'Скопировано в буфер обмена'
    });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-telegram-sidebar rounded-lg shadow-xl max-w-4xl w-full mx-4 max-h-[90vh] overflow-hidden">
        {/* Заголовок */}
        <div className="flex items-center justify-between p-4 border-b border-telegram-border">
          <h2 className="text-lg font-medium text-telegram-text">
            Управление ботами
          </h2>
          <button
            onClick={onClose}
            className="p-1 text-telegram-secondary hover:text-telegram-text transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Содержимое */}
        <div className="p-4 overflow-y-auto max-h-[calc(90vh-120px)]">
          {/* Ошибка */}
          {error && (
            <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
              <p className="text-red-500 text-sm">{error}</p>
            </div>
          )}

          {/* Кнопка создания */}
          <div className="mb-4">
            <button
              onClick={() => setShowCreateModal(true)}
              className="flex items-center px-4 py-2 bg-telegram-primary text-white rounded-lg hover:bg-telegram-primary/80 transition-colors"
            >
              <Plus className="w-4 h-4 mr-2" />
              Создать бота
            </button>
          </div>

          {/* Список ботов */}
          {isLoading ? (
            <div className="text-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-telegram-primary mx-auto mb-4"></div>
              <p className="text-telegram-text-secondary">Загрузка ботов...</p>
            </div>
          ) : bots.length === 0 ? (
            <div className="text-center py-8">
              <Bot className="w-12 h-12 text-telegram-secondary mx-auto mb-3" />
              <h3 className="text-telegram-text font-medium mb-1">Нет ботов</h3>
              <p className="text-telegram-text-secondary text-sm">
                Создайте первого бота для начала работы
              </p>
            </div>
          ) : (
            <div className="grid gap-4">
              {bots.map((bot) => (
                <div
                  key={bot.id}
                  className="p-4 border border-telegram-border rounded-lg hover:bg-telegram-primary/5 transition-colors"
                >
                  <div className="flex items-center justify-between mb-3">
                    <div className="flex items-center">
                      <Bot className="w-6 h-6 text-blue-500 mr-3" />
                      <div>
                        <h3 className="text-telegram-text font-medium">{bot.name}</h3>
                        <p className="text-sm text-telegram-text-secondary">@{bot.username}</p>
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <button
                        onClick={() => handleToggleBot(bot.id, bot.is_active)}
                        className={`p-2 rounded-lg transition-colors ${
                          bot.is_active
                            ? 'bg-green-500/10 text-green-500 hover:bg-green-500/20'
                            : 'bg-gray-500/10 text-gray-500 hover:bg-gray-500/20'
                        }`}
                        title={bot.is_active ? 'Деактивировать' : 'Активировать'}
                      >
                        {bot.is_active ? <Pause className="w-4 h-4" /> : <Play className="w-4 h-4" />}
                      </button>
                    </div>
                  </div>

                  {/* Информация о боте */}
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                    <div>
                      <p className="text-telegram-text-secondary mb-1">Токен:</p>
                      <div className="flex items-center">
                        <code className="flex-1 bg-telegram-bg px-2 py-1 rounded text-xs font-mono text-telegram-text">
                          {bot.token}
                        </code>
                        <button
                          onClick={() => copyToClipboard(bot.token)}
                          className="ml-2 p-1 text-telegram-secondary hover:text-telegram-text transition-colors"
                          title="Копировать токен"
                        >
                          <Copy className="w-3 h-3" />
                        </button>
                      </div>
                    </div>

                    {bot.webhook_url && (
                      <div>
                        <p className="text-telegram-text-secondary mb-1">Webhook URL:</p>
                        <div className="flex items-center">
                          <code className="flex-1 bg-telegram-bg px-2 py-1 rounded text-xs font-mono text-telegram-text">
                            {bot.webhook_url}
                          </code>
                          <button
                            onClick={() => copyToClipboard(bot.webhook_url)}
                            className="ml-2 p-1 text-telegram-secondary hover:text-telegram-text transition-colors"
                            title="Копировать URL"
                          >
                            <Copy className="w-3 h-3" />
                          </button>
                        </div>
                      </div>
                    )}
                  </div>

                  {/* Статус */}
                  <div className="mt-3 flex items-center justify-between">
                    <div className="flex items-center">
                      <div className={`w-2 h-2 rounded-full mr-2 ${
                        bot.is_active ? 'bg-green-500' : 'bg-gray-500'
                      }`}></div>
                      <span className={`text-xs ${
                        bot.is_active ? 'text-green-500' : 'text-gray-500'
                      }`}>
                        {bot.is_active ? 'Активен' : 'Неактивен'}
                      </span>
                    </div>
                    <span className="text-xs text-telegram-text-secondary">
                      Создан: {new Date(bot.created_at).toLocaleDateString('ru-RU')}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Модальное окно создания бота */}
        <CreateBotModal
          isOpen={showCreateModal}
          onClose={() => setShowCreateModal(false)}
          onBotCreated={(bot) => {
            setBots(prev => [...prev, bot]);
            setShowCreateModal(false);
          }}
        />
      </div>
    </div>
  );
};

export default BotManager;
