// Типы данных для Telegram Emulator

export const MessageStatus = {
  SENDING: 'sending',
  SENT: 'sent',
  DELIVERED: 'delivered',
  READ: 'read'
};

export const MessageType = {
  TEXT: 'text',
  FILE: 'file',
  VOICE: 'voice',
  PHOTO: 'photo'
};

export const ChatType = {
  PRIVATE: 'private',
  GROUP: 'group'
};

export const DebugEventType = {
  MESSAGE: 'message',
  API_CALL: 'api_call',
  ERROR: 'error',
  INFO: 'info'
};

/**
 * @typedef {Object} User
 * @property {string} id - Уникальный идентификатор пользователя
 * @property {string} username - Имя пользователя
 * @property {string} first_name - Имя
 * @property {string} last_name - Фамилия
 * @property {boolean} is_bot - Является ли ботом
 * @property {boolean} is_online - Онлайн статус
 * @property {string} last_seen - Время последнего появления
 * @property {string} created_at - Время создания
 * @property {string} updated_at - Время обновления
 */

/**
 * @typedef {Object} Chat
 * @property {string} id - Уникальный идентификатор чата
 * @property {ChatType} type - Тип чата
 * @property {string} title - Название чата
 * @property {string} username - Имя пользователя чата
 * @property {string} description - Описание чата
 * @property {User[]} members - Участники чата
 * @property {Message|null} last_message - Последнее сообщение
 * @property {number} unread_count - Количество непрочитанных сообщений
 * @property {string} created_at - Время создания
 * @property {string} updated_at - Время обновления
 */

/**
 * @typedef {Object} Message
 * @property {string} id - Уникальный идентификатор сообщения
 * @property {string} chat_id - ID чата
 * @property {string} from_id - ID отправителя
 * @property {User} from - Информация об отправителе
 * @property {string} text - Текст сообщения
 * @property {MessageType} type - Тип сообщения
 * @property {MessageStatus} status - Статус сообщения
 * @property {boolean} is_outgoing - Исходящее ли сообщение
 * @property {string} timestamp - Время отправки
 * @property {string} created_at - Время создания
 */

/**
 * @typedef {Object} Bot
 * @property {string} id - Уникальный идентификатор бота
 * @property {string} name - Имя бота
 * @property {string} username - Имя пользователя бота
 * @property {string} token - Токен бота
 * @property {string} webhook_url - URL webhook
 * @property {boolean} is_active - Активен ли бот
 * @property {string} created_at - Время создания
 * @property {string} updated_at - Время обновления
 */

/**
 * @typedef {Object} DebugEvent
 * @property {string} id - Уникальный идентификатор события
 * @property {string} timestamp - Время события
 * @property {DebugEventType} type - Тип события
 * @property {string} description - Описание события
 * @property {any} data - Дополнительные данные
 */


