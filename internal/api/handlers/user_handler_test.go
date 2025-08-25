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
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(&models.User{}, &models.Chat{}, &models.Message{}, &models.Bot{}, &models.ChatMember{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func setupTestRouter(t *testing.T) *gin.Engine {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	botRepo := repository.NewBotRepository(db)
	userManager := emulator.NewUserManager(userRepo, botRepo)
	userHandler := NewUserHandler(userManager)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/users", userHandler.CreateUser)
	router.GET("/api/users", userHandler.GetAllUsers)
	router.GET("/api/users/:id", userHandler.GetUser)
	router.PUT("/api/users/:id", userHandler.UpdateUser)
	router.DELETE("/api/users/:id", userHandler.DeleteUser)

	return router
}

func TestUserHandler_CreateUser(t *testing.T) {
	router := setupTestRouter(t)

	// Test creating a user
	userData := map[string]interface{}{
		"username":   "testuser",
		"first_name": "Test",
		"last_name":  "User",
		"is_bot":     false,
	}

	jsonData, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["user"])
}

func TestUserHandler_CreateBot(t *testing.T) {
	router := setupTestRouter(t)

	// Test creating a bot
	botData := map[string]interface{}{
		"username":   "testbot",
		"first_name": "Test",
		"last_name":  "Bot",
		"is_bot":     true,
	}

	jsonData, _ := json.Marshal(botData)
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["user"])
}

func TestUserHandler_CreateUserInvalidData(t *testing.T) {
	router := setupTestRouter(t)

	// Test creating a user with invalid data
	userData := map[string]interface{}{
		"username": "", // Empty username should fail
	}

	jsonData, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	router := setupTestRouter(t)

	// Create some users first
	userData1 := map[string]interface{}{
		"username":   "user1",
		"first_name": "User",
		"last_name":  "One",
		"is_bot":     false,
	}

	userData2 := map[string]interface{}{
		"username":   "user2",
		"first_name": "User",
		"last_name":  "Two",
		"is_bot":     false,
	}

	// Create first user
	jsonData1, _ := json.Marshal(userData1)
	req1, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	// Create second user
	jsonData2, _ := json.Marshal(userData2)
	req2, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Get all users
	req, _ := http.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["users"])

	users := response["users"].([]interface{})
	assert.GreaterOrEqual(t, len(users), 2)
}

func TestUserHandler_GetUser(t *testing.T) {
	router := setupTestRouter(t)

	// Create a user first
	userData := map[string]interface{}{
		"username":   "testuser",
		"first_name": "Test",
		"last_name":  "User",
		"is_bot":     false,
	}

	jsonData, _ := json.Marshal(userData)
	req1, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	// Get the created user's ID from the response
	var createResponse map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &createResponse)
	user := createResponse["user"].(map[string]interface{})
	userID := int(user["id"].(float64))

	// Get the user by ID
	req2, _ := http.NewRequest("GET", "/api/users/"+string(rune(userID)), nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["user"])
}

func TestUserHandler_GetUserNotFound(t *testing.T) {
	router := setupTestRouter(t)

	// Try to get non-existent user
	req, _ := http.NewRequest("GET", "/api/users/999999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_UpdateUser(t *testing.T) {
	router := setupTestRouter(t)

	// Create a user first
	userData := map[string]interface{}{
		"username":   "testuser",
		"first_name": "Test",
		"last_name":  "User",
		"is_bot":     false,
	}

	jsonData, _ := json.Marshal(userData)
	req1, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	// Get the created user's ID from the response
	var createResponse map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &createResponse)
	user := createResponse["user"].(map[string]interface{})
	userID := int(user["id"].(float64))

	// Update the user
	updateData := map[string]interface{}{
		"first_name": "Updated",
		"last_name":  "Name",
	}

	jsonUpdateData, _ := json.Marshal(updateData)
	req2, _ := http.NewRequest("PUT", "/api/users/"+string(rune(userID)), bytes.NewBuffer(jsonUpdateData))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["user"])

	updatedUser := response["user"].(map[string]interface{})
	assert.Equal(t, "Updated", updatedUser["first_name"])
	assert.Equal(t, "Name", updatedUser["last_name"])
}

func TestUserHandler_UpdateUserNotFound(t *testing.T) {
	router := setupTestRouter(t)

	// Try to update non-existent user
	updateData := map[string]interface{}{
		"first_name": "Updated",
		"last_name":  "Name",
	}

	jsonUpdateData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/api/users/999999", bytes.NewBuffer(jsonUpdateData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_DeleteUser(t *testing.T) {
	router := setupTestRouter(t)

	// Create a user first
	userData := map[string]interface{}{
		"username":   "testuser",
		"first_name": "Test",
		"last_name":  "User",
		"is_bot":     false,
	}

	jsonData, _ := json.Marshal(userData)
	req1, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	// Get the created user's ID from the response
	var createResponse map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &createResponse)
	user := createResponse["user"].(map[string]interface{})
	userID := int(user["id"].(float64))

	// Delete the user
	req2, _ := http.NewRequest("DELETE", "/api/users/"+string(rune(userID)), nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))

	// Try to get the deleted user
	req3, _ := http.NewRequest("GET", "/api/users/"+string(rune(userID)), nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusNotFound, w3.Code)
}

func TestUserHandler_DeleteUserNotFound(t *testing.T) {
	router := setupTestRouter(t)

	// Try to delete non-existent user
	req, _ := http.NewRequest("DELETE", "/api/users/999999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
