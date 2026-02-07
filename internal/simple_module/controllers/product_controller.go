package controllers

import (
	"net/http"

	"github.com/refortunato/go_app_base/internal/shared/dto"
	"github.com/refortunato/go_app_base/internal/shared/web/advisor"
	"github.com/refortunato/go_app_base/internal/shared/web/context"
	"github.com/refortunato/go_app_base/internal/simple_module/services"
)

// ProductController handles HTTP requests for products
type ProductController struct {
	service *services.ProductService
}

// NewProductController creates a new product controller instance
func NewProductController(service *services.ProductService) *ProductController {
	return &ProductController{service: service}
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" example:"Laptop Dell XPS 15"`
	Description string  `json:"description" example:"High-performance laptop"`
	Price       float64 `json:"price" example:"5499.99"`
	Stock       int     `json:"stock" example:"10"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name        string  `json:"name" example:"Laptop Dell XPS 15 (Updated)"`
	Description string  `json:"description" example:"Updated description"`
	Price       float64 `json:"price" example:"4999.99"`
	Stock       int     `json:"stock" example:"15"`
}

// GetProduct godoc
// @Summary      Get product by ID
// @Description  Retrieves a specific product from the database
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID (UUID format)"
// @Success      200  {object}  models.Product
// @Failure      404  {object}  errors.ProblemDetails  "Product not found"
// @Failure      500  {object}  errors.ProblemDetails  "Internal server error"
// @Router       /products/{id} [get]
func (c *ProductController) GetProduct(ctx context.WebContext) {
	id := ctx.Param("id")

	product, err := c.service.GetProduct(id)
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// ListProducts godoc
// @Summary      List all products
// @Description  Returns a paginated list of products
// @Tags         products
// @Produce      json
// @Param        page   query  int  false  "Page number" default(1)
// @Param        limit  query  int  false  "Items per page" default(10)
// @Success      200    {object}  services.ListProductsResponse
// @Failure      400    {object}  errors.ProblemDetails   "Invalid pagination parameters"
// @Failure      500    {object}  errors.ProblemDetails   "Internal server error"
// @Router       /products [get]
func (c *ProductController) ListProducts(ctx context.WebContext) {
	// Parse pagination parameters from query string
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")

	pagination, err := dto.NewPaginationRequestDTO(pageStr, limitStr)
	if err != nil {
		advisor.ReturnBadRequestError(ctx, err)
		return
	}

	result, err := c.service.ListProducts(pagination.Page, pagination.Limit)
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// CreateProduct godoc
// @Summary      Create new product
// @Description  Creates a new product in the system
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request  body      CreateProductRequest  true  "Product data"
// @Success      201      {object}  models.Product
// @Failure      400      {object}  errors.ProblemDetails  "Invalid input"
// @Failure      500      {object}  errors.ProblemDetails  "Internal server error"
// @Router       /products [post]
func (c *ProductController) CreateProduct(ctx context.WebContext) {
	var request CreateProductRequest

	if err := ctx.BindJSON(&request); err != nil {
		advisor.ReturnBadRequestError(ctx, err)
		return
	}

	product, err := c.service.CreateProduct(
		request.Name,
		request.Description,
		request.Price,
		request.Stock,
	)
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary      Update product
// @Description  Updates an existing product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      string                 true  "Product ID"
// @Param        request  body      UpdateProductRequest   true  "Updated product data"
// @Success      200      {object}  models.Product
// @Failure      400      {object}  errors.ProblemDetails  "Invalid input"
// @Failure      404      {object}  errors.ProblemDetails  "Product not found"
// @Failure      500      {object}  errors.ProblemDetails  "Internal server error"
// @Router       /products/{id} [put]
func (c *ProductController) UpdateProduct(ctx context.WebContext) {
	id := ctx.Param("id")

	var request UpdateProductRequest

	if err := ctx.BindJSON(&request); err != nil {
		advisor.ReturnBadRequestError(ctx, err)
		return
	}

	product, err := c.service.UpdateProduct(
		id,
		request.Name,
		request.Description,
		request.Price,
		request.Stock,
	)
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary      Delete product
// @Description  Removes a product from the system
// @Tags         products
// @Param        id   path  string  true  "Product ID"
// @Success      204  "No content"
// @Failure      404  {object}  errors.ProblemDetails  "Product not found"
// @Failure      500  {object}  errors.ProblemDetails  "Internal server error"
// @Router       /products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx context.WebContext) {
	id := ctx.Param("id")

	if err := c.service.DeleteProduct(id); err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
