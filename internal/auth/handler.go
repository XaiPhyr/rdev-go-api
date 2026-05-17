package auth

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/helpers"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc AuthService
}

func NewAuthHandler(svc AuthService) *Handler {
	return &Handler{svc: svc}
}

func (s *Handler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "invalid request")
		return
	}

	token, err := s.svc.Login(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusUnauthorized, "invalid credentials")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Handler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "invalid request")
		return
	}

	err := s.svc.Register(ctx.Request.Context(), req)
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "invalid request")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
}
