package services

import (
	"context"
	"producer-service/internal/core/models"
	"producer-service/internal/infrastructure/cache"
	logger "producer-service/internal/interfaces/logger"
	"producer-service/internal/interfaces/repositories"
	"producer-service/internal/usecases/services"
	"producer-service/internal/usecases/validators"
	"proto/generated/ecommerce/inventory"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryServiceServer struct {
	inventorypb.UnimplementedInventoryServiceServer
	productService *services.ProductService
	cache cache.CacheService
	logger logger.Logger
}

func NewInventoryGrpcServer(productRepo repositories.ProductRepository, logger logger.Logger, cache cache.CacheService) *InventoryServiceServer {
	return &InventoryServiceServer{
		productService: services.NewProductService(productRepo, logger, cache),
		logger: logger,
	}
}

func (s *InventoryServiceServer) CreateProduct(ctx context.Context, req *inventorypb.CreateProductRequest) (*inventorypb.ProductResponse, error) {
	product := req.GetProduct()
	err := validators.ValidateProductForCreation(models.Product{
		Name:        product.GetName(),
		Description: product.GetDescription(),
		Price:       float64(product.GetPrice()),
		Stock:       int(product.GetStock()),
		CategoryID:  product.GetCategoryId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	createdProduct, err := s.productService.CreateProduct(ctx, models.Product{
		ID:          product.GetId(),
		Name:        product.GetName(),
		Description: product.GetDescription(),
		Price:       float64(product.GetPrice()),
		Stock:       int(product.GetStock()),
		CategoryID:  product.GetCategoryId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &inventorypb.ProductResponse{
		Product: &inventorypb.Product{
			ID:          createdProduct.ID,
			Name:        createdProduct.Name,
			Description: createdProduct.Description,
			Price:       float32(createdProduct.Price),
			Stock:       int32(createdProduct.Stock),
			CategoryID:  createdProduct.CategoryID,
			CreatedAt:   createdProduct.CreatedAt.String(),
			UpdatedAt:   createdProduct.UpdatedAt.String(),
		},
		Message: "Product created successfully",
	}, nil
}

func (s *InventoryServiceServer) GetProductByID(ctx context.Context, req *inventorypb.GetProductByIDRequest) (*inventorypb.ProductResponse, error) {
	product, err := s.productService.GetProductByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &inventorypb.ProductResponse{
		Product: &inventorypb.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       float32(product.Price),
			Stock:       int32(product.Stock),
			CategoryID:  product.CategoryID,
			CreatedAt:   product.CreatedAt.String(),
			UpdatedAt:   product.UpdatedAt.String(),
		},
		Message: "Product retrieved successfully",
	}, nil
}

func (s *InventoryServiceServer) UpdateProduct(ctx context.Context, req *inventorypb.UpdateProductRequest) (*inventorypb.ProductResponse, error) {
	product := req.GetProduct()
	updatedProduct, err := s.productService.UpdateProduct(ctx, req.GetId(), models.Product{
		ID:          product.GetId(),
		Name:        product.GetName(),
		Description: product.GetDescription(),
		Price:       float64(product.GetPrice()),
		Stock:       int(product.GetStock()),
		CategoryID:  product.GetCategoryId(),
	})
	if err != nil {
		return nil, err
	}

	return &inventorypb.ProductResponse{
		Product: &inventorypb.Product{
			ID:          updatedProduct.ID,
			Name:        updatedProduct.Name,
			Description: updatedProduct.Description,
			Price:       float32(updatedProduct.Price),
			Stock:       int32(updatedProduct.Stock),
			CategoryID:  updatedProduct.CategoryID,
			UpdatedAt:   updatedProduct.UpdatedAt.String(),
		},
		Message: "Product updated successfully",
	}, nil
}

func (s *InventoryServiceServer) DeleteProduct(ctx context.Context, req *inventorypb.DeleteProductRequest) (*inventorypb.DeleteProductResponse, error) {
	err := s.productService.DeleteProduct(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &inventorypb.DeleteProductResponse{
		Success: true,
	}, nil
}

func (s *InventoryServiceServer) ListProducts(ctx context.Context, req *inventorypb.ListProductsRequest) (*inventorypb.ListProductsResponse, error) {
	filter := make(map[string]interface{})
	for key, value := range req.GetFilter() {
		filter[key] = value
	}

	products, err := s.productService.ListProducts(ctx, filter, req.GetSkip(), req.GetLimit())
	if err != nil {
		return nil, err
	}

	var grpcProducts []*inventorypb.Product
	for _, product := range products {
		grpcProducts = append(grpcProducts, &inventorypb.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       float32(product.Price),
			Stock:       int32(product.Stock),
			CategoryID:  product.CategoryID,
			CreatedAt:   product.CreatedAt.String(),
			UpdatedAt:   product.UpdatedAt.String(),
		})
	}

	return &inventorypb.ListProductsResponse{
		Products: grpcProducts,
	}, nil
}

func (s *InventoryServiceServer) CheckStock(ctx context.Context, req *inventorypb.CheckStockRequest) (*inventorypb.CheckStockResponse, error) {
	inStock, availableStock, err := s.productService.CheckStock(ctx, req.GetProductId(), req.GetQuantity())
	if err != nil {
		return nil, err
	}

	return &inventorypb.CheckStockResponse{
		InStock:        inStock,
		AvailableStock: availableStock,
	}, nil
}

func (s *InventoryServiceServer) DecreaseStock(ctx context.Context, req *inventorypb.DecreaseStockRequest) (*inventorypb.DecreaseStockResponse, error) {
	updatedProduct, err := s.productService.DecreaseStock(ctx, req.GetProductId(), req.GetQuantity())
	if err != nil {
		return &inventorypb.DecreaseStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &inventorypb.DecreaseStockResponse{
		Success: true,
		Message: "stock decreased successfully",
		Product: &inventorypb.Product{
			ID:          updatedProduct.ID,
			Name:        updatedProduct.Name,
			Description: updatedProduct.Description,
			Price:       float32(updatedProduct.Price),
			Stock:       int32(updatedProduct.Stock),
			CategoryID:  updatedProduct.CategoryID,
			CreatedAt:   updatedProduct.CreatedAt.String(),
			UpdatedAt:   updatedProduct.UpdatedAt.String(),
		},
	}, nil
}

