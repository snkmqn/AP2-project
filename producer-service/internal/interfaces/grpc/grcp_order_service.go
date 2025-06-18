package services

import (
	"context"
	"producer-service/internal/infrastructure/cache"
	logger "producer-service/internal/interfaces/logger"
	"producer-service/internal/interfaces/mappers"
	"producer-service/internal/usecases/services"
	"producer-service/internal/usecases/validators"
	"proto/generated/ecommerce/order"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGrpcServer struct {
	orderpb.UnimplementedOrderServiceServer
	orderService *services.OrderService
	natsClient   *nats.Conn
	logger       logger.Logger
	cache        cache.CacheService
}

func NewOrderGrpcServer(logger logger.Logger, orderService *services.OrderService, natsClient *nats.Conn, cache cache.CacheService) *OrderGrpcServer {
	return &OrderGrpcServer{
		orderService: orderService,
		natsClient:   natsClient,
		logger:       logger,
		cache:        cache,
	}
}

func (s *OrderGrpcServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	if err := validators.ValidateOrderStatus(req.GetStatus()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	order, err := s.orderService.CreateOrderFromProto(ctx, req)
	if err != nil {
		return nil, err
	}

	orderData := mappers.ToProtoOrder(order)

	orderJson, err := json.Marshal(orderData)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to marshal order data")
	}

	err = s.natsClient.Publish("order.created", orderJson)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to publish event to NATS")
	}

	return &orderpb.CreateOrderResponse{
		Order: mappers.ToProtoOrder(order),
	}, nil
}

func (s *OrderGrpcServer) GetOrderByID(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	order, err := s.orderService.GetOrderByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	protoOrder := mappers.ToProtoOrder(order)

	return &orderpb.GetOrderResponse{
		Order: protoOrder,
	}, nil
}

func (s *OrderGrpcServer) UpdateOrder(ctx context.Context, req *orderpb.UpdateOrderRequest) (*orderpb.UpdateOrderResponse, error) {
	err := s.orderService.UpdateOrder(ctx, req.GetId(), req.GetStatus())

	if err != nil {
		return nil, err
	}

	return &orderpb.UpdateOrderResponse{
		Message: "Order status updated successfully",
	}, nil
}

func (s *OrderGrpcServer) GetOrdersByUserID(ctx context.Context, req *orderpb.GetOrdersByUserRequest) (*orderpb.GetOrdersByUserResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	orders, err := s.orderService.GetOrdersByUserID(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &orderpb.GetOrdersByUserResponse{
		Orders: mappers.ToProtoOrders(orders),
	}, nil
}
