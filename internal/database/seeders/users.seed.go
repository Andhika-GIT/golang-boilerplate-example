package seeders

import (
	"app/internal/models"
	"time"

	"gorm.io/gorm"
)

func timePtr(t time.Time) *time.Time {
	return &t
}

type UserSeeder struct{}

func (s *UserSeeder) GetName() string {
	return "UserSeeder"
}

func (s *UserSeeder) Seed(db *gorm.DB) error {
	users := []models.User{
		{
			ID:              "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			Username:        "linustorvalds",
			Email:           "torvalds@linuxfoundation.org",
			Phone:           "+1-555-0101",
			PasswordHash:    "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
			IsActive:        true,
			EmailVerified:   true,
			EmailVerifiedAt: timePtr(time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)),
			LastLoginAt:     timePtr(time.Date(2024, 1, 20, 14, 30, 0, 0, time.UTC)),
		},
		{
			ID:              "b2c3d4e5-f6a7-8901-bcde-f23456789012",
			Username:        "satyanadella",
			Email:           "satya@microsoft.com",
			Phone:           "+1-555-0102",
			PasswordHash:    "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
			IsActive:        true,
			EmailVerified:   true,
			EmailVerifiedAt: timePtr(time.Date(2024, 1, 10, 8, 0, 0, 0, time.UTC)),
			LastLoginAt:     timePtr(time.Date(2024, 1, 21, 16, 45, 0, 0, time.UTC)),
		},
		{
			ID:              "c3d4e5f6-a7b8-9012-cdef-345678901234",
			Username:        "timcook",
			Email:           "tcook@apple.com",
			Phone:           "+1-555-0103",
			PasswordHash:    "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
			IsActive:        true,
			EmailVerified:   true,
			EmailVerifiedAt: timePtr(time.Date(2024, 1, 5, 11, 0, 0, 0, time.UTC)),
			LastLoginAt:     timePtr(time.Date(2024, 1, 22, 9, 15, 0, 0, time.UTC)),
		},
	}

	for _, user := range users {
		var existingUser models.User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&user).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
