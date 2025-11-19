package controllers

import (
	userController "app/internal/controllers/user"
	"app/internal/services"
)

// IControllerRegistry defines the interface for the controller registry.
// It exposes getter methods for each controller in the application.
type IControllerRegistry interface {
	GetUserController() userController.IUserController
	// Additional controllers can be registered here.
	// GetProductController() productController.IProductController
}

// ControllerRegistry stores all controller instances used by the application.
// This helps centralize dependency injection and makes the controllers easier to manage.
type ControllerRegistry struct {
	userController userController.IUserController
	// productController productController.IProductController
}

// NewControllerRegistry creates a new instance of ControllerRegistry.
// All required services should already be initialized in main.go and passed here.
func NewControllerRegistry(
	userService services.IUserService,
	// productService services.IProductService,
) IControllerRegistry {
	return &ControllerRegistry{
		userController: userController.NewUserController(userService),
		// productController: productController.NewProductController(productService),
	}
}

// GetUserController returns the user controller instance.
func (cr *ControllerRegistry) GetUserController() userController.IUserController {
	return cr.userController
}
