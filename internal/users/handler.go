package users

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/helpers"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc UserService
}

func NewUserHandler(svc UserService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetUserByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	user, err := h.svc.GetUserByUUID(ctx, uuid)

	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if user == nil {
		helpers.ResponseErr(ctx, http.StatusNotFound, "user not found")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) GetUsers(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	users, count, err := h.svc.GetUsers(ctx.Request.Context(), req)
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users, "count": count})
}

func (h *Handler) CreateUser(ctx *gin.Context) {
	var req UserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.CreateUser(ctx.Request.Context(), req, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) UpdateUser(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req UserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.UpdateUser(ctx.Request.Context(), uuid, req, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) DeleteUser(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.DeleteUser(ctx.Request.Context(), uuid, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) UpdateUserStatus(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.UpdateUserStatus(ctx.Request.Context(), uuid, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
