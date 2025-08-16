# Исправление бага: Боты не появляются в списке ботов

## 🐛 Проблема

При создании пользователя с типом "бот" через интерфейс, пользователь создавался только в таблице `users` с флагом `is_bot: true`, но не создавалась запись в таблице `bots`. Поэтому бот не появлялся в списке ботов в панели управления.

## 🔧 Решение

### 1. Обновление UserManager

Добавлена зависимость от `BotRepository` в `UserManager`:

```go
type UserManager struct {
    userRepo *repository.UserRepository
    botRepo  *repository.BotRepository  // Добавлено
    logger   *zap.Logger
}
```

### 2. Обновление конструктора UserManager

```go
func NewUserManager(userRepo *repository.UserRepository, botRepo *repository.BotRepository) *UserManager {
    return &UserManager{
        userRepo: userRepo,
        botRepo:  botRepo,  // Добавлено
        logger:   logger.GetLogger(),
    }
}
```

### 3. Обновление метода CreateUser

При создании пользователя-бота теперь также создается запись в таблице `bots`:

```go
// Если пользователь является ботом, создаем запись в таблице ботов
if user.IsBot {
    bot := &models.Bot{
        ID:         user.ID,
        Name:       user.GetFullName(),
        Username:   user.Username,
        Token:      "", // Токен можно будет установить позже
        WebhookURL: "",
        IsActive:   true,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    if err := m.botRepo.Create(bot); err != nil {
        m.logger.Error("Ошибка создания записи бота", zap.Error(err))
        // Не удаляем пользователя, просто логируем ошибку
    } else {
        m.logger.Info("Создана запись бота", 
            zap.String("bot_id", bot.ID),
            zap.String("username", bot.Username))
    }
}
```

### 4. Обновление метода DeleteUser

При удалении пользователя-бота теперь также удаляется запись из таблицы `bots`:

```go
// Если пользователь является ботом, удаляем запись бота
if user.IsBot {
    if err := m.botRepo.Delete(id); err != nil {
        m.logger.Error("Ошибка удаления записи бота", zap.String("id", id), zap.Error(err))
        // Не прерываем удаление пользователя, просто логируем ошибку
    } else {
        m.logger.Info("Запись бота удалена", zap.String("id", id))
    }
}
```

### 5. Обновление главного файла приложения

Обновлен вызов конструктора `UserManager` в `cmd/emulator/main.go`:

```go
userManager := emulator.NewUserManager(userRepo, botRepo)  // Добавлен botRepo
```

### 6. Добавление компонента CreateBotModal

Создан новый компонент `CreateBotModal.jsx` для создания ботов через интерфейс с полными настройками (токен, webhook URL).

### 7. Обновление BotManager

Добавлена кнопка создания бота и интеграция с `CreateBotModal`.

## ✅ Результат

Теперь при создании пользователя с типом "бот":

1. **Создается пользователь** в таблице `users` с `is_bot: true`
2. **Создается запись бота** в таблице `bots` с пустым токеном и webhook URL
3. **Бот появляется** в списке ботов в панели управления
4. **Можно обновить** токен и webhook URL через API или интерфейс

## 🧪 Тестирование

### Создание бота через API
```bash
curl -X POST http://localhost:3001/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_bot",
    "first_name": "Test",
    "last_name": "Bot",
    "is_bot": true
  }'
```

### Проверка списка ботов
```bash
curl http://localhost:3001/api/bots
```

### Обновление токена бота
```bash
curl -X PUT http://localhost:3001/api/bots/{bot_id} \
  -H "Content-Type: application/json" \
  -d '{
    "token": "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"
  }'
```

## 🔄 Синхронизация данных

Теперь данные между таблицами `users` и `bots` синхронизированы:

- При создании пользователя-бота → создается запись в `bots`
- При удалении пользователя-бота → удаляется запись из `bots`
- При обновлении пользователя-бота → можно обновить данные в `bots`

## 📝 Примечания

- Токен и webhook URL изначально пустые и могут быть установлены позже
- При ошибке создания записи бота пользователь все равно создается
- Логирование помогает отслеживать процесс создания/удаления ботов
