package server

import (
	"net/http"
	"strings"

	"github.com/XaiPhyr/rdev-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthRequired(authSvc *service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			abortErr(ctx, http.StatusUnauthorized, "authorization header missing or invalid")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := authSvc.ParseToken(tokenString)
		if err != nil {
			abortErr(ctx, http.StatusUnauthorized, "invalid token")
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
			abortErr(ctx, http.StatusInternalServerError, "user ID not found in context")
			return
		}

		hasAccess, err := authSvc.CanAccess(ctx.Request.Context(), userID.(int64), requiredRole)
		if err != nil {
			abortErr(ctx, http.StatusInternalServerError, "failed to check permissions")
			return
		}

		if !hasAccess {
			abortErr(ctx, http.StatusForbidden, "insufficient permissions")
			return
		}

		ctx.Next()
	}
}
