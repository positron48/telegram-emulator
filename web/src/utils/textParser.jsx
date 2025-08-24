import React from 'react';

/**
 * Парсит текст в стиле Telegram и возвращает React элементы с форматированием
 * Поддерживаемые форматы:
 * - *текст* - жирный
 * - _текст_ - курсив
 * - ~текст~ - зачеркнутый
 * - `код` - моноширинный шрифт
 * - [текст](ссылка) - ссылки
 * 
 * @param {string} text - Исходный текст
 * @returns {Array} Массив React элементов
 */
export const parseTelegramText = (text) => {
  if (!text) return null;

  const parts = [];
  let i = 0;
  let currentText = '';
  
  while (i < text.length) {
    // Кодблоки `code`
    if (text[i] === '`') {
      // Добавляем накопленный текст
      if (currentText) {
        parts.push(currentText);
        currentText = '';
      }
      
      const start = i;
      i++;
      while (i < text.length && text[i] !== '`') {
        i++;
      }
      if (i < text.length) {
        const codeContent = text.substring(start + 1, i);
        parts.push(
          <code
            key={`code-${parts.length}`}
            className="bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100 px-1 py-0.5 rounded text-sm font-mono"
          >
            {codeContent}
          </code>
        );
        i++;
        continue;
      } else {
        // Если не нашли закрывающий символ, добавляем ` как обычный текст
        currentText += text[start];
        i = start + 1;
        continue;
      }
    }
    
    // Жирный текст *bold* - только в начале/конце строк или после пробелов
    if (text[i] === '*') {
      // Проверяем, что перед * есть пробел, начало строки или это первый символ
      const isStartOfLine = i === 0 || text[i - 1] === '\n' || text[i - 1] === '\r';
      const isAfterSpace = i > 0 && /\s/.test(text[i - 1]);
      
      if (isStartOfLine || isAfterSpace) {
        // Добавляем накопленный текст
        if (currentText) {
          parts.push(currentText);
          currentText = '';
        }
        
        const start = i;
        i++;
        while (i < text.length && text[i] !== '*') {
          i++;
        }
        if (i < text.length) {
          const boldContent = text.substring(start + 1, i);
          parts.push(
            <strong key={`bold-${parts.length}`} className="font-bold">
              {boldContent}
            </strong>
          );
          i++;
          continue;
        } else {
          // Если не нашли закрывающий символ, добавляем * как обычный текст
          currentText += text[start];
          i = start + 1;
          continue;
        }
      }
    }
    
    // Курсив _italic_ - только в начале/конце строк или после пробелов
    if (text[i] === '_') {
      // Проверяем, что перед _ есть пробел, начало строки или это первый символ
      const isStartOfLine = i === 0 || text[i - 1] === '\n' || text[i - 1] === '\r';
      const isAfterSpace = i > 0 && /\s/.test(text[i - 1]);
      
      if (isStartOfLine || isAfterSpace) {
        // Добавляем накопленный текст
        if (currentText) {
          parts.push(currentText);
          currentText = '';
        }
        
        const start = i;
        i++;
        while (i < text.length && text[i] !== '_') {
          i++;
        }
        if (i < text.length) {
          const italicContent = text.substring(start + 1, i);
          parts.push(
            <em key={`italic-${parts.length}`} className="italic">
              {italicContent}
            </em>
          );
          i++;
          continue;
        } else {
          // Если не нашли закрывающий символ, добавляем _ как обычный текст
          currentText += text[start];
          i = start + 1;
          continue;
        }
      }
    }
    
    // Зачеркнутый ~strikethrough~ - только в начале/конце строк или после пробелов
    if (text[i] === '~') {
      // Проверяем, что перед ~ есть пробел, начало строки или это первый символ
      const isStartOfLine = i === 0 || text[i - 1] === '\n' || text[i - 1] === '\r';
      const isAfterSpace = i > 0 && /\s/.test(text[i - 1]);
      
      if (isStartOfLine || isAfterSpace) {
        // Добавляем накопленный текст
        if (currentText) {
          parts.push(currentText);
          currentText = '';
        }
        
        const start = i;
        i++;
        while (i < text.length && text[i] !== '~') {
          i++;
        }
        if (i < text.length) {
          const strikeContent = text.substring(start + 1, i);
          parts.push(
            <del key={`strike-${parts.length}`} className="line-through">
              {strikeContent}
            </del>
          );
          i++;
          continue;
        } else {
          // Если не нашли закрывающий символ, добавляем ~ как обычный текст
          currentText += text[start];
          i = start + 1;
          continue;
        }
      }
    }
    
    // Ссылки [text](url)
    if (text[i] === '[') {
      // Добавляем накопленный текст
      if (currentText) {
        parts.push(currentText);
        currentText = '';
      }
      
      const start = i;
      i++;
      while (i < text.length && text[i] !== ']') {
        i++;
      }
      if (i < text.length && text[i + 1] === '(') {
        const linkText = text.substring(start + 1, i);
        i += 2;
        const urlStart = i;
        while (i < text.length && text[i] !== ')') {
          i++;
        }
        if (i < text.length) {
          const url = text.substring(urlStart, i);
          parts.push(
            <a
              key={`link-${parts.length}`}
              href={url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-500 hover:text-blue-600 dark:text-blue-400 dark:hover:text-blue-300 underline"
            >
              {linkText}
            </a>
          );
          i++;
          continue;
        }
      }
    }
    
    // Обычный текст
    currentText += text[i];
    i++;
  }
  
  // Добавляем оставшийся текст
  if (currentText) {
    parts.push(currentText);
  }
  
  return parts;
};

/**
 * Обрабатывает команды в уже отформатированном тексте
 * @param {Array} formattedElements - Массив React элементов
 * @param {Function} onSendMessage - Функция для отправки команды
 * @returns {Array} Массив React элементов с обработанными командами
 */
export const processCommandsInFormattedText = (formattedElements, onSendMessage) => {
  if (!formattedElements) return formattedElements;
  
  // Если это строка, обрабатываем команды
  if (typeof formattedElements === 'string') {
    return processCommandsInText(formattedElements, 'text', onSendMessage);
  }
  
  // Если это не массив, возвращаем как есть
  if (!Array.isArray(formattedElements)) {
    return formattedElements;
  }
  
  return formattedElements.map((element, index) => {
    // Если это строка, обрабатываем команды
    if (typeof element === 'string') {
      return processCommandsInText(element, `text-${index}`, onSendMessage);
    }
    
    // Если это React элемент с кодом, НЕ обрабатываем команды внутри
    // Кодблоки должны отображаться как есть
    if (element && element.type === 'code') {
      return element;
    }
    
    // Если это объект с типом 'code', НЕ обрабатываем команды внутри
    if (element && typeof element === 'object' && element.type === 'code') {
      return element;
    }
    
    // Если это React элемент с детьми, рекурсивно обрабатываем
    if (element && element.props && element.props.children) {
      const processedChildren = processCommandsInFormattedText(element.props.children, onSendMessage);
      return React.cloneElement(element, { key: `element-${index}` }, processedChildren);
    }
    
    // Если это объект с контентом, рекурсивно обрабатываем
    if (element && typeof element === 'object' && element.content) {
      const processedContent = processCommandsInFormattedText(element.content, onSendMessage);
      return {
        ...element,
        content: processedContent
      };
    }
    
    return element;
  });
};

/**
 * Обрабатывает команды в тексте
 * @param {string} text - Текст для обработки
 * @param {string} keyPrefix - Префикс для ключей
 * @param {Function} onSendMessage - Функция для отправки команды
 * @returns {Array} Массив React элементов
 */
const processCommandsInText = (text, keyPrefix, onSendMessage) => {
  if (!text) return text;
  
  // Команды должны быть в начале строки или после пробела
  // Более точный regex для команд - исправлен для лучшего распознавания
  const commandRegex = /(^|\s)(\/[a-zA-Z0-9_]+(?:\s+[a-zA-Z0-9_\-\.]+)*)(?=\s|$|\.|,|;|!|\?)/g;
  const commandParts = text.split(commandRegex);
  
  const result = [];
  
  for (let i = 0; i < commandParts.length; i++) {
    const commandPart = commandParts[i];
    
    if (commandRegex.test(commandPart)) {
      // Это пробел + команда, разделяем
      const space = commandPart.charAt(0) === ' ' ? ' ' : '';
      const command = commandPart.substring(space.length);
      
      result.push(
        <React.Fragment key={`${keyPrefix}-command-${i}`}>
          {space}
          <button
            style={{
              color: 'rgb(59, 130, 246)',
              cursor: 'pointer',
              fontWeight: '500',
              transition: 'color 0.2s'
            }}
            className="hover:text-blue-600 dark:hover:text-[rgb(156,197,255)] inline-flex items-center"
            onClick={() => {
              if (onSendMessage) {
                onSendMessage(command);
              }
              console.log('Command clicked:', command);
            }}
            title={`Click to send: ${command}`}
          >
            <span 
              style={{
                color: 'rgb(59, 130, 246)',
                fontWeight: '500'
              }}
              className="dark:text-[rgb(121,171,252)]"
            >
              {command}
            </span>
          </button>
        </React.Fragment>
      );
    } else {
      result.push(commandPart);
    }
  }
  
  return result;
};

/**
 * Проверяет, содержит ли текст команды Telegram
 * @param {string} text - Текст для проверки
 * @returns {boolean} true если содержит команды
 */
export const hasCommands = (text) => {
  if (!text) return false;
  const commandRegex = /(^|\s)(\/[a-zA-Z0-9_]+)/g;
  return commandRegex.test(text);
};

/**
 * Извлекает команды из текста
 * @param {string} text - Текст для извлечения команд
 * @returns {Array} Массив найденных команд
 */
export const extractCommands = (text) => {
  if (!text) return [];
  const commandRegex = /(^|\s)(\/[a-zA-Z0-9_]+)/g;
  const commands = [];
  let match;
  
  while ((match = commandRegex.exec(text)) !== null) {
    commands.push(match[2]); // match[2] содержит команду без пробела
  }
  
  return commands;
};
