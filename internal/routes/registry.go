// registry.go placeholder
package routes

import (
	"app/internal/controllers"
	userRoutes "app/internal/routes/user"

	"github.com/gin-gonic/gin"
)

// IRouteRegistry defines the interface for the route registry.
type IRouteRegistry interface {
	Serve(addr string) error
}

// RouteRegistry holds the Gin router and the controller registry.
type RouteRegistry struct {
	router             *gin.Engine
	controllerRegistry controllers.IControllerRegistry
	userRoutes         userRoutes.IUserRoutes
}

// NewRouteRegistry creates a new instance of RouteRegistry.
func NewRouteRegistry(
	router *gin.Engine,
	controllerRegistry controllers.IControllerRegistry,
) IRouteRegistry {
	registry := &RouteRegistry{
		router:             router,
		controllerRegistry: controllerRegistry,
	}

	// Initialize user routes (controllers already prepared by the registry)
	registry.userRoutes = userRoutes.NewUserRoutes(
		registry.router,
		controllerRegistry.GetUserController(),
	)

	// Register all routes
	registry.setupRoutes()

	return registry
}

// setupRoutes registers all API routes in the application.
func (rr *RouteRegistry) setupRoutes() {
	// Health check endpoint
	rr.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Register user routes
	rr.userRoutes.Run()

	// Additional route modules can be added here
	// rr.productRoutes.Run()
	// rr.orderRoutes.Run()
}

// Serve starts the HTTP server on the specified address.
func (rr *RouteRegistry) Serve(addr string) error {
	return rr.router.Run(addr)
}
