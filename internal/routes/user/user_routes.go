// user_routes.go placeholder
package user

import (
	"app/internal/controllers/user"

	"github.com/gin-gonic/gin"
)

// IUserRoutes defines the interface for user route registration.
type IUserRoutes interface {
	Run()
}

// UserRoutes handles the registration of all user-related routes.
type UserRoutes struct {
	router         *gin.Engine
	userController user.IUserController
}

// NewUserRoutes creates a new instance of UserRoutes.
// It receives a Gin router and a user controller that will handle the incoming requests.
func NewUserRoutes(
	router *gin.Engine,
	userController user.IUserController,
) IUserRoutes {
	return &UserRoutes{
		router:         router,
		userController: userController,
	}
}

// Run registers all user routes under the prefix /api/v1/users.
func (ur *UserRoutes) Run() {
	// Group routes with the prefix /api/v1/users
	userGroup := ur.router.Group("/api/v1/users")
	{
		userGroup.GET("", ur.userController.List)          // GET /api/v1/users
		userGroup.GET("/:id", ur.userController.GetByID)   // GET /api/v1/users/:id
		userGroup.POST("", ur.userController.Create)       // POST /api/v1/users
		userGroup.PUT("/:id", ur.userController.Update)    // PUT /api/v1/users/:id
		userGroup.DELETE("/:id", ur.userController.Delete) // DELETE /api/v1/users/:id
	}
}
