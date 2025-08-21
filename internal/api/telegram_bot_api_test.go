package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"telegram-emulator/internal/emulator"
	"telegram-emulator/internal/models"
	"telegram-emulator/internal/repository"

	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockBotManager представляет мок для BotManager
type MockBotManager struct {
	mock.Mock
}

func (m *MockBotManager) GetAllBots() ([]models.Bot, error) {
	args := m.Called()
	return args.Get(0).([]models.Bot), args.Error(1)
}

func (m *MockBotManager) GetBotUpdates(botID string, offset, limit int) ([]models.Update, error) {
	args := m.Called(botID, offset, limit)
	return args.Get(0).([]models.Update), args.Error(1)
}

func (m *MockBotManager) AddUpdate(botID string, update *models.Update) error {
	args := m.Called(botID, update)
	return args.Error(0)
}

func (m *MockBotManager) UpdateBot(bot *models.Bot) error {
	args := m.Called(bot)
	return args.Error(0)
}

func (m *MockBotManager) GetLogger() *zap.Logger {
	return nil
}

// MockUserManager представляет мок для UserManager
type MockUserManager struct {
	mock.Mock
}

func (m *MockUserManager) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

// MockChatManager представляет мок для ChatManager
type MockChatManager struct {
	mock.Mock
}

func (m *MockChatManager) GetAllChats() ([]models.Chat, error) {
	args := m.Called()
	return args.Get(0).([]models.Chat), args.Error(1)
}

// MockMessageManager представляет мок для MessageManager
type MockMessageManager struct {
	mock.Mock
}

func (m *MockMessageManager) SendMessage(chatID, fromUserID, text, messageType string, replyMarkup interface{}) (*models.Message, error) {
	args := m.Called(chatID, fromUserID, text, messageType, replyMarkup)
	return args.Get(0).(*models.Message), args.Error(1)
}

func TestTelegramBotAPI_SendMessage(t *testing.T) {
	// Настраиваем Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Создаем моки
	mockBotManager := &MockBotManager{}
	mockUserManager := &MockUserManager{}
	mockChatManager := &MockChatManager{}
	mockMessageManager := &MockMessageManager{}

	// Создаем API
	api := NewTelegramBotAPI(mockBotManager, mockUserManager, mockChatManager, mockMessageManager)

	// Создаем тестовый бот
	testBot := models.Bot{
		ID:       "bot-1",
		Name:     "Test Bot",
		Username: "test_bot",
		Token:    "test-token",
	}

	// Создаем тестового пользователя-бота
	testBotUser := &models.User{
		ID:        "user-1",
		FirstName: "Test",
		LastName:  "Bot",
		Username:  "test_bot",
		IsBot:     true,
	}

	// Создаем тестовое сообщение
	testMessage := &models.Message{
		ID:        "msg-1",
		ChatID:    "chat-1",
		FromID:    "user-1",
		From:      *testBotUser,
		Text:      "/start",
		Type:      "text",
		Status:    "sent",
		Timestamp: time.Now(),
	}

	// Настраиваем ожидания моков
	mockBotManager.On("GetAllBots").Return([]models.Bot{testBot}, nil)
	mockUserManager.On("GetUserByUsername", "test_bot").Return(testBotUser, nil)
	mockChatManager.On("GetAllChats").Return([]models.Chat{}, nil)
	mockMessageManager.On("SendMessage", "chat-1", "user-1", "/start", "text", nil).Return(testMessage, nil)

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"chat_id": "chat-1",
		"text":    "/start",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Создаем HTTP запрос
	req, _ := http.NewRequest("POST", "/bot:test-token/sendMessage", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Создаем HTTP recorder
	w := httptest.NewRecorder()

	// Создаем Gin router
	router := gin.New()
	api.SetupTelegramBotRoutes(router)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["ok"].(bool))
	assert.NotNil(t, response["result"])

	// Проверяем, что все моки были вызваны
	mockBotManager.AssertExpectations(t)
	mockUserManager.AssertExpectations(t)
	mockChatManager.AssertExpectations(t)
	mockMessageManager.AssertExpectations(t)
}

func TestTelegramBotAPI_SendMessageWithReplyMarkup(t *testing.T) {
	// Настраиваем Gin в тестовом режиме
	gin.SetMode(gin.TestMode)

	// Создаем моки
	mockBotManager := &MockBotManager{}
	mockUserManager := &MockUserManager{}
	mockChatManager := &MockChatManager{}
	mockMessageManager := &MockMessageManager{}

	// Создаем API
	api := NewTelegramBotAPI(mockBotManager, mockUserManager, mockChatManager, mockMessageManager)

	// Создаем тестовый бот
	testBot := models.Bot{
		ID:       "bot-1",
		Name:     "Test Bot",
		Username: "test_bot",
		Token:    "test-token",
	}

	// Создаем тестового пользователя-бота
	testBotUser := &models.User{
		ID:        "user-1",
		FirstName: "Test",
		LastName:  "Bot",
		Username:  "test_bot",
		IsBot:     true,
	}

	// Создаем тестовое сообщение
	testMessage := &models.Message{
		ID:        "msg-1",
		ChatID:    "chat-1",
		FromID:    "user-1",
		From:      *testBotUser,
		Text:      "Выберите опцию:",
		Type:      "text",
		Status:    "sent",
		Timestamp: time.Now(),
	}

	// Создаем reply markup
	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "Опция 1", "callback_data": "option1"},
				{"text": "Опция 2", "callback_data": "option2"},
			},
		},
	}

	// Настраиваем ожидания моков
	mockBotManager.On("GetAllBots").Return([]models.Bot{testBot}, nil)
	mockUserManager.On("GetUserByUsername", "test_bot").Return(testBotUser, nil)
	mockChatManager.On("GetAllChats").Return([]models.Chat{}, nil)
	mockMessageManager.On("SendMessage", "chat-1", "user-1", "Выберите опцию:", "text", replyMarkup).Return(testMessage, nil)

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"chat_id":     "chat-1",
		"text":        "Выберите опцию:",
		"reply_markup": replyMarkup,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Создаем HTTP запрос
	req, _ := http.NewRequest("POST", "/bot:test-token/sendMessage", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Создаем HTTP recorder
	w := httptest.NewRecorder()

	// Создаем Gin router
	router := gin.New()
	api.SetupTelegramBotRoutes(router)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["ok"].(bool))
	assert.NotNil(t, response["result"])

	// Проверяем, что все моки были вызваны
	mockBotManager.AssertExpectations(t)
	mockUserManager.AssertExpectations(t)
	mockChatManager.AssertExpectations(t)
	mockMessageManager.AssertExpectations(t)
}

func TestTelegramBotAPI_ValidateReplyMarkup(t *testing.T) {
	api := &TelegramBotAPI{}

	// Тест валидной inline клавиатуры
	validInlineKeyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "Кнопка 1", "callback_data": "btn1"},
				{"text": "Кнопка 2", "callback_data": "btn2"},
			},
		},
	}
	err := api.validateReplyMarkup(validInlineKeyboard)
	assert.NoError(t, err)

	// Тест невалидной клавиатуры (без callback_data)
	invalidInlineKeyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "Кнопка без callback_data"},
			},
		},
	}
	err = api.validateReplyMarkup(invalidInlineKeyboard)
	assert.Error(t, err)

	// Тест валидной обычной клавиатуры
	validKeyboard := map[string]interface{}{
		"keyboard": [][]map[string]interface{}{
			{
				{"text": "Кнопка 1"},
				{"text": "Кнопка 2"},
			},
		},
	}
	err = api.validateReplyMarkup(validKeyboard)
	assert.NoError(t, err)

	// Тест невалидной клавиатуры (пустой текст)
	invalidKeyboard := map[string]interface{}{
		"keyboard": [][]map[string]interface{}{
			{
				{"text": ""},
			},
		},
	}
	err = api.validateReplyMarkup(invalidKeyboard)
	assert.Error(t, err)
}
