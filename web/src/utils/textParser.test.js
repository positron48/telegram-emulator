import { parseTelegramText } from './textParser.jsx';

// Тестовые данные из реального сообщения
const testMessage = `📊 *Статистика и отчеты*

\`/stats [период]\` - Общая статистика
Показывает доходы и расходы за период

*Варианты периода:*
• /stats - Текущий месяц
• \`/stats 2023-12\` - Конкретный месяц (YYYY-MM)
• \`/stats week\` - Текущая неделя

\`/top_categories [период] [лимит]\` - Топ категорий
Показывает категории с наибольшими расходами

*Примеры:*
• /top_categories - Топ-5 за текущий месяц
• \`/top_categories 2023-12\` - Топ-5 за декабрь 2023
• \`/top_categories week 10\` - Топ-10 за неделю

\`/recent [лимит]\` - Последние транзакции
Показывает последние транзакции

*Примеры:*
• /recent - Последние 10 транзакций
• \`/recent 20\` - Последние 20 транзакций

\`/export [период] [лимит]\` - Экспорт данных
Экспортирует транзакции в CSV формат

*Примеры:*
• /export - Экспорт за текущий месяц
• \`/export 2023-12\` - Экспорт за декабрь 2023
• \`/export week 100\` - Экспорт 100 транзакций за неделю`;

// Функция для извлечения текста из React элементов
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

// Функция для проверки наличия элементов
const hasElement = (elements, tagName, className = null) => {
  if (!elements) return false;
  if (Array.isArray(elements)) {
    return elements.some(el => hasElement(el, tagName, className));
  }
  if (elements.type === tagName) {
    if (className) {
      return elements.props && elements.props.className && elements.props.className.includes(className);
    }
    return true;
  }
  if (elements.props && elements.props.children) {
    return hasElement(elements.props.children, tagName, className);
  }
  return false;
};

console.log("🧪 Запуск полного теста парсера...\n");

// Тест 1: Проверка заголовка
console.log("📝 Тест 1: Заголовок");
const headerTest = parseTelegramText("📊 *Статистика и отчеты*");
const headerText = extractTextFromReactElements(headerTest);
const hasBold = hasElement(headerTest, 'strong');
console.log(`  ${hasBold ? '✅' : '❌'} Заголовок жирный: ${hasBold ? 'PASS' : 'FAIL'}`);
console.log(`  Текст: "${headerText}"`);

// Тест 2: Проверка команд в коде
console.log("\n📝 Тест 2: Команды в коде");
const codeTest = parseTelegramText("`/stats [период]` - Общая статистика");
const hasCode = hasElement(codeTest, 'code');
const codeText = extractTextFromReactElements(codeTest);
console.log(`  ${hasCode ? '✅' : '❌'} Код обработан: ${hasCode ? 'PASS' : 'FAIL'}`);
console.log(`  Текст: "${codeText}"`);

// Тест 3: Проверка жирного текста
console.log("\n📝 Тест 3: Жирный текст");
const boldTest = parseTelegramText("*Варианты периода:*");
const hasBold2 = hasElement(boldTest, 'strong');
const boldText = extractTextFromReactElements(boldTest);
console.log(`  ${hasBold2 ? '✅' : '❌'} Жирный текст: ${hasBold2 ? 'PASS' : 'FAIL'}`);
console.log(`  Текст: "${boldText}"`);

// Тест 4: Проверка команд в списке
console.log("\n📝 Тест 4: Команды в списке");
const commandTest = parseTelegramText("• /stats - Текущий месяц");
const commandText = extractTextFromReactElements(commandTest);
console.log(`  Текст: "${commandText}"`);

// Тест 5: Проверка команд в коде в списке
console.log("\n📝 Тест 5: Команды в коде в списке");
const codeCommandTest = parseTelegramText("• `/stats 2023-12` - Конкретный месяц (YYYY-MM)");
const hasCode2 = hasElement(codeCommandTest, 'code');
const codeCommandText = extractTextFromReactElements(codeCommandTest);
console.log(`  ${hasCode2 ? '✅' : '❌'} Код в списке: ${hasCode2 ? 'PASS' : 'FAIL'}`);
console.log(`  Текст: "${codeCommandText}"`);

// Тест 6: Проверка команд с подчеркиванием
console.log("\n📝 Тест 6: Команды с подчеркиванием");
const underscoreTest = parseTelegramText("• /top_categories - Топ-5 за текущий месяц");
const underscoreText = extractTextFromReactElements(underscoreTest);
console.log(`  Текст: "${underscoreText}"`);

// Тест 7: Проверка команд с параметрами в коде
console.log("\n📝 Тест 7: Команды с параметрами в коде");
const paramTest = parseTelegramText("• `/top_categories 2023-12` - Топ-5 за декабрь 2023");
const hasCode3 = hasElement(paramTest, 'code');
const paramText = extractTextFromReactElements(paramTest);
console.log(`  ${hasCode3 ? '✅' : '❌'} Код с параметрами: ${hasCode3 ? 'PASS' : 'FAIL'}`);
console.log(`  Текст: "${paramText}"`);

// Тест 8: Проверка команд с пробелами в коде
console.log("\n📝 Тест 8: Команды с пробелами в коде");
const spaceTest = parseTelegramText("• `/top categories week 10` - Топ-10 за неделю");
const hasCode4 = hasElement(spaceTest, 'code');
const spaceText = extractTextFromReactElements(spaceTest);
console.log(`  ${hasCode4 ? '✅' : '❌'} Код с пробелами: ${hasCode4 ? 'PASS' : 'FAIL'}`);
console.log(`  Текст: "${spaceText}"`);

// Тест 9: Проверка курсива
console.log("\n📝 Тест 9: Курсив");
const italicTest = parseTelegramText("_курсивный текст_");
const hasItalic = hasElement(italicTest, 'em');
const italicText = extractTextFromReactElements(italicTest);
console.log(`  ${hasItalic ? '✅' : '❌'} Курсив: ${hasItalic ? 'PASS' : 'FAIL'}`);
console.log(`  Текст: "${italicText}"`);

// Тест 10: Проверка зачеркнутого текста
console.log("\n📝 Тест 10: Зачеркнутый текст");
const strikeTest = parseTelegramText("~зачеркнутый текст~");
const hasStrike = hasElement(strikeTest, 'span');
const strikeText = extractTextFromReactElements(strikeTest);
console.log(`  ${hasStrike ? '✅' : '❌'} Зачеркнутый: ${hasStrike ? 'PASS' : 'FAIL'}`);
console.log(`  Текст: "${strikeText}"`);

// Тест 11: Проверка ссылок
console.log("\n📝 Тест 11: Ссылки");
const linkTest = parseTelegramText("[Google](https://google.com)");
const hasLink = hasElement(linkTest, 'a');
const linkText = extractTextFromReactElements(linkTest);
console.log(`  ${hasLink ? '✅' : '❌'} Ссылка: ${hasLink ? 'PASS' : 'FAIL'}`);
console.log(`  Текст: "${linkText}"`);

console.log("\n✨ Полный тест завершен!");

// Экспортируем для использования в браузере
if (typeof window !== 'undefined') {
  window.runParserTests = () => {
    console.log("🧪 Запуск тестов в браузере...");
    // Здесь можно добавить визуальные тесты
  };
}
