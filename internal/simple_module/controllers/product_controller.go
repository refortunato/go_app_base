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

// GetProduct handles GET /products/:id
func (c *ProductController) GetProduct(ctx context.WebContext) {
	id := ctx.Param("id")

	product, err := c.service.GetProduct(id)
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// ListProducts handles GET /products
func (c *ProductController) ListProducts(ctx context.WebContext) {
	// Parse pagination parameters from query string
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")

	pagination, err := dto.NewPaginationRequestDTO(pageStr, limitStr)
	if err != nil {
		advisor.ReturnBadRequestError(ctx, err)
		return
	}

	products, err := c.service.ListProducts(pagination.Limit, pagination.Offset)
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": products,
		"pagination": map[string]int{
			"page":   pagination.Page,
			"limit":  pagination.Limit,
			"offset": pagination.Offset,
		},
	})
}

// CreateProduct handles POST /products
func (c *ProductController) CreateProduct(ctx context.WebContext) {
	var request struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}

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

// UpdateProduct handles PUT /products/:id
func (c *ProductController) UpdateProduct(ctx context.WebContext) {
	id := ctx.Param("id")

	var request struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}

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

// DeleteProduct handles DELETE /products/:id
func (c *ProductController) DeleteProduct(ctx context.WebContext) {
	id := ctx.Param("id")

	if err := c.service.DeleteProduct(id); err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
