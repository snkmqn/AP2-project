package handlers

import (
	"context"
	"net/http"

	orderpb "proto/generated/ecommerce/order"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Client orderpb.OrderServiceClient
}

func NewOrderHandler(client orderpb.OrderServiceClient) *OrderHandler {
	return &OrderHandler{Client: client}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.Client.CreateOrder(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
