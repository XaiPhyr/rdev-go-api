package server

import (
	"net/http"
	"rdev-go-api/internal/dto"
	"rdev-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (s *UserHandler) GetUserByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	user, err := s.svc.GetUserByUUID(ctx, uuid)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (s *UserHandler) GetUsers(ctx *gin.Context) {
	var req dto.UserQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	users, count, err := s.svc.GetUsers(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users, "count": count})
}

func (s *UserHandler) UpdateUser(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req dto.UserRequestUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
		return
	}

	err := s.svc.UpdateUser(ctx.Request.Context(), uuid, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s *UserHandler) DeleteUser(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := s.svc.DeleteUser(ctx.Request.Context(), uuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
