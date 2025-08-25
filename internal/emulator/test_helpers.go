package emulator

import (
	"testing"

	"telegram-emulator/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB создает тестовую базу данных в памяти
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Bot{}, &models.User{}, &models.Chat{}, &models.Message{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
