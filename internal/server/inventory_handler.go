package server

import (
	"net/http"

	"github.com/XaiPhyr/rdev-go-api/internal/dto"
	"github.com/XaiPhyr/rdev-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	svc *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

func (h *InventoryHandler) GetInventoryByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	inventory, err := h.svc.GetInventoryByUUID(ctx, uuid)

	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if inventory == nil {
		responseErr(ctx, http.StatusNotFound, "inventory not found")
		return
	}

	ctx.JSON(http.StatusOK, inventory)
}

func (h *InventoryHandler) GetInventories(ctx *gin.Context) {
	var req dto.Query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "invalid query parameters")
		return
	}

	inventories, count, err := h.svc.GetInventories(ctx.Request.Context(), req)
	if err != nil {
		responseErr(ctx, http.StatusInternalServerError, "failed to fetch inventories")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": inventories, "count": count})
}

func (h *InventoryHandler) CreateInventory(ctx *gin.Context) {
	var req dto.InventoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.CreateInventory(ctx.Request.Context(), req)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *InventoryHandler) UpdateInventory(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req dto.InventoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	err := h.svc.UpdateInventory(ctx.Request.Context(), uuid, req)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *InventoryHandler) DeleteInventory(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.DeleteInventory(ctx.Request.Context(), uuid)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *InventoryHandler) UpdateInventoryStatus(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.svc.UpdateInventoryStatus(ctx.Request.Context(), uuid)
	if err != nil {
		responseErr(ctx, http.StatusBadRequest, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
