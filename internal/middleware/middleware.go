package middleware

import (
	"net/http"
	"strings"

	"github.com/XaiPhyr/rdev-go-api/internal/auth"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/helpers"

	"github.com/gin-gonic/gin"
)

func AuthRequired(authSvc auth.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			helpers.AbortErr(ctx, http.StatusUnauthorized, "authorization header missing or invalid")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := authSvc.ParseToken(tokenString)
		if err != nil {
			helpers.AbortErr(ctx, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}

func PermissionRequired(authSvc auth.AuthService, requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			helpers.AbortErr(ctx, http.StatusInternalServerError, "user ID not found in context")
			return
		}

		hasAccess, err := authSvc.CanAccess(ctx.Request.Context(), userID.(int64), requiredRole)
		if err != nil {
			helpers.AbortErr(ctx, http.StatusInternalServerError, "failed to check permissions")
			return
		}

		if !hasAccess {
			helpers.AbortErr(ctx, http.StatusForbidden, "insufficient permissions")
			return
		}

		ctx.Next()
	}
}
