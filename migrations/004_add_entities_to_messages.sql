-- Добавление поля entities в таблицу messages
ALTER TABLE messages ADD COLUMN entities TEXT;

-- Создание индекса для быстрого поиска по entities
CREATE INDEX idx_messages_entities ON messages(entities);
