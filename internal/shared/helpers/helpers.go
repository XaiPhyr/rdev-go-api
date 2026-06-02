package helpers

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func ValidateStruct(s interface{}) error {
	validate := validator.New()

	return validate.Struct(s)
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

func SaveImage(path, file string, data []byte) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return fmt.Errorf("could not create directory: %w", err)
		}
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
