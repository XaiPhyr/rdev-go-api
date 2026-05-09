package server

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) GetCategoryByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	category, err := h.svc.GetCategoryByUUID(ctx, uuid)

	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if category == nil {
		responseErr(ctx, http.StatusNotFound, "category not found")
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) GetCategories(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	categories, count, err := h.svc.GetCategories(ctx.Request.Context(), req)
	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "failed to fetch categories")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": categories, "count": count})
}

func (h *CategoryHandler) GetCategoryTree(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	categories, err := h.svc.GetCategoryTree(ctx.Request.Context(), req)
	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "failed to fetch categories")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req dto.CategoryRequestUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.UpdateCategory(ctx.Request.Context(), uuid, req)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.DeleteCategory(ctx.Request.Context(), uuid)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *CategoryHandler) UpdateCategoryStatus(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.UpdateCategoryStatus(ctx.Request.Context(), uuid)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
