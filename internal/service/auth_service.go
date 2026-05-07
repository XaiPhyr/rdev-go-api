package service

import (
	"context"
	"log"
	"rdev-go-api/internal/config"
	"rdev-go-api/internal/data"
	"rdev-go-api/internal/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	r  *data.UserRepository
	es *EmailService
	c  *config.Config
}

func NewAuthService(r *data.UserRepository, es *EmailService, c *config.Config) *AuthService {
	return &AuthService{r: r, es: es, c: c}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.r.GetUserByUsernameOrEmail(ctx, username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	return s.GenerateToken(user.ID)
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := &data.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(passwordHash),
	}

	err = s.r.CreateUser(ctx, user)
	if err == nil {
		go func() {
			if err := s.es.SendEmail(req.Email); err != nil {
				log.Printf("Failed to send email: %v", err)
			}
		}()
	}

	return err
}

func (s *AuthService) GenerateToken(userID int64) (string, error) {
	jwtKey := []byte(s.c.JWT.SecretKey)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func (s *AuthService) ParseToken(token string) (int64, error) {
	jwtKey := []byte(s.c.JWT.SecretKey)

	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	if err != nil || !t.Valid {
		return 0, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["user_id"].(float64); ok {
			return int64(userID), nil
		}
	}

	return 0, jwt.ErrTokenInvalidClaims
}

func (s *AuthService) CanAccess(ctx context.Context, userID int64, requiredRole string) (bool, error) {
	return s.r.CheckUserPermission(ctx, userID, requiredRole)
}
