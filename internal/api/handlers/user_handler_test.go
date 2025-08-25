package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"telegram-emulator/internal/emulator"
	"telegram-emulator/internal/models"
	"telegram-emulator/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockUserManager is a mock implementation of UserManager
type MockUserManager struct {
	mock.Mock
}

func (m *MockUserManager) CreateUser(username, firstName, lastName string, isBot bool) (*models.User, error) {
	args := m.Called(username, firstName, lastName, isBot)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserManager) GetUser(id int64) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserManager) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserManager) UpdateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserManager) DeleteUser(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserManager) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Chat{}, &models.Message{}, &models.Bot{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func setupTestRouter(t *testing.T) (*gin.Engine, *MockUserManager) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	mockUserManager := &MockUserManager{}
	userHandler := NewUserHandler(mockUserManager)
	
	// Setup routes
	router.GET("/users", userHandler.GetAll)
	router.POST("/users", userHandler.Create)
	router.GET("/users/:id", userHandler.GetByID)
	router.PUT("/users/:id", userHandler.Update)
	router.DELETE("/users/:id", userHandler.Delete)
	
	return router, mockUserManager
}

func TestUserHandler_Create(t *testing.T) {
	router, mockUserManager := setupTestRouter(t)
	
	// Setup mock expectations
	expectedUser := &models.User{
		ID:        1,
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		IsBot:     false,
	}
	mockUserManager.On("CreateUser", "testuser", "Test", "User", false).Return(expectedUser, nil)
	
	// Create request
	reqBody := map[string]interface{}{
		"username":   "testuser",
		"first_name": "Test",
		"last_name":  "User",
		"is_bot":     false,
	}
	reqJSON, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assertions
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
	
	mockUserManager.AssertExpectations(t)
}

func TestUserHandler_CreateBot(t *testing.T) {
	router, mockUserManager := setupTestRouter(t)
	
	// Setup mock expectations
	expectedUser := &models.User{
		ID:        2,
		Username:  "testbot",
		FirstName: "Test",
		LastName:  "Bot",
		IsBot:     true,
	}
	mockUserManager.On("CreateUser", "testbot", "Test", "Bot", true).Return(expectedUser, nil)
	
	// Create request
	reqBody := map[string]interface{}{
		"username":   "testbot",
		"first_name": "Test",
		"last_name":  "Bot",
		"is_bot":     true,
	}
	reqJSON, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assertions
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
	
	mockUserManager.AssertExpectations(t)
}

func TestUserHandler_CreateInvalidData(t *testing.T) {
	router, _ := setupTestRouter(t)
	
	// Create request with invalid data
	reqBody := map[string]interface{}{
		"username": "", // Missing required field
	}
	reqJSON, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assertions
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUserHandler_GetAll(t *testing.T) {
	router, mockUserManager := setupTestRouter(t)
	
	// Setup mock expectations
	expectedUsers := []models.User{
		{ID: 1, Username: "user1", FirstName: "User", LastName: "One"},
		{ID: 2, Username: "user2", FirstName: "User", LastName: "Two"},
	}
	mockUserManager.On("GetAllUsers").Return(expectedUsers, nil)
	
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	mockUserManager.AssertExpectations(t)
}

func TestUserHandler_GetByID(t *testing.T) {
	router, mockUserManager := setupTestRouter(t)
	
	// Setup mock expectations
	expectedUser := &models.User{
		ID:        1,
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	}
	mockUserManager.On("GetUser", int64(1)).Return(expectedUser, nil)
	
	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	mockUserManager.AssertExpectations(t)
}

func TestUserHandler_GetByIDNotFound(t *testing.T) {
	router, mockUserManager := setupTestRouter(t)
	
	// Setup mock expectations
	mockUserManager.On("GetUser", int64(999)).Return((*models.User)(nil), &models.UserNotFoundError{})
	
	req, _ := http.NewRequest("GET", "/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assertions
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
	
	mockUserManager.AssertExpectations(t)
}

func TestUserHandler_Update(t *testing.T) {
	router, mockUserManager := setupTestRouter(t)
	
	// Setup mock expectations
	expectedUser := &models.User{
		ID:        1,
		Username:  "testuser",
		FirstName: "Updated",
		LastName:  "Name",
	}
	mockUserManager.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(expectedUser, nil)
	
	// Create request
	reqBody := map[string]interface{}{
		"first_name": "Updated",
		"last_name":  "Name",
	}
	reqJSON, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	mockUserManager.AssertExpectations(t)
}

func TestUserHandler_UpdateNotFound(t *testing.T) {
	router, mockUserManager := setupTestRouter(t)
	
	// Setup mock expectations
	mockUserManager.On("UpdateUser", mock.AnythingOfType("*models.User")).Return((*models.User)(nil), &models.UserNotFoundError{})
	
	// Create request
	reqBody := map[string]interface{}{
		"first_name": "Updated",
		"last_name":  "Name",
	}
	reqJSON, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("PUT", "/users/999", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assertions
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
	
	mockUserManager.AssertExpectations(t)
}

func TestUserHandler_Delete(t *testing.T) {
	router, mockUserManager := setupTestRouter(t)
	
	// Setup mock expectations
	mockUserManager.On("DeleteUser", int64(1)).Return(nil)
	
	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	mockUserManager.AssertExpectations(t)
}

func TestUserHandler_DeleteNotFound(t *testing.T) {
	router, mockUserManager := setupTestRouter(t)
	
	// Setup mock expectations
	mockUserManager.On("DeleteUser", int64(999)).Return(&models.UserNotFoundError{})
	
	req, _ := http.NewRequest("DELETE", "/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assertions
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
	
	mockUserManager.AssertExpectations(t)
}
