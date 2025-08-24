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

  // Сначала обрабатываем ссылки в формате [текст](ссылка)
  const linkRegex = /\[([^\]]+)\]\(([^)]+)\)/g;
  let processedText = text.replace(linkRegex, (match, linkText, url) => {
    return `__LINK_START__${linkText}__LINK_URL__${url}__LINK_END__`;
  });

  // Обрабатываем жирный текст: *текст*
  const boldRegex = /\*([^*]+)\*/g;
  processedText = processedText.replace(boldRegex, (match, content) => {
    return `__BOLD_START__${content}__BOLD_END__`;
  });

  // Обрабатываем курсив: _текст_
  const italicRegex = /_([^_]+)_/g;
  processedText = processedText.replace(italicRegex, (match, content) => {
    return `__ITALIC_START__${content}__ITALIC_END__`;
  });

  // Обрабатываем зачеркнутый текст: ~текст~
  const strikethroughRegex = /~([^~]+)~/g;
  processedText = processedText.replace(strikethroughRegex, (match, content) => {
    return `__STRIKE_START__${content}__STRIKE_END__`;
  });

  // Обрабатываем код: `код`
  const codeRegex = /`([^`]+)`/g;
  processedText = processedText.replace(codeRegex, (match, content) => {
    return `__CODE_START__${content}__CODE_END__`;
  });

  // Разбиваем текст на части по маркерам
  const parts = processedText.split(/(__[A-Z_]+__)/g);
  
  // Исправлено: используем цикл for вместо map для правильной обработки индексов
  // Это предотвращает дублирование текста в элементах <code>, <strong>, etc.
  const result = [];
  
  for (let i = 0; i < parts.length; i++) {
    const part = parts[i];
    
    if (part === '__BOLD_START__') {
      result.push(<strong key={i} className="font-bold">{parts[i + 1]}</strong>);
      i++; // Пропускаем следующий элемент (содержимое)
    } else if (part === '__ITALIC_START__') {
      result.push(<em key={i} className="italic">{parts[i + 1]}</em>);
      i++; // Пропускаем следующий элемент (содержимое)
    } else if (part === '__STRIKE_START__') {
      result.push(<span key={i} className="line-through">{parts[i + 1]}</span>);
      i++; // Пропускаем следующий элемент (содержимое)
    } else if (part === '__CODE_START__') {
      // Исправлено: добавлены цвета текста для темной темы
      // Светлая тема: серый фон, темный текст
      // Темная тема: темно-серый фон, светлый текст
      result.push(<code key={i} className="bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100 px-1 py-0.5 rounded text-sm font-mono">{parts[i + 1]}</code>);
      i++; // Пропускаем следующий элемент (содержимое)
    } else if (part === '__LINK_START__') {
      const linkText = parts[i + 1];
      const url = parts[i + 3]; // __LINK_URL__ находится на index + 2
      result.push(
        <a
          key={i}
          href={url}
          target="_blank"
          rel="noopener noreferrer"
          className="text-blue-500 hover:text-blue-600 dark:text-blue-400 dark:hover:text-blue-300 underline"
        >
          {linkText}
        </a>
      );
      i += 3; // Пропускаем текст ссылки, URL и маркер окончания
    } else if (part.match(/__[A-Z_]+_END__/)) {
      // Пропускаем маркеры окончания
      continue;
    } else if (part.match(/__LINK_URL__/)) {
      // Пропускаем URL части ссылок
      continue;
    } else {
      // Обычный текст
      result.push(part);
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
