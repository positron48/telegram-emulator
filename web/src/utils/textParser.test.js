import { parseTelegramText, hasCommands, extractCommands } from './textParser.jsx';

// Простые тесты для проверки работы парсера
// Для запуска: node textParser.test.js

const testCases = [
  {
    name: "Жирный текст",
    input: "Это *жирный текст*",
    expected: "Это жирный текст" // Упрощенная проверка
  },
  {
    name: "Курсив",
    input: "Это _курсивный текст_",
    expected: "Это курсивный текст"
  },
  {
    name: "Зачеркнутый текст",
    input: "Это ~зачеркнутый текст~",
    expected: "Это зачеркнутый текст"
  },
  {
    name: "Код",
    input: "Это `код`",
    expected: "Это код"
  },
  {
    name: "Ссылка",
    input: "Это [ссылка](https://example.com)",
    expected: "Это ссылка"
  },
  {
    name: "Комбинированное форматирование",
    input: "Это *жирный _с курсивом_* текст",
    expected: "Это жирный с курсивом текст"
  },
  {
    name: "Команда",
    input: "Используйте /start",
    expected: "Используйте /start"
  }
];

const commandTests = [
  {
    name: "Текст с командой",
    input: "Используйте /start для начала",
    hasCommands: true,
    commands: ["/start"]
  },
  {
    name: "Текст без команд",
    input: "Обычный текст без команд",
    hasCommands: false,
    commands: []
  },
  {
    name: "Множественные команды",
    input: "Команды: /start /help /stop",
    hasCommands: true,
    commands: ["/start", "/help", "/stop"]
  }
];

// Функция для извлечения текста из React элементов (упрощенная)
const extractTextFromReactElements = (elements) => {
  if (!elements) return "";
  if (typeof elements === 'string') return elements;
  if (Array.isArray(elements)) {
    return elements.map(extractTextFromReactElements).join('');
  }
  if (elements.props && elements.props.children) {
    return extractTextFromReactElements(elements.props.children);
  }
  return "";
};

console.log("🧪 Запуск тестов для textParser...\n");

// Тестирование parseTelegramText
console.log("📝 Тестирование parseTelegramText:");
testCases.forEach((testCase, index) => {
  try {
    const result = parseTelegramText(testCase.input);
    const resultText = extractTextFromReactElements(result);
    
    // Упрощенная проверка - проверяем, что результат содержит ожидаемый текст
    const passed = resultText.includes(testCase.expected.replace(/\*|_|~|`|\[.*?\]\(.*?\)/g, ''));
    
    console.log(`  ${passed ? '✅' : '❌'} ${testCase.name}: ${passed ? 'PASS' : 'FAIL'}`);
    if (!passed) {
      console.log(`    Ожидалось: ${testCase.expected}`);
      console.log(`    Получено: ${resultText}`);
    }
  } catch (error) {
    console.log(`  ❌ ${testCase.name}: ERROR - ${error.message}`);
  }
});

console.log("\n🔍 Тестирование hasCommands:");
commandTests.forEach((testCase) => {
  try {
    const result = hasCommands(testCase.input);
    const passed = result === testCase.hasCommands;
    
    console.log(`  ${passed ? '✅' : '❌'} ${testCase.name}: ${passed ? 'PASS' : 'FAIL'}`);
    if (!passed) {
      console.log(`    Ожидалось: ${testCase.hasCommands}`);
      console.log(`    Получено: ${result}`);
    }
  } catch (error) {
    console.log(`  ❌ ${testCase.name}: ERROR - ${error.message}`);
  }
});

console.log("\n📋 Тестирование extractCommands:");
commandTests.forEach((testCase) => {
  try {
    const result = extractCommands(testCase.input);
    const passed = JSON.stringify(result) === JSON.stringify(testCase.commands);
    
    console.log(`  ${passed ? '✅' : '❌'} ${testCase.name}: ${passed ? 'PASS' : 'FAIL'}`);
    if (!passed) {
      console.log(`    Ожидалось: ${JSON.stringify(testCase.commands)}`);
      console.log(`    Получено: ${JSON.stringify(result)}`);
    }
  } catch (error) {
    console.log(`  ❌ ${testCase.name}: ERROR - ${error.message}`);
  }
});

console.log("\n✨ Тестирование завершено!");
