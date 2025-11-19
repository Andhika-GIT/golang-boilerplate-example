package repositories

import "app/internal/models"

// IUserRepository interface untuk user repository
type IUserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(id int) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int) error
}
