import React, { useState, useEffect } from 'react';
import { X, UserPlus, UserMinus, Users } from 'lucide-react';
import apiService from '../services/api';
import useStore from '../store';
import { t, getCurrentLanguage } from '../locales';

const ChatMembersModal = ({ isOpen, onClose, chat }) => {
  const [members, setMembers] = useState([]);
  const [availableUsers, setAvailableUsers] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isAddingMember, setIsAddingMember] = useState(false);
  const [isRemovingMember, setIsRemovingMember] = useState(false);
  const [error, setError] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  const { users, addDebugEvent } = useStore();

  useEffect(() => {
    if (isOpen && chat) {
      loadMembers();
    }
  }, [isOpen, chat]);

  // Update available users list when members change
  useEffect(() => {
    if (isOpen && chat) {
      loadAvailableUsers();
    }
  }, [isOpen, chat, members, users]);

  const loadMembers = async () => {
    if (!chat) return;
    
    setIsLoading(true);
    setError('');
    
    try {
      // Get chat members
      const response = await apiService.getChatMembers(chat.id);
      setMembers(response.members || []);
    } catch (error) {
      // console.error('Failed to load chat members:', error);
      setError(t('failedToLoadMembers', getCurrentLanguage()));
    } finally {
      setIsLoading(false);
    }
  };

  const loadAvailableUsers = async () => {
    try {
      // Get all users who are not in the chat
      // Use current members list from state
      const chatMemberIds = members.map(m => m.id);
      const available = users.filter(user => !chatMemberIds.includes(user.id));
      setAvailableUsers(available);
    } catch (error) {
      // console.error('Failed to load available users:', error);
    }
  };

  const handleAddMember = async (userId) => {
    if (!chat) return;
    
    setIsAddingMember(true);
    setError('');
    setSuccessMessage('');
    
    try {
      await apiService.addChatMember(chat.id, userId);
      
      addDebugEvent({
        id: `member-added-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: new Date().toLocaleTimeString(getCurrentLanguage() === 'en' ? 'en-US' : 'ru-RU'),
        type: 'info',
        description: t('memberAdded', getCurrentLanguage())
      });
      
      setSuccessMessage(t('memberAddedSuccess', getCurrentLanguage()));
      
      // Reload data
      await loadMembers();
      
      // Clear notification after 3 seconds
      setTimeout(() => setSuccessMessage(''), 3000);
    } catch (error) {
      // console.error('Failed to add member:', error);
      setError(t('failedToAddMember', getCurrentLanguage()));
    } finally {
      setIsAddingMember(false);
    }
  };

  const handleRemoveMember = async (userId) => {
    if (!chat) return;
    
    const language = getCurrentLanguage();
    if (!confirm(t('confirmRemoveMember', language))) {
      return;
    }
    
    setIsRemovingMember(true);
    setError('');
    setSuccessMessage('');
    
    try {
      await apiService.removeChatMember(chat.id, userId);
      
      addDebugEvent({
        id: `member-removed-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: new Date().toLocaleTimeString(language === 'en' ? 'en-US' : 'ru-RU'),
        type: 'warning',
        description: t('memberRemoved', language)
      });
      
      setSuccessMessage(t('memberRemovedSuccess', language));
      
      // Reload data
      await loadMembers();
      
      // Clear notification after 3 seconds
      setTimeout(() => setSuccessMessage(''), 3000);
    } catch (error) {
      // console.error('Failed to remove member:', error);
      setError(t('failedToRemoveMember', language));
    } finally {
      setIsRemovingMember(false);
    }
  };

  const getMemberRole = (member) => {
    if (chat?.type === 'private') {
      return t('participant', getCurrentLanguage());
    }
    
    // For groups and channels, role logic can be added
    if (member.id === chat?.created_by) {
      return t('creator', getCurrentLanguage());
    }
    
    return t('member', getCurrentLanguage());
  };

  if (!isOpen || !chat) return null;

  const language = getCurrentLanguage();

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-telegram-sidebar rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-telegram-border">
          <div className="flex items-center">
            <Users className="w-5 h-5 text-telegram-primary mr-2" />
            <h2 className="text-lg font-medium text-telegram-text">
              {t('chatMembers', language)} - {chat.title}
            </h2>
          </div>
          <button
            onClick={onClose}
            className="p-1 text-telegram-secondary hover:text-telegram-text transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-4 overflow-y-auto max-h-[calc(90vh-120px)]">
          {/* Error */}
          {error && (
            <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
              <p className="text-red-500 text-sm">{error}</p>
            </div>
          )}

          {/* Success */}
          {successMessage && (
            <div className="mb-4 p-3 bg-green-500/10 border border-green-500/20 rounded-lg">
              <p className="text-green-500 text-sm">{successMessage}</p>
            </div>
          )}

          {/* Current members */}
          <div className="mb-6">
            <h3 className="text-md font-medium text-telegram-text mb-3">
              {t('currentMembers', language)} ({members.length})
            </h3>
            
            {isLoading ? (
              <div className="text-center py-4">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-telegram-primary mx-auto"></div>
                <p className="text-telegram-text-secondary text-sm mt-2">
                  {t('loadingMembers', language)}
                </p>
              </div>
            ) : members.length === 0 ? (
              <div className="text-center py-4">
                <Users className="w-8 h-8 text-telegram-secondary mx-auto mb-2" />
                <p className="text-telegram-text-secondary text-sm">
                  {t('noMembers', language)}
                </p>
              </div>
            ) : (
              <div className="space-y-2">
                {members.map((member) => (
                  <div
                    key={member.id}
                    className="flex items-center justify-between p-3 bg-telegram-bg rounded-lg border border-telegram-border"
                  >
                    <div className="flex items-center">
                      <div className="w-8 h-8 rounded-full bg-gradient-to-br from-telegram-primary to-blue-600 flex items-center justify-center text-white font-medium mr-3 shadow-sm">
                        {member.first_name?.charAt(0).toUpperCase() || 'U'}
                      </div>
                      <div>
                        <p className="text-sm font-medium text-telegram-text">
                          {member.first_name} {member.last_name || ''}
                        </p>
                        <p className="text-xs text-telegram-text-secondary">
                          @{member.username} • {getMemberRole(member)}
                        </p>
                      </div>
                    </div>
                    
                    {chat.type === 'group' && (
                      <button
                        onClick={() => handleRemoveMember(member.id)}
                        disabled={isRemovingMember}
                        className="p-2 text-red-500 hover:text-red-600 hover:bg-red-500/10 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                        title={t('removeMember', language)}
                      >
                        {isRemovingMember ? (
                          <div className="w-4 h-4 border-2 border-red-500 border-t-transparent rounded-full animate-spin"></div>
                        ) : (
                          <UserMinus className="w-4 h-4" />
                        )}
                      </button>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Adding members (only for groups) */}
          {chat.type === 'group' && (
            <div>
              <h3 className="text-md font-medium text-telegram-text mb-3">
                {t('addMembers', language)}
              </h3>
              
              {availableUsers.length === 0 ? (
                <div className="text-center py-4">
                  <UserPlus className="w-8 h-8 text-telegram-secondary mx-auto mb-2" />
                  <p className="text-telegram-text-secondary text-sm">
                    {t('noAvailableUsers', language)}
                  </p>
                </div>
              ) : (
                <div className="space-y-2">
                  {availableUsers.map((user) => (
                    <div
                      key={user.id}
                      className="flex items-center justify-between p-3 bg-telegram-bg rounded-lg border border-telegram-border"
                    >
                      <div className="flex items-center">
                        <div className="w-8 h-8 rounded-full bg-gradient-to-br from-telegram-primary to-blue-600 flex items-center justify-center text-white font-medium mr-3 shadow-sm">
                          {user.first_name?.charAt(0).toUpperCase() || 'U'}
                        </div>
                        <div>
                          <p className="text-sm font-medium text-telegram-text">
                            {user.first_name} {user.last_name || ''}
                          </p>
                          <p className="text-xs text-telegram-text-secondary">
                            @{user.username}
                          </p>
                        </div>
                      </div>
                      
                      <button
                        onClick={() => handleAddMember(user.id)}
                        disabled={isAddingMember}
                        className="p-2 text-green-500 hover:text-green-600 hover:bg-green-500/10 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                        title={t('addMember', language)}
                      >
                        {isAddingMember ? (
                          <div className="w-4 h-4 border-2 border-green-500 border-t-transparent rounded-full animate-spin"></div>
                        ) : (
                          <UserPlus className="w-4 h-4" />
                        )}
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ChatMembersModal;
