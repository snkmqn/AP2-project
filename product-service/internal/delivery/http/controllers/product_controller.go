package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"product-service/internal/core/models"
	"product-service/internal/usecases/services"
	"strconv"
	"strings"
)

type ProductController struct {
	productService *services.ProductService
}

func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	log.Printf("Received product: %+v", product)
	ctx := c.Request.Context()
	createdProduct, err := pc.productService.CreateProduct(ctx, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdProduct)
}

func (pc *ProductController) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	product, err := pc.productService.GetProductByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	ctx := c.Request.Context()
	updatedProduct, err := pc.productService.UpdateProduct(ctx, id, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	if err := pc.productService.DeleteProduct(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (pc *ProductController) ListProducts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	categoryID := c.Query("category_id")
	categoryName := c.Query("category_name")

	log.Printf("Received query parameters - page: %s, limit: %s", page, limit)

	limit = strings.TrimSpace(limit)
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	limit = strings.TrimSpace(limit)
	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum <= 0 {
		limitNum = 10
	}

	skip := (pageNum - 1) * limitNum

	log.Printf("Fetching products with skip=%d, limit=%d", skip, limitNum)

	filter := make(map[string]interface{})

	if categoryID != "" {
		filter["category_id"] = categoryID
	}

	if categoryName != "" {
		filter["category_name"] = categoryName
		log.Printf("Filtering by category_name: %s", categoryName)
	}

	ctx := c.Request.Context()
	products, err := pc.productService.ListProducts(ctx, filter, int64(skip), int64(limitNum))
	if err != nil {
		log.Println("ListProducts error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products"})
		return
	}
	log.Printf("Found %d products", len(products))
	if len(products) == 0 {
		log.Println("No products found in the database")
	}

	c.JSON(http.StatusOK, products)
}
