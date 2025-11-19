// user_service.go placeholder
package services

import (
	"app/internal/models"
	"app/internal/repositories"
)

type UserService struct {
	repo repositories.IUserRepository
}

// NewUserService constructor untuk inisialisasi service
func NewUserService(repo repositories.IUserRepository) IUserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) Create(user *models.User) error {
	return s.repo.Create(user)
}

func (s *UserService) Update(user *models.User) error {
	return s.repo.Update(user)
}

func (s *UserService) Delete(id int) error {
	return s.repo.Delete(id)
}
