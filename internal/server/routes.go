package server

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api/internal/auth"
	"github.com/XaiPhyr/rdev-go-api/internal/categories"
	"github.com/XaiPhyr/rdev-go-api/internal/config"
	"github.com/XaiPhyr/rdev-go-api/internal/inventories"
	"github.com/XaiPhyr/rdev-go-api/internal/middleware"
	"github.com/XaiPhyr/rdev-go-api/internal/products"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api/internal/stock_movements"
	"github.com/XaiPhyr/rdev-go-api/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

// Import Flow Server -> Service -> Data
// "server" package imports "service" (for the middleware).
// "service" package imports "data" (for the repository).
// Never let "data" import "service" or "server".
func Container(r *gin.Engine, db *bun.DB, redis *redis.Client, cfg *config.Config) {
	emailSvc := email.NewEmailService(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.From)

	auditLogRepo := audit_logs.NewAuditLogRepository(db)
	auditLogSvc := audit_logs.NewAuditLogService(auditLogRepo)
	go auditLogSvc.QueAuditLog(context.Background())

	userRepo := users.NewUserRepository(db)
	authSvc := auth.NewAuthService(userRepo, emailSvc, redis, cfg)

	categoryRepo := categories.NewCategoryRepository(db)
	productRepo := products.NewProductRepository(db)
	inventoryRepo := inventories.NewInventoryRepository(db)
	stockMovementRepo := stock_movements.NewStockMovementRepository(db)

	apiVersion := r.Group("/api/v1")
	setupAuthRoutes(apiVersion, authSvc)
	setupUserRoutes(apiVersion, userRepo, authSvc, emailSvc, redis, auditLogSvc)
	setupCategoryRoutes(apiVersion, categoryRepo, authSvc, emailSvc, redis, auditLogSvc)
	setupProductRoutes(apiVersion, productRepo, authSvc, emailSvc, redis, auditLogSvc)
	setupInventoryRoutes(apiVersion, inventoryRepo, authSvc, emailSvc, redis, auditLogSvc)
	setupStockMovementRoutes(apiVersion, stockMovementRepo, authSvc, emailSvc, redis, auditLogSvc)
}

func setupAuthRoutes(rg *gin.RouterGroup, authSvc auth.AuthService) {
	authHandler := auth.NewAuthHandler(authSvc)

	rg.POST("/login", authHandler.Login)
	rg.POST("/register", authHandler.Register)
}

func setupUserRoutes(rg *gin.RouterGroup, userRepo users.UserRepository, authSvc auth.AuthService, emailSvc email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) {
	userSvc := users.NewUserService(userRepo, emailSvc, redis, auditLog)
	userHandler := users.NewUserHandler(userSvc)

	userRoute := rg.Group("/users")
	userRoute.Use(middleware.AuthRequired(authSvc))

	userRoute.GET("", middleware.PermissionRequired(authSvc, "users:view"), userHandler.GetUsers)
	userRoute.GET("/:uuid", middleware.PermissionRequired(authSvc, "users:view"), userHandler.GetUserByUUID)
	userRoute.POST("", middleware.PermissionRequired(authSvc, "users:create"), userHandler.CreateUser)
	userRoute.PUT("/:uuid", middleware.PermissionRequired(authSvc, "users:edit"), userHandler.UpdateUser)
	userRoute.DELETE("/:uuid", middleware.PermissionRequired(authSvc, "users:delete"), userHandler.DeleteUser)
	userRoute.POST("/updatestatus/:uuid", middleware.PermissionRequired(authSvc, "users:status"), userHandler.UpdateUserStatus)
}

func setupCategoryRoutes(rg *gin.RouterGroup, categoryRepo categories.CategoryRepository, authSvc auth.AuthService, emailSvc email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) {
	categorySvc := categories.NewCategoryService(categoryRepo, emailSvc, redis, auditLog)
	categoryHandler := categories.NewCategoryHandler(categorySvc)

	categoryRoute := rg.Group("/categories")
	categoryRoute.Use(middleware.AuthRequired(authSvc))

	categoryRoute.GET("", middleware.PermissionRequired(authSvc, "categories:view"), categoryHandler.GetCategories)
	categoryRoute.GET("/:uuid", middleware.PermissionRequired(authSvc, "categories:view"), categoryHandler.GetCategoryByUUID)
	categoryRoute.POST("", middleware.PermissionRequired(authSvc, "categories:create"), categoryHandler.CreateCategory)
	categoryRoute.PUT("/:uuid", middleware.PermissionRequired(authSvc, "categories:edit"), categoryHandler.UpdateCategory)
	categoryRoute.DELETE("/:uuid", middleware.PermissionRequired(authSvc, "categories:delete"), categoryHandler.DeleteCategory)
	categoryRoute.POST("/updatestatus/:uuid", middleware.PermissionRequired(authSvc, "categories:status"), categoryHandler.UpdateCategoryStatus)
	categoryRoute.GET("/tree", middleware.PermissionRequired(authSvc, "categories:view"), categoryHandler.GetCategoryTree)
}

func setupProductRoutes(rg *gin.RouterGroup, productRepo products.ProductRepository, authSvc auth.AuthService, emailSvc email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) {
	productSvc := products.NewProductService(productRepo, emailSvc, redis, auditLog)
	productHandler := products.NewProductHandler(productSvc)

	productRoute := rg.Group("/products")
	productRoute.Use(middleware.AuthRequired(authSvc))

	productRoute.GET("/public", productHandler.GetProductsPublic)
	productRoute.GET("", middleware.PermissionRequired(authSvc, "products:view"), productHandler.GetProducts)
	productRoute.GET("/:uuid", middleware.PermissionRequired(authSvc, "products:view"), productHandler.GetProductByUUID)
	productRoute.POST("", middleware.PermissionRequired(authSvc, "products:create"), productHandler.CreateProduct)
	productRoute.PUT("/:uuid", middleware.PermissionRequired(authSvc, "products:edit"), productHandler.UpdateProduct)
	productRoute.DELETE("/:uuid", middleware.PermissionRequired(authSvc, "products:delete"), productHandler.DeleteProduct)
	productRoute.POST("/updatestatus/:uuid", middleware.PermissionRequired(authSvc, "products:status"), productHandler.UpdateProductStatus)
	productRoute.GET("/backoffice", middleware.PermissionRequired(authSvc, "products:view"), productHandler.GetProductsBackoffice)
}

func setupInventoryRoutes(rg *gin.RouterGroup, inventoryRepo inventories.InventoryRepository, authSvc auth.AuthService, emailSvc email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) {
	inventorySvc := inventories.NewInventoryService(inventoryRepo, emailSvc, redis, auditLog)
	inventoryHandler := inventories.NewInventoryHandler(inventorySvc)

	inventoryRoute := rg.Group("/inventories")
	inventoryRoute.Use(middleware.AuthRequired(authSvc))

	inventoryRoute.GET("", middleware.PermissionRequired(authSvc, "inventories:view"), inventoryHandler.GetInventories)
	inventoryRoute.GET("/:uuid", middleware.PermissionRequired(authSvc, "inventories:view"), inventoryHandler.GetInventoryByUUID)
	inventoryRoute.POST("", middleware.PermissionRequired(authSvc, "inventories:create"), inventoryHandler.CreateInventory)
	inventoryRoute.PUT("/:uuid", middleware.PermissionRequired(authSvc, "inventories:edit"), inventoryHandler.UpdateInventory)
	inventoryRoute.DELETE("/:uuid", middleware.PermissionRequired(authSvc, "inventories:delete"), inventoryHandler.DeleteInventory)
	inventoryRoute.POST("/updatestatus/:uuid", middleware.PermissionRequired(authSvc, "inventories:status"), inventoryHandler.UpdateInventoryStatus)
}

func setupStockMovementRoutes(rg *gin.RouterGroup, stockMovementRepo stock_movements.StockMovementRepository, authSvc auth.AuthService, emailSvc email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) {
	stockMovementSvc := stock_movements.NewStockMovementService(stockMovementRepo, emailSvc, redis, auditLog)
	stockMovementHandler := stock_movements.NewStockMovementHandler(stockMovementSvc)

	stockMovementRoute := rg.Group("/stock_movements")
	stockMovementRoute.Use(middleware.AuthRequired(authSvc))

	stockMovementRoute.GET("", middleware.PermissionRequired(authSvc, "stock_movements:view"), stockMovementHandler.GetStockMovements)
	stockMovementRoute.GET("/:uuid", middleware.PermissionRequired(authSvc, "stock_movements:view"), stockMovementHandler.GetStockMovementByUUID)
	stockMovementRoute.POST("", middleware.PermissionRequired(authSvc, "stock_movements:create"), stockMovementHandler.CreateStockMovement)
	stockMovementRoute.PUT("/:uuid", middleware.PermissionRequired(authSvc, "stock_movements:edit"), stockMovementHandler.UpdateStockMovement)
	stockMovementRoute.DELETE("/:uuid", middleware.PermissionRequired(authSvc, "stock_movements:delete"), stockMovementHandler.DeleteStockMovement)
	stockMovementRoute.POST("/updatestatus/:uuid", middleware.PermissionRequired(authSvc, "stock_movements:status"), stockMovementHandler.UpdateStockMovementStatus)
	stockMovementRoute.POST("/bulkupload", middleware.PermissionRequired(authSvc, "stock_movements:upload"), stockMovementHandler.BulkUpload)
	stockMovementRoute.POST("/processbulkupload", middleware.PermissionRequired(authSvc, "stock_movements:process"), stockMovementHandler.ProcessBulkUpload)
}
