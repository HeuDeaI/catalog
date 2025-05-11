package routes

import (
	"catalog/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, productHandler *handlers.ProductHandler) {
	api := router.Group("/products")
	{
		api.POST("/", productHandler.CreateProduct)
		api.GET("/", productHandler.GetAllProducts)
		api.GET("/:id", productHandler.GetProductByID)
		api.PUT("/:id", productHandler.UpdateProduct)
		api.DELETE("/:id", productHandler.DeleteProduct)
	}
}
