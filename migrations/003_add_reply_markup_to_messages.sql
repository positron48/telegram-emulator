-- Миграция для добавления поля reply_markup в таблицу messages
-- Добавляем поле для хранения клавиатур в JSON формате

ALTER TABLE messages ADD COLUMN reply_markup TEXT;
