import React, { useState } from 'react';
import { X, Trash2, Download, RefreshCw, AlertCircle, Info, MessageCircle, Activity } from 'lucide-react';
import clsx from 'clsx';
import { t, getCurrentLanguage } from '../locales';

const DebugPanel = ({ events, statistics, onClose }) => {
  const [activeTab, setActiveTab] = useState('events');

  const getEventIcon = (type) => {
    switch (type) {
      case 'error':
        return <AlertCircle className="w-4 h-4 text-red-500" />;
      case 'message':
        return <MessageCircle className="w-4 h-4 text-blue-500" />;
      case 'api_call':
        return <Activity className="w-4 h-4 text-green-500" />;
      default:
        return <Info className="w-4 h-4 text-telegram-secondary" />;
    }
  };

  const getEventColor = (type) => {
    switch (type) {
      case 'error':
        return 'text-red-500';
      case 'message':
        return 'text-blue-500';
      case 'api_call':
        return 'text-green-500';
      default:
        return 'text-telegram-text-secondary';
    }
  };

  const exportLogs = () => {
    const logData = {
      timestamp: new Date().toISOString(),
      events: events,
      statistics: statistics
    };

    const blob = new Blob([JSON.stringify(logData, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `telegram-emulator-logs-${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const clearEvents = () => {
    // Здесь нужно добавить действие для очистки событий
    console.log('Clear events');
  };

  return (
    <div className="w-96 bg-telegram-sidebar border-l border-telegram-border flex flex-col">
      {/* Заголовок */}
      <div className="p-4 border-b border-telegram-border bg-telegram-bg flex items-center justify-between">
        <h2 className="text-lg font-medium text-telegram-text">{t('debug', getCurrentLanguage())}</h2>
        <button
          onClick={onClose}
          className="p-1 text-telegram-secondary hover:text-telegram-text transition-colors"
        >
          <X className="w-5 h-5" />
        </button>
      </div>

      {/* Табы */}
      <div className="flex border-b border-telegram-border">
        <button
          onClick={() => setActiveTab('events')}
          className={clsx(
            'flex-1 px-4 py-2 text-sm font-medium transition-colors',
            activeTab === 'events'
              ? 'text-telegram-text border-b-2 border-telegram-primary'
              : 'text-telegram-text-secondary hover:text-telegram-text'
          )}
        >
          {t('events', getCurrentLanguage())}
        </button>
        <button
          onClick={() => setActiveTab('statistics')}
          className={clsx(
            'flex-1 px-4 py-2 text-sm font-medium transition-colors',
            activeTab === 'statistics'
              ? 'text-telegram-text border-b-2 border-telegram-primary'
              : 'text-telegram-text-secondary hover:text-telegram-text'
          )}
        >
          {t('statistics', getCurrentLanguage())}
        </button>
      </div>

      {/* Контент */}
      <div className="flex-1 overflow-hidden">
        {activeTab === 'events' && (
          <div className="h-full flex flex-col">
            {/* Панель действий */}
            <div className="p-3 border-b border-telegram-border flex space-x-2">
              <button
                onClick={clearEvents}
                className="flex items-center px-3 py-1 bg-telegram-bg border border-telegram-border rounded text-telegram-text hover:bg-telegram-primary hover:border-telegram-primary transition-colors"
              >
                <Trash2 className="w-4 h-4 mr-1" />
                Очистить
              </button>
              <button
                onClick={exportLogs}
                className="flex items-center px-3 py-1 bg-telegram-bg border border-telegram-border rounded text-telegram-text hover:bg-telegram-primary hover:border-telegram-primary transition-colors"
              >
                <Download className="w-4 h-4 mr-1" />
                Экспорт
              </button>
            </div>

            {/* Список событий */}
            <div className="flex-1 overflow-y-auto p-3 space-y-2">
              {events.length === 0 ? (
                <div className="text-center py-8">
                  <Info className="w-8 h-8 text-telegram-secondary mx-auto mb-2" />
                  <p className="text-telegram-text-secondary text-sm">
                    Нет событий для отображения
                  </p>
                </div>
              ) : (
                events.map((event, index) => (
                  <div
                    key={`${event.id}-${index}`}
                    className="p-3 bg-telegram-bg rounded-lg border border-telegram-border"
                  >
                    <div className="flex items-start space-x-2">
                      {getEventIcon(event.type)}
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between mb-1">
                          <span className={clsx('text-sm font-medium', getEventColor(event.type))}>
                            {event.type.toUpperCase()}
                          </span>
                          <span className="text-xs text-telegram-text-secondary">
                            {event.timestamp}
                          </span>
                        </div>
                        <p className="text-sm text-telegram-text">
                          {event.description}
                        </p>
                        {event.data && (
                          <details className="mt-2">
                            <summary className="text-xs text-telegram-text-secondary cursor-pointer">
                              Детали
                            </summary>
                            <pre className="text-xs text-telegram-text-secondary mt-1 p-2 bg-telegram-sidebar rounded overflow-x-auto">
                              {JSON.stringify(event.data, null, 2)}
                            </pre>
                          </details>
                        )}
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        )}

        {activeTab === 'statistics' && (
          <div className="p-4 space-y-4">
            {/* Основная статистика */}
            <div className="grid grid-cols-2 gap-4">
              <div className="p-3 bg-telegram-bg rounded-lg border border-telegram-border">
                <div className="text-2xl font-bold text-telegram-text">
                  {statistics.messages_count || 0}
                </div>
                <div className="text-xs text-telegram-text-secondary">
                  Сообщений
                </div>
              </div>
              
              <div className="p-3 bg-telegram-bg rounded-lg border border-telegram-border">
                <div className="text-2xl font-bold text-telegram-text">
                  {statistics.response_time || 0}ms
                </div>
                <div className="text-xs text-telegram-text-secondary">
                  Время ответа
                </div>
              </div>
              
              <div className="p-3 bg-telegram-bg rounded-lg border border-telegram-border">
                <div className="text-2xl font-bold text-red-500">
                  {statistics.errors_count || 0}
                </div>
                <div className="text-xs text-telegram-text-secondary">
                  Ошибок
                </div>
              </div>
              
              <div className="p-3 bg-telegram-bg rounded-lg border border-telegram-border">
                <div className="text-2xl font-bold text-telegram-text">
                  {statistics.users_count || 0}
                </div>
                <div className="text-xs text-telegram-text-secondary">
                  Пользователей
                </div>
              </div>
            </div>

            {/* Дополнительная информация */}
            <div className="p-3 bg-telegram-bg rounded-lg border border-telegram-border">
              <h3 className="text-sm font-medium text-telegram-text mb-2">
                Системная информация
              </h3>
              <div className="space-y-1 text-xs text-telegram-text-secondary">
                <div className="flex justify-between">
                  <span>Версия:</span>
                  <span>1.0.0</span>
                </div>
                <div className="flex justify-between">
                  <span>Время работы:</span>
                  <span>00:15:30</span>
                </div>
                <div className="flex justify-between">
                  <span>Память:</span>
                  <span>45.2 MB</span>
                </div>
                <div className="flex justify-between">
                  <span>CPU:</span>
                  <span>2.3%</span>
                </div>
              </div>
            </div>

            {/* Действия */}
            <div className="space-y-2">
              <button className="w-full flex items-center justify-center px-3 py-2 bg-telegram-bg border border-telegram-border rounded text-telegram-text hover:bg-telegram-primary hover:border-telegram-primary transition-colors">
                <RefreshCw className="w-4 h-4 mr-2" />
                Обновить статистику
              </button>
              
              <button className="w-full flex items-center justify-center px-3 py-2 bg-red-600 border border-red-600 rounded text-white hover:bg-red-700 transition-colors">
                <Trash2 className="w-4 h-4 mr-2" />
                Сбросить все данные
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default DebugPanel;
