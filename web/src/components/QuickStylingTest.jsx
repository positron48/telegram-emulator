import React, { useState } from 'react';
import { parseTelegramText } from '../utils/textParser.jsx';

const QuickStylingTest = ({ onClose }) => {
  const [testText, setTestText] = useState('Привет! Это *жирный текст* и _курсив_. Попробуйте команду /start');

  const quickExamples = [
    "Привет! Это *жирный текст*",
    "Это _курсивный текст_",
    "Это ~зачеркнутый текст~",
    "Это `код` в сообщении",
    "Это [ссылка на Google](https://google.com)",
    "Команды: /start /help /stop",
    "*Жирный* с _курсивом_ и `кодом`"
  ];

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
              🎨 Тест стилизации текста
            </h2>
            <button
              onClick={onClose}
              className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
            >
              ✕
            </button>
          </div>

          <div className="space-y-4">
            {/* Быстрые примеры */}
            <div>
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-3">
                Быстрые примеры:
              </h3>
              <div className="flex flex-wrap gap-2">
                {quickExamples.map((example, index) => (
                  <button
                    key={index}
                    onClick={() => setTestText(example)}
                    className="px-3 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 rounded-full text-sm hover:bg-blue-200 dark:hover:bg-blue-800 transition-colors"
                  >
                    {example.substring(0, 20)}...
                  </button>
                ))}
              </div>
            </div>

            {/* Поле ввода */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Введите текст для тестирования:
              </label>
              <textarea
                value={testText}
                onChange={(e) => setTestText(e.target.value)}
                className="w-full h-24 p-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Попробуйте: *жирный*, _курсив_, ~зачеркнутый~, `код`, [ссылка](https://example.com)"
              />
            </div>

            {/* Предварительный просмотр */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Результат:
              </label>
              <div className="p-4 border border-gray-300 dark:border-gray-600 rounded-lg bg-gray-50 dark:bg-gray-700 min-h-24">
                <div className="whitespace-pre-wrap break-words text-gray-900 dark:text-white">
                  {parseTelegramText(testText)}
                </div>
              </div>
            </div>

            {/* Справка */}
            <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
              <h4 className="text-sm font-medium text-blue-900 dark:text-blue-100 mb-2">
                💡 Доступные форматы:
              </h4>
                             <div className="grid grid-cols-2 gap-2 text-xs text-blue-800 dark:text-blue-200">
                 <div>• <code className="bg-blue-100 dark:bg-blue-800 text-blue-900 dark:text-blue-100 px-1 rounded">*текст*</code> - жирный</div>
                 <div>• <code className="bg-blue-100 dark:bg-blue-800 text-blue-900 dark:text-blue-100 px-1 rounded">_текст_</code> - курсив</div>
                 <div>• <code className="bg-blue-100 dark:bg-blue-800 text-blue-900 dark:text-blue-100 px-1 rounded">~текст~</code> - зачеркнутый</div>
                 <div>• <code className="bg-blue-100 dark:bg-blue-800 text-blue-900 dark:text-blue-100 px-1 rounded">`код`</code> - моноширинный</div>
                 <div>• <code className="bg-blue-100 dark:bg-blue-800 text-blue-900 dark:text-blue-100 px-1 rounded">[текст](ссылка)</code> - ссылка</div>
                 <div>• <code className="bg-blue-100 dark:bg-blue-800 text-blue-900 dark:text-blue-100 px-1 rounded">/команда</code> - команда</div>
               </div>
            </div>
          </div>

          <div className="mt-6 flex justify-end">
            <button
              onClick={onClose}
              className="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors"
            >
              Закрыть
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default QuickStylingTest;
