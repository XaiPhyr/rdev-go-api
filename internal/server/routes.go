package server

import (
	"github.com/XaiPhyr/rdev-go-api/internal/config"
	"github.com/XaiPhyr/rdev-go-api/internal/data"
	"github.com/XaiPhyr/rdev-go-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

// Import Flow Server -> Service -> Data
// "server" package imports "service" (for the middleware).
// "service" package imports "data" (for the repository).
// Never let "data" import "service" or "server".
func Container(r *gin.Engine, db *bun.DB, redis *redis.Client, cfg *config.Config) {
	emailSvc := service.NewEmailService(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.From)

	userRepo := data.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, emailSvc, redis, cfg)

	categoryRepo := data.NewCategoryRepository(db)
	productRepo := data.NewProductRepository(db)
	inventoryRepo := data.NewInventoryRepository(db)

	apiVersion := r.Group("/api/v1")
	setupAuthRoutes(apiVersion, authSvc)
	setupUserRoutes(apiVersion, userRepo, authSvc, emailSvc, redis)
	setupCategoryRoutes(apiVersion, categoryRepo, authSvc, emailSvc, redis)
	setupProductRoutes(apiVersion, productRepo, authSvc, emailSvc, redis)
	setupInventoryRoutes(apiVersion, inventoryRepo, authSvc, emailSvc, redis)
}

func setupAuthRoutes(rg *gin.RouterGroup, authSvc *service.AuthService) {
	authHandler := NewAuthHandler(authSvc)

	rg.POST("/login", authHandler.Login)
	rg.POST("/register", authHandler.Register)
}

func setupUserRoutes(rg *gin.RouterGroup, userRepo *data.UserRepository, authSvc *service.AuthService, emailSvc *service.EmailService, redis *redis.Client) {
	userSvc := service.NewUserService(userRepo, emailSvc, redis)
	userHandler := NewUserHandler(userSvc)

	userRoute := rg.Group("/users")
	userRoute.Use(AuthRequired(authSvc))

	userRoute.GET("", PermissionRequired(authSvc, "users:view"), userHandler.GetUsers)
	userRoute.GET("/:uuid", PermissionRequired(authSvc, "users:view"), userHandler.GetUserByUUID)
	userRoute.PUT("/:uuid", PermissionRequired(authSvc, "users:edit"), userHandler.UpdateUser)
	userRoute.DELETE("/:uuid", PermissionRequired(authSvc, "users:delete"), userHandler.DeleteUser)
	userRoute.POST("/:uuid", PermissionRequired(authSvc, "users:status"), userHandler.UpdateUserStatus)
}

func setupCategoryRoutes(rg *gin.RouterGroup, categoryRepo *data.CategoryRepository, authSvc *service.AuthService, emailSvc *service.EmailService, redis *redis.Client) {
	categorySvc := service.NewCategory(categoryRepo, emailSvc, redis)
	categoryHandler := NewCategoryHandler(categorySvc)

	categoryRoute := rg.Group("/categories")
	categoryRoute.Use(AuthRequired(authSvc))

	categoryRoute.GET("", PermissionRequired(authSvc, "categories:view"), categoryHandler.GetCategories)
	categoryRoute.GET("/:uuid", PermissionRequired(authSvc, "categories:view"), categoryHandler.GetCategoryByUUID)
	categoryRoute.PUT("/:uuid", PermissionRequired(authSvc, "categories:edit"), categoryHandler.UpdateCategory)
	categoryRoute.DELETE("/:uuid", PermissionRequired(authSvc, "categories:delete"), categoryHandler.DeleteCategory)
	categoryRoute.POST("/:uuid", PermissionRequired(authSvc, "categories:status"), categoryHandler.UpdateCategoryStatus)
	categoryRoute.GET("/tree", PermissionRequired(authSvc, "categories:view"), categoryHandler.GetCategoryTree)
}

func setupProductRoutes(rg *gin.RouterGroup, productRepo *data.ProductRepository, authSvc *service.AuthService, emailSvc *service.EmailService, redis *redis.Client) {
	productSvc := service.NewProduct(productRepo, emailSvc, redis)
	productHandler := NewProductHandler(productSvc)

	productRoute := rg.Group("/products")
	productRoute.Use(AuthRequired(authSvc))

	productRoute.GET("/public", productHandler.GetProductsPublic)
	productRoute.GET("", PermissionRequired(authSvc, "products:view"), productHandler.GetProducts)
	productRoute.GET("/:uuid", PermissionRequired(authSvc, "products:view"), productHandler.GetProductByUUID)
	productRoute.PUT("/:uuid", PermissionRequired(authSvc, "products:edit"), productHandler.UpdateProduct)
	productRoute.DELETE("/:uuid", PermissionRequired(authSvc, "products:delete"), productHandler.DeleteProduct)
	productRoute.POST("/:uuid", PermissionRequired(authSvc, "products:status"), productHandler.UpdateProductStatus)
	productRoute.GET("/backoffice", PermissionRequired(authSvc, "products:view"), productHandler.GetProductsBackoffice)
}

func setupInventoryRoutes(rg *gin.RouterGroup, inventoryRepo *data.InventoryRepository, authSvc *service.AuthService, emailSvc *service.EmailService, redis *redis.Client) {
	inventorySvc := service.NewInventory(inventoryRepo, emailSvc, redis)
	inventoryHandler := NewInventoryHandler(inventorySvc)

	inventoryRoute := rg.Group("/inventories")
	inventoryRoute.Use(AuthRequired(authSvc))

	inventoryRoute.GET("", PermissionRequired(authSvc, "inventories:view"), inventoryHandler.GetInventories)
	inventoryRoute.GET("/:uuid", PermissionRequired(authSvc, "inventories:view"), inventoryHandler.GetInventoryByUUID)
	inventoryRoute.PUT("/:uuid", PermissionRequired(authSvc, "inventories:edit"), inventoryHandler.UpdateInventory)
	inventoryRoute.DELETE("/:uuid", PermissionRequired(authSvc, "inventories:delete"), inventoryHandler.DeleteInventory)
	inventoryRoute.POST("/:uuid", PermissionRequired(authSvc, "inventories:status"), inventoryHandler.UpdateInventoryStatus)
}
