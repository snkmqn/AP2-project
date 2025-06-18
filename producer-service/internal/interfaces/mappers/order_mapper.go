package mappers

import (
	"producer-service/internal/core/models"
	orderpb "proto/generated/ecommerce/order"
	"time"
)

func ToProtoOrderItem(item models.OrderItem) *orderpb.OrderItem {
	return &orderpb.OrderItem{
		ProductId:    item.ProductID,
		Quantity:     int32(item.Quantity),
		PricePerUnit: item.PricePerUnit,
	}
}

func ToProtoOrder(order *models.Order) *orderpb.Order {
	var items []*orderpb.OrderItem
	for _, item := range order.Items {
		items = append(items, ToProtoOrderItem(item))
	}

	return &orderpb.Order{
		Id:         order.ID,
		OrderId:    order.OrderID,
		UserId:     order.UserID,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		Items:      items,
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
	}
}

func ToProtoOrders(orders []*models.Order) []*orderpb.Order {
	var protoOrders []*orderpb.Order
	for _, order := range orders {
		protoOrders = append(protoOrders, ToProtoOrder(order))
	}
	return protoOrders
}
