package http

import (
	"github.com/Putthakun/face-recognition-api-go/internal/delivery/http/handler"
	"github.com/Putthakun/face-recognition-api-go/internal/delivery/http/middleware"
	"github.com/Putthakun/face-recognition-api-go/pkg/config"
	"github.com/Putthakun/face-recognition-api-go/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth        *handler.AuthHandler
	Employee    *handler.EmployeeHandler
	Camera      *handler.CameraHandler
	Transaction *handler.TransactionHandler
}

func NewRouter(cfg *config.Config, jwtService jwt.Service, h Handlers) *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(corsMiddleware(cfg.AllowedOrigins))

	authMiddleware := middleware.Auth(jwtService)
	adminOnly := middleware.RequireRole("Admin")
	adminOrSupervisor := middleware.RequireRole("Admin", "Supervisor")

	api := r.Group("/api")
	{
		// Auth — public
		api.POST("/auth/login", h.Auth.Login)

		// Transactions — POST is internal (no auth), GET requires role
		api.POST("/transactions", h.Transaction.Create)
		api.GET("/transactions", authMiddleware, adminOrSupervisor, h.Transaction.GetAll)

		// Admin routes
		admin := api.Group("/admin", authMiddleware, adminOnly)
		{
			admin.GET("/employees", h.Employee.GetAll)
			admin.POST("/employees", h.Employee.Create)
			admin.PUT("/employees/:empId", h.Employee.Update)
			admin.DELETE("/employees/:empId", h.Employee.Delete)

			admin.GET("/cameras", h.Camera.GetAll)
			admin.POST("/cameras", h.Camera.Create)
			admin.PUT("/cameras/:cameraId", h.Camera.Update)
			admin.DELETE("/cameras/:cameraId", h.Camera.Delete)
		}
	}

	return r
}

func corsMiddleware(origins []string) gin.HandlerFunc {
	originSet := make(map[string]bool, len(origins))
	for _, o := range origins {
		originSet[o] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if originSet[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
