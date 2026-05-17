package products

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/helpers"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc ProductService
}

func NewProductHandler(svc ProductService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetProductByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	product, err := h.svc.GetProductByUUID(ctx, uuid)

	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if product == nil {
		helpers.ResponseErr(ctx, http.StatusNotFound, "product not found")
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (h *Handler) GetProducts(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	products, count, err := h.svc.GetProducts(ctx.Request.Context(), req)
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products, "count": count})
}

func (h *Handler) GetProductsPublic(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	products, count, err := h.svc.GetProductsPublic(ctx.Request.Context(), req)
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products, "count": count})
}

func (h *Handler) GetProductsBackoffice(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	products, count, err := h.svc.GetProductsBackoffice(ctx.Request.Context(), req)
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products, "count": count})
}

func (h *Handler) CreateProduct(ctx *gin.Context) {
	var req ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.CreateProduct(ctx, req, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func (h *Handler) UpdateProduct(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.UpdateProduct(ctx.Request.Context(), uuid, req, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) DeleteProduct(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.DeleteProduct(ctx.Request.Context(), uuid, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) UpdateProductStatus(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.UpdateProductStatus(ctx.Request.Context(), uuid, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
