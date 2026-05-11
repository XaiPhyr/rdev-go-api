package server

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type StockMovementHandler struct {
	svc service.StockMovementService
}

func NewStockMovementHandler(svc service.StockMovementService) *StockMovementHandler {
	return &StockMovementHandler{svc: svc}
}

func (h *StockMovementHandler) GetStockMovementByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	sm, err := h.svc.GetStockMovementByUUID(ctx, uuid)

	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if sm == nil {
		responseErr(ctx, http.StatusNotFound, "stock movement not found")
		return
	}

	ctx.JSON(http.StatusOK, sm)
}

func (h *StockMovementHandler) GetStockMovements(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	sm, count, err := h.svc.GetStockMovements(ctx.Request.Context(), req)
	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "failed to fetch stock movements")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": sm, "count": count})
}

func (h *StockMovementHandler) CreateStockMovement(ctx *gin.Context) {
	var req dto.StockMovementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.CreateStockMovement(ctx.Request.Context(), req)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *StockMovementHandler) UpdateStockMovement(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req dto.StockMovementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.UpdateStockMovement(ctx.Request.Context(), uuid, req)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *StockMovementHandler) DeleteStockMovement(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.DeleteStockMovement(ctx.Request.Context(), uuid)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *StockMovementHandler) UpdateStockMovementStatus(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.UpdateStockMovementStatus(ctx.Request.Context(), uuid)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
