package helpers

import (
	"log/slog"
	"os"
	"regexp"
	"testing"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/gin-gonic/gin"
)

func ResponseErr(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{"error": message})
}

func AbortErr(ctx *gin.Context, code int, message string) {
	ctx.AbortWithStatusJSON(code, gin.H{"error": message})
}

func ParseAuditLog(ctx *gin.Context) models.AuditLogRequest {
	userID, exists := ctx.Get("userID")
	audit := models.AuditLogRequest{}
	if exists {
		audit.UserID = userID.(int64)
		audit.Path = ctx.Request.URL.String()
		audit.Action = ctx.Request.Method
		audit.IPAddress = ctx.ClientIP()
		audit.UserAgent = ctx.Request.UserAgent()
	}

	return audit
}

func LogInfo(msg string, args ...any) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info(msg, args...)
}

func CleanSpecialChars(s string) string {
	re := regexp.MustCompile(`[!@#$%^&*()\_\+\=]`)
	clean := re.ReplaceAll([]byte(s), []byte(``))

	return string(clean)
}

func CheckUUID(t testing.TB, uuid string) string {
	t.Helper()

	if uuid == "" {
		t.Error("Expected UUID to be provided")
	}

	return uuid
}
