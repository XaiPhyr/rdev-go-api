package server

import (
	"github.com/XaiPhyr/rdev-go-api/internal/dto"
	"github.com/gin-gonic/gin"
)

func responseErr(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{"error": message})
}

func abortErr(ctx *gin.Context, code int, message string) {
	ctx.AbortWithStatusJSON(code, gin.H{"error": message})
}

func parseAuditLog(ctx *gin.Context) dto.AuditLogRequest {
	userID, exists := ctx.Get("userID")
	audit := dto.AuditLogRequest{}
	if exists {
		audit.UserID = userID.(int64)
		audit.Path = ctx.Request.URL.String()
		audit.Action = ctx.Request.Method
		audit.IPAddress = ctx.ClientIP()
		audit.UserAgent = ctx.Request.UserAgent()
	}

	return audit
}
