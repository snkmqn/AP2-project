package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	inventorypb "proto/generated/ecommerce/inventory"
)

type ProductHandler struct {
	Client inventorypb.InventoryServiceClient
}

func NewProductHandler(client inventorypb.InventoryServiceClient) *ProductHandler {
	return &ProductHandler{Client: client}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {

	var req inventorypb.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.Client.CreateProduct(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
