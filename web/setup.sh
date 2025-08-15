#!/bin/bash

echo "🚀 Настройка Telegram Emulator Web Interface..."

# Проверка наличия Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js не установлен. Пожалуйста, установите Node.js 18+"
    exit 1
fi

# Проверка версии Node.js
NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 18 ]; then
    echo "❌ Требуется Node.js версии 18 или выше. Текущая версия: $(node -v)"
    exit 1
fi

echo "✅ Node.js версии $(node -v) найден"

# Установка зависимостей
echo "📦 Установка зависимостей..."
npm install

if [ $? -eq 0 ]; then
    echo "✅ Зависимости установлены успешно"
else
    echo "❌ Ошибка при установке зависимостей"
    exit 1
fi

echo ""
echo "🎉 Настройка завершена!"
echo ""
echo "Для запуска в режиме разработки выполните:"
echo "  npm run dev"
echo ""
echo "Для сборки проекта выполните:"
echo "  npm run build"
echo ""
echo "Приложение будет доступно по адресу: http://localhost:3000"
