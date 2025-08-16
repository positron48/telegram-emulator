import React from 'react';
import { X, Plus, Trash2 } from 'lucide-react';
import clsx from 'clsx';

const UserModal = ({ isOpen, onClose, users, currentUser, onUserSelect, onCreateUser, onDeleteUser }) => {
  if (!isOpen) return null;

  const handleUserSelect = (user) => {
    onUserSelect(user);
    onClose();
  };

  const handleDeleteUser = (e, userId) => {
    e.stopPropagation();
    if (confirm('Вы уверены, что хотите удалить этого пользователя?')) {
      onDeleteUser(userId);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-telegram-sidebar rounded-lg shadow-xl max-w-md w-full mx-4 max-h-[80vh] overflow-hidden">
        {/* Заголовок */}
        <div className="flex items-center justify-between p-4 border-b border-telegram-border">
          <h2 className="text-lg font-medium text-telegram-text">
            Выберите пользователя
          </h2>
          <button
            onClick={onClose}
            className="p-1 text-telegram-secondary hover:text-telegram-text transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Список пользователей */}
        <div className="overflow-y-auto max-h-[60vh]">
          {users.map((user) => (
            <button
              key={user.id}
              onClick={() => handleUserSelect(user)}
              className={clsx(
                'flex items-center w-full p-4 hover:bg-telegram-primary transition-colors border-b border-telegram-border last:border-b-0',
                currentUser?.id === user.id && 'bg-telegram-primary'
              )}
            >
              <div className="w-10 h-10 rounded-full bg-gradient-to-br from-telegram-primary to-blue-600 flex items-center justify-center text-white font-medium mr-3 shadow-sm">
                {user.first_name?.charAt(0).toUpperCase() || 'U'}
              </div>
              <div className="flex-1 text-left">
                <p className="text-sm font-medium text-telegram-text">
                  {user.first_name} {user.last_name || ''}
                </p>
                <p className="text-xs text-telegram-text-secondary">
                  @{user.username}
                </p>
              </div>
              {user.is_bot && (
                <span className="text-xs bg-blue-500 text-white px-2 py-1 rounded">
                  Бот
                </span>
              )}
              <button
                onClick={(e) => handleDeleteUser(e, user.id)}
                className="p-1 text-red-500 hover:text-red-600 transition-colors"
                title="Удалить пользователя"
              >
                <Trash2 className="w-4 h-4" />
              </button>
            </button>
          ))}
        </div>

        {/* Кнопка создания нового пользователя */}
        <div className="p-4 border-t border-telegram-border">
          <button
            onClick={() => {
              onClose();
              onCreateUser();
            }}
            className="flex items-center justify-center w-full p-3 bg-telegram-primary text-white rounded-lg hover:bg-telegram-primary/80 transition-colors"
          >
            <Plus className="w-4 h-4 mr-2" />
            <span>Создать нового пользователя</span>
          </button>
        </div>
      </div>
    </div>
  );
};

export default UserModal;
