import React, { useState } from 'react';
import { parseTelegramText } from '../utils/textParser.jsx';

const StylingDemoPage = () => {
  const [inputText, setInputText] = useState('');
  const [previewText, setPreviewText] = useState('');

  const examples = [
    {
      title: "Жирный текст",
      text: "Это *жирный текст* в стиле Telegram",
      description: "Используйте *текст* для жирного начертания"
    },
    {
      title: "Курсив",
      text: "Это _курсивный текст_ в стиле Telegram",
      description: "Используйте _текст_ для курсива"
    },
    {
      title: "Зачеркнутый текст",
      text: "Это ~зачеркнутый текст~ в стиле Telegram",
      description: "Используйте ~текст~ для зачеркивания"
    },
    {
      title: "Код",
      text: "Это `код` в стиле Telegram",
      description: "Используйте `код` для моноширинного шрифта"
    },
    {
      title: "Ссылки",
      text: "Это [ссылка на Google](https://google.com) в стиле Telegram",
      description: "Используйте [текст](ссылка) для создания ссылок"
    },
    {
      title: "Комбинированное форматирование",
      text: "Это *жирный текст* с _курсивом_ и `кодом`",
      description: "Можно комбинировать разные стили"
    },
    {
      title: "Команды",
      text: "Это команда /start в стиле Telegram",
      description: "Команды автоматически выделяются синим цветом"
    }
  ];

  const handleInputChange = (e) => {
    const text = e.target.value;
    setInputText(text);
    setPreviewText(text);
  };

  const handleExampleClick = (exampleText) => {
    setInputText(exampleText);
    setPreviewText(exampleText);
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 p-6">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 mb-6">
          <h1 className="text-3xl font-bold mb-6 text-gray-900 dark:text-white">
            Демонстрация стилизации текста в стиле Telegram
          </h1>
          
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Интерактивная область */}
            <div className="space-y-4">
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
                Попробуйте сами
              </h2>
              
              <div>
                <label htmlFor="demo-textarea" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Введите текст с форматированием:
                </label>
                <textarea
                  id="demo-textarea"
                  value={inputText}
                  onChange={handleInputChange}
                  placeholder="Попробуйте: *жирный*, _курсив_, ~зачеркнутый~, `код`, [ссылка](https://example.com)"
                  className="w-full h-32 p-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              
              <div>
                <label htmlFor="demo-preview" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Предварительный просмотр:
                </label>
                <div id="demo-preview" className="w-full min-h-32 p-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white">
                  <div className="whitespace-pre-wrap break-words">
                    {parseTelegramText(previewText)}
                  </div>
                </div>
              </div>
            </div>

            {/* Примеры */}
            <div>
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
                Примеры форматирования
              </h2>
              <div className="space-y-3">
                {examples.map((example, index) => (
                  <div 
                    key={index} 
                    className="p-3 border border-gray-200 dark:border-gray-600 rounded-lg cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                    onClick={() => handleExampleClick(example.text)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter' || e.key === ' ') {
                        e.preventDefault();
                        handleExampleClick(example.text);
                      }
                    }}
                    tabIndex={0}
                    role="button"
                    aria-label={`Use example: ${example.title}`}
                  >
                    <h3 className="text-sm font-medium text-gray-900 dark:text-white mb-1">
                      {example.title}
                    </h3>
                    <p className="text-xs text-gray-600 dark:text-gray-400 mb-2">
                      {example.description}
                    </p>
                    <div className="text-sm text-gray-700 dark:text-gray-300">
                      {parseTelegramText(example.text)}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Справка по форматам */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
          <h2 className="text-xl font-semibold mb-4 text-gray-900 dark:text-white">
            Поддерживаемые форматы
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <h3 className="text-lg font-medium mb-3 text-gray-900 dark:text-white">
                Основное форматирование
              </h3>
                             <ul className="space-y-2 text-sm text-gray-700 dark:text-gray-300">
                 <li>• <code className="bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-2 py-1 rounded">*текст*</code> - жирный текст</li>
                 <li>• <code className="bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-2 py-1 rounded">_текст_</code> - курсив</li>
                 <li>• <code className="bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-2 py-1 rounded">~текст~</code> - зачеркнутый текст</li>
                 <li>• <code className="bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-2 py-1 rounded">`код`</code> - моноширинный шрифт</li>
               </ul>
            </div>
            <div>
              <h3 className="text-lg font-medium mb-3 text-gray-900 dark:text-white">
                Ссылки и команды
              </h3>
                             <ul className="space-y-2 text-sm text-gray-700 dark:text-gray-300">
                 <li>• <code className="bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-2 py-1 rounded">[текст](ссылка)</code> - ссылки</li>
                 <li>• <code className="bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-2 py-1 rounded">/команда</code> - команды (автоматически выделяются)</li>
                 <li>• <code className="bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-2 py-1 rounded">https://example.com</code> - автоматические ссылки</li>
               </ul>
            </div>
          </div>
          
          <div className="mt-6 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
            <h3 className="text-lg font-medium mb-2 text-blue-900 dark:text-blue-100">
              💡 Подсказка
            </h3>
                         <p className="text-sm text-blue-800 dark:text-blue-200">
               Форматы можно комбинировать! Например: <code className="bg-blue-100 dark:bg-blue-800 text-blue-900 dark:text-blue-100 px-1 rounded">*жирный _с курсивом_*</code>
             </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default StylingDemoPage;
