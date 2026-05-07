package server

import (
	"net/http"
	"rdev-go-api/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired(authSvc *service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing or invalid"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := authSvc.ParseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}

func PermissionRequired(authSvc *service.AuthService, requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
			return
		}

		hasAccess, err := authSvc.CanAccess(ctx.Request.Context(), userID.(int64), requiredRole)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to check permissions"})
			return
		}

		if !hasAccess {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		ctx.Next()
	}
}
