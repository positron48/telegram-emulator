import React, { useState } from 'react';
import { Settings } from 'lucide-react';
import UserModal from './UserModal';
import { t, getCurrentLanguage } from '../locales';

const UserSelector = ({ users, currentUser, onUserSelect, onCreateUser, onDeleteUser }) => {
  const [isModalOpen, setIsModalOpen] = useState(false);

  return (
    <>
      {/* Кнопка выбора пользователя */}
      <button
        onClick={() => setIsModalOpen(true)}
        className="flex items-center p-3 bg-telegram-sidebar rounded-lg text-telegram-text hover:bg-telegram-primary transition-colors w-full"
      >
        <div className="w-8 h-8 rounded-full bg-gradient-to-br from-telegram-primary to-blue-600 flex items-center justify-center text-white font-medium mr-3 shadow-sm">
          {currentUser?.first_name?.charAt(0).toUpperCase() || 'U'}
        </div>
        <div className="flex-1 text-left">
          <p className="text-sm font-medium truncate">
            {currentUser ? `${currentUser.first_name} ${currentUser.last_name || ''}`.trim() : t('selectUser', getCurrentLanguage())}
          </p>
          <p className="text-xs text-telegram-text-secondary truncate">
            @{currentUser?.username || 'username'}
          </p>
        </div>
        <Settings className="w-4 h-4 text-telegram-secondary" />
      </button>

      {/* Модальное окно выбора пользователя */}
      <UserModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        users={users}
        currentUser={currentUser}
        onUserSelect={onUserSelect}
        onCreateUser={onCreateUser}
        onDeleteUser={onDeleteUser}
      />
    </>
  );
};

export default UserSelector;


