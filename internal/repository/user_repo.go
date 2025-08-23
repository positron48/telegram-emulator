package repository

import (
	"telegram-emulator/internal/models"

	"gorm.io/gorm"
)

// UserRepository представляет репозиторий для работы с пользователями
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create создает нового пользователя
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID получает пользователя по ID
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername получает пользователя по username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll получает всех пользователей
func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// Update обновляет пользователя
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete удаляет пользователя
func (r *UserRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}

// SetOnlineStatus устанавливает статус онлайн для пользователя
func (r *UserRepository) SetOnlineStatus(id int64, isOnline bool) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_online": isOnline,
		"last_seen": gorm.Expr("CASE WHEN ? = false THEN NOW() ELSE last_seen END", isOnline),
	}).Error
}

// GetOnlineUsers получает всех онлайн пользователей
func (r *UserRepository) GetOnlineUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("is_online = ?", true).Find(&users).Error
	return users, err
}

// GetBots получает всех ботов
func (r *UserRepository) GetBots() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("is_bot = ?", true).Find(&users).Error
	return users, err
}
