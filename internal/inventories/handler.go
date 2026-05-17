package inventories

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/shared/helpers"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc InventoryService
}

func NewInventoryHandler(svc InventoryService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetInventoryByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	inventory, err := h.svc.GetInventoryByUUID(ctx, uuid)

	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if inventory == nil {
		helpers.ResponseErr(ctx, http.StatusNotFound, "inventory not found")
		return
	}

	ctx.JSON(http.StatusOK, inventory)
}

func (h *Handler) GetInventories(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	inventories, count, err := h.svc.GetInventories(ctx.Request.Context(), req)
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusInternalServerError, "failed to fetch inventories")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": inventories, "count": count})
}

func (h *Handler) CreateInventory(ctx *gin.Context) {
	var req InventoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.CreateInventory(ctx.Request.Context(), req, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) UpdateInventory(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req InventoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.UpdateInventory(ctx.Request.Context(), uuid, req, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) DeleteInventory(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.DeleteInventory(ctx.Request.Context(), uuid, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) UpdateInventoryStatus(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.UpdateInventoryStatus(ctx.Request.Context(), uuid, helpers.ParseAuditLog(ctx))
	if err != nil {
		helpers.ResponseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
