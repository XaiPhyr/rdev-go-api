package auth

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/XaiPhyr/rdev-go-api/internal/config"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api/internal/users"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, req RegisterRequest) error
	GenerateToken(userID int64) (string, error)
	ParseToken(token string) (int64, error)
	CanAccess(ctx context.Context, userID int64, requiredRole string) (bool, error)
}

type service struct {
	r     users.UserRepository
	es    email.EmailService
	redis *redis.Client
	c     *config.Config
}

func NewAuthService(r users.UserRepository, es email.EmailService, redis *redis.Client, c *config.Config) *service {
	return &service{r: r, es: es, redis: redis, c: c}
}

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
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

func (s *service) Register(ctx context.Context, req RegisterRequest) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := &users.User{
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

func (s *service) GenerateToken(userID int64) (string, error) {
	jwtKey := []byte(s.c.JWTSecretKey)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func (s *service) ParseToken(token string) (int64, error) {
	jwtKey := []byte(s.c.JWTSecretKey)

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

func (s *service) CanAccess(ctx context.Context, userID int64, requiredRole string) (bool, error) {
	// when updating user_roles and user_groups delete cache after
	cacheKey := fmt.Sprintf("user:perms:%d", userID)

	existCount, err := s.redis.Exists(ctx, cacheKey).Result()
	if err == nil && existCount > 0 {
		isSuperAdmin, _ := s.redis.SIsMember(ctx, cacheKey, "super_admin").Result()
		if isSuperAdmin {
			return true, nil
		}

		hasRole, _ := s.redis.SIsMember(ctx, cacheKey, requiredRole).Result()
		return hasRole, nil
	}

	allPerms, err := s.r.CheckUserPermission(ctx, userID, requiredRole)
	if err != nil {
		log.Println(fmt.Errorf("user permission error: %w", err))
		return false, err
	}

	if len(allPerms) > 0 {
		pipe := s.redis.Pipeline()
		pipe.SAdd(ctx, cacheKey, allPerms)
		pipe.Expire(ctx, cacheKey, 1*time.Hour)
		_, err := pipe.Exec(ctx)
		if err != nil {
			log.Printf("failed to update redis: %v", err)
		}
	} else {
		s.redis.SAdd(ctx, cacheKey, "NONE")
		s.redis.Expire(ctx, cacheKey, 5*time.Minute)
	}

	return slices.Contains(allPerms, requiredRole) || slices.Contains(allPerms, "super_admin"), nil
}
