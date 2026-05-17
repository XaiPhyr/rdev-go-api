package categories

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/helpers"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc CategoryService
}

func NewCategoryHandler(svc CategoryService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetCategoryByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	category, err := h.svc.GetCategoryByUUID(ctx, uuid)

	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if category == nil {
		helpers.ResponseErr(ctx, http.StatusNotFound, "category not found")
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (h *Handler) GetCategories(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	categories, count, err := h.svc.GetCategories(ctx.Request.Context(), req)
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "failed to fetch categories")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": categories, "count": count})
}

func (h *Handler) GetCategoryTree(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	categories, err := h.svc.GetCategoryTree(ctx.Request.Context(), req)
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "failed to fetch categories")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *Handler) CreateCategory(ctx *gin.Context) {
	var req CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.CreateCategory(ctx.Request.Context(), req, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) UpdateCategory(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.UpdateCategory(ctx.Request.Context(), uuid, req, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) DeleteCategory(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.DeleteCategory(ctx.Request.Context(), uuid, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) UpdateCategoryStatus(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.UpdateCategoryStatus(ctx.Request.Context(), uuid, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
