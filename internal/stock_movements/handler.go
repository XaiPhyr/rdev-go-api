package stock_movements

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/helpers"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc StockMovementService
}

func NewStockMovementHandler(svc StockMovementService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetStockMovementByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	sm, err := h.svc.GetStockMovementByUUID(ctx, uuid)

	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if sm == nil {
		helpers.ResponseErr(ctx, http.StatusNotFound, "stock movement not found")
		return
	}

	ctx.JSON(http.StatusOK, sm)
}

func (h *Handler) GetStockMovements(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	sm, count, err := h.svc.GetStockMovements(ctx.Request.Context(), req)
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "failed to fetch stock movements")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": sm, "count": count})
}

func (h *Handler) CreateStockMovement(ctx *gin.Context) {
	var req StockMovementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.CreateStockMovement(ctx.Request.Context(), req, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) UpdateStockMovement(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req StockMovementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.UpdateStockMovement(ctx.Request.Context(), uuid, req, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) DeleteStockMovement(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.DeleteStockMovement(ctx.Request.Context(), uuid, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) UpdateStockMovementStatus(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.UpdateStockMovementStatus(ctx.Request.Context(), uuid, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) BulkUpload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err = h.svc.BulkUpload(ctx, file, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusNotFound, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) ProcessBulkUpload(ctx *gin.Context) {
	var req BulkUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.ProcessBulkUpload(ctx.Request.Context(), req.File, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
