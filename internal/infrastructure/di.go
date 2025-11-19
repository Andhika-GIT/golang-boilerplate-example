// di.go placeholder
package infrastructure

import (
	"app/internal/config"
	"app/internal/controllers"
	"app/internal/repositories"
	"app/internal/services"
	"log"
)

// ✅ HIGH-LEVEL: Business logic dependencies
type Dependencies struct {
	UserService        services.IUserService
	ControllerRegistry controllers.IControllerRegistry
}

// ✅ HIGH-LEVEL: Wire business logic
func InitializeDependencies(
	conn *Connections,
	cfg *config.Config,
) (*Dependencies, error) {
	log.Println("Initializing dependencies...")

	// 1. Repositories
	userRepo := repositories.NewUserRepository(conn.DB)
	log.Println("✅ Repositories initialized")

	// 2. Services
	userService := services.NewUserService(userRepo)
	log.Println("✅ Services initialized")

	// 3. Controllers
	controllerRegistry := controllers.NewControllerRegistry(
		userService,
	)
	log.Println("✅ Controllers initialized")

	return &Dependencies{
		UserService:        userService,
		ControllerRegistry: controllerRegistry,
	}, nil
}
