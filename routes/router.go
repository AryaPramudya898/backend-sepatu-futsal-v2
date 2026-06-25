package routes

import (
	"github.com/AryaPramudya898/backend-sepatu-futsal.git/handlers"
	"github.com/AryaPramudya898/backend-sepatu-futsal.git/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authHandler := handlers.NewAuthHandler()
	productHandler := handlers.NewProductHandler()
	cartHandler := handlers.NewCartHandler()
	transactionHandler := handlers.NewTransactionHandler()

	v1 := r.Group("/v1")

	v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "service": "sepatu-futsal"})
		})


	auth := v1.Group("/auth")
	{
		auth.POST("/verify-token", authHandler.VerifyToken)
	}

	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		products := protected.Group("/products")
		{
			products.GET("", productHandler.GetAll)
			products.GET("/:id", productHandler.GetByID)

			adminProducts := products.Group("")
			adminProducts.Use(middleware.AdminOnly())
			{
				adminProducts.POST("", productHandler.Create)
				adminProducts.PUT("/:id", productHandler.Update)
				adminProducts.DELETE("/:id", productHandler.Delete)
			}
		}

		cart := protected.Group("/cart")
		{
			cart.GET("", cartHandler.GetCart)
			cart.POST("", cartHandler.AddToCart)
			cart.PUT("", cartHandler.UpdateQuantity)
			cart.DELETE("/:product_id", cartHandler.RemoveFromCart)
			cart.DELETE("", cartHandler.ClearCart)
		}

		transactions := protected.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("", transactionHandler.GetTransactions)
			transactions.PUT("/:reference/status", transactionHandler.UpdateTransactionStatus)
		}
	}
	return r
}
