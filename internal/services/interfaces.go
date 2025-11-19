// interfaces.go placeholder
package services

import "app/internal/models"

// IUserService interface untuk user service
type IUserService interface {
	GetAll() ([]models.User, error)
	GetByID(id int) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int) error
}
