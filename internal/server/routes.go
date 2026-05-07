package server

import (
	"rdev-go-api/internal/config"
	"rdev-go-api/internal/data"
	"rdev-go-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func RegisterRouters(r *gin.Engine, db *bun.DB, cfg *config.Config) {
	// Import Flow Server -> Service -> Data
	// "server" package imports "service" (for the middleware).
	// "service" package imports "data" (for the repository).
	// Never let "data" import "service" or "server".

	emailSvc := service.NewEmailService(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.From)

	userRepo := data.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, emailSvc, cfg)

	apiVersion := r.Group("/api/v1")
	setupAuthRoutes(apiVersion, authSvc)
	setupUserRoutes(apiVersion, userRepo, authSvc)
}

func setupAuthRoutes(rg *gin.RouterGroup, authSvc *service.AuthService) {
	authHandler := NewAuthHandler(authSvc)

	rg.POST("/login", authHandler.Login)
	rg.POST("/register", authHandler.Register)
}

func setupUserRoutes(rg *gin.RouterGroup, userRepo *data.UserRepository, authSvc *service.AuthService) {
	userSvc := service.NewUserService(userRepo)
	userHandler := NewUserHandler(userSvc)

	userRoute := rg.Group("/users")
	userRoute.Use(AuthRequired(authSvc))
	{
		userRoute.GET("", PermissionRequired(authSvc, "users:view"), userHandler.GetUsers)
		userRoute.GET("/:uuid", PermissionRequired(authSvc, "users:view"), userHandler.GetUserByUUID)
		userRoute.PUT("/:uuid", PermissionRequired(authSvc, "users:edit"), userHandler.UpdateUser)
		userRoute.DELETE("/:uuid", PermissionRequired(authSvc, "users:delete"), userHandler.DeleteUser)
	}
}
