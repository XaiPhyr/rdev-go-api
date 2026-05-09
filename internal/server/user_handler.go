package server

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetUserByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	user, err := h.svc.GetUserByUUID(ctx, uuid)

	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if user == nil {
		responseErr(ctx, http.StatusNotFound, "user not found")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUsers(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	users, count, err := h.svc.GetUsers(ctx.Request.Context(), req)
	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users, "count": count})
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req dto.UserRequestUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.UpdateUser(ctx.Request.Context(), uuid, req)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.DeleteUser(ctx.Request.Context(), uuid)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *UserHandler) UpdateUserStatus(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.UpdateUserStatus(ctx.Request.Context(), uuid)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
