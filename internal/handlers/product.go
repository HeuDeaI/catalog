package handlers

import (
	"catalog/internal/models"
	"catalog/internal/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/products", h.CreateProduct)
	router.GET("/products/:id", h.GetProductByID)
	router.GET("/products", h.GetProducts)
	router.PUT("/products", h.UpdateProduct)
	router.DELETE("/products/:id", h.DeleteProduct)
}

// CreateProduct godoc
// @Summary     Create a new product
// @Description Uploads product info and image via multipart form
// @Tags        Products
// @Accept      multipart/form-data
// @Produce     json
// @Param       name           formData  string   true   "Product name"
// @Param       price          formData  number   true   "Price"
// @Param       brand_id       formData  int      false  "Brand ID"
// @Param       skin_type_ids  formData  []int    false  "Skin type IDs"
// @Param       image          formData  file     true   "Image file"
// @Success     201            {object}  models.Product
// @Failure     400            {object}  map[string]string
// @Failure     500            {object}  map[string]string
// @Router      /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.Bind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file required"})
		return
	}

	filePath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
		return
	}

	if err := h.service.CreateProduct(&product, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

// GetProductByID godoc
// @Summary     Return product by ID
// @Tags        Products
// @Produce     json
// @Param       id   path      int  true  "Product ID"
// @Success     200  {object}  models.Product
// @Failure     400            {object}  map[string]string
// @Failure     500            {object}  map[string]string
// @Router      /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetProducts godoc
// @Summary     Return products
// @Description Returns all products or filters them by parameters (price, skin type)
// @Tags        Products
// @Accept      json
// @Produce     json
// @Param       min_price   query    number  false  "The minimum price of the product"
// @Param       max_price   query    number  false  "The maximum price of the product"
// @Param       skin_type   query    string  false  "Skin type ID (comma-separated)"
// @Success     200         {array}  models.Product
// @Failure     400         {object} map[string]string
// @Router      /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	if c.Query("min_price") != "" || c.Query("max_price") != "" || c.Query("skin_type") != "" {
		h.GetProductsByFilter(c)
		return
	}
	h.GetAllProducts(c)
}

func (h *ProductHandler) GetProductsByFilter(c *gin.Context) {
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")
	skinTypeStr := c.Query("skin_type")
	var minPrice, maxPrice float64
	var err error

	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_price"})
			return
		}
	}

	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max_price"})
			return
		}
	}

	if minPrice < 0 || maxPrice < 0 || minPrice > maxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price range"})
		return
	}

	var skinTypeIDs []uint
	if skinTypeStr != "" {
		parts := strings.Split(skinTypeStr, ",")
		for _, part := range parts {
			id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skin_type"})
				return
			}
			skinTypeIDs = append(skinTypeIDs, uint(id))
		}
	}

	products, err := h.service.GetProductsByFilter(minPrice, maxPrice, skinTypeIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// UpdateProduct godoc
// @Summary     Update product
// @Description Update product info with optional new image
// @Tags        Products
// @Accept      multipart/form-data
// @Produce     json
// @Param       id             formData  int      true   "Product ID"
// @Param       name           formData  string   false  "Name"
// @Param       price          formData  number   false  "Price"
// @Param       brand_id       formData  int      false  "Brand ID"
// @Param       skin_type_ids  formData  []int    false  "Skin type IDs"
// @Param       image          formData  file     false  "New image file"
// @Success     200            {object}  models.Product
// @Failure     400  {object}  map[string]string
// @Failure     404  {object}  map[string]string
// @Router      /products [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var product models.Product
	if err := c.Bind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("image")
	var filePath string
	if err == nil {
		filePath = "/tmp/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
			return
		}
	}

	if err := h.service.UpdateProduct(&product, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product updated", "product": product})
}

// DeleteProduct godoc
// @Summary     Delete product
// @Tags        Products
// @Param       id   path  int  true  "Product ID"
// @Success     200  {object}  map[string]string
// @Failure     400  {object}  map[string]string
// @Failure     404  {object}  map[string]string
// @Router      /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.service.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
