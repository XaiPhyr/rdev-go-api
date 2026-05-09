package server

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) GetProductByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	product, err := h.svc.GetProductByUUID(ctx, uuid)

	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if product == nil {
		responseErr(ctx, http.StatusNotFound, "product not found")
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetProducts(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	products, count, err := h.svc.GetProducts(ctx.Request.Context(), req)
	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products, "count": count})
}

func (h *ProductHandler) GetProductsPublic(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	products, count, err := h.svc.GetProductsPublic(ctx.Request.Context(), req)
	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products, "count": count})
}

func (h *ProductHandler) GetProductsBackoffice(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	products, count, err := h.svc.GetProductsBackoffice(ctx.Request.Context(), req)
	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products, "count": count})
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var req dto.ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.CreateProduct(ctx, req)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req dto.ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.UpdateProduct(ctx.Request.Context(), uuid, req)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.DeleteProduct(ctx.Request.Context(), uuid)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *ProductHandler) UpdateProductStatus(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.UpdateProductStatus(ctx.Request.Context(), uuid)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
