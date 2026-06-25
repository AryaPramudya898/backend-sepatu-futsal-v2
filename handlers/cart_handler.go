package handlers

import (
	"net/http"
	"strconv"

	"github.com/AryaPramudya898/backend-sepatu-futsal.git/config"
	"github.com/AryaPramudya898/backend-sepatu-futsal.git/models"
	"github.com/gin-gonic/gin"
)

type CartHandler struct{}

func NewCartHandler() *CartHandler {
	return &CartHandler{}
}

func (h *CartHandler) getUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User ID tidak ditemukan di token",
		})
		return 0, false
	}
	if f, ok := val.(float64); ok {
		return uint(f), true
	}
	if u, ok := val.(uint); ok {
		return u, true
	}
	if i, ok := val.(int); ok {
		return uint(i), true
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": "User ID memiliki tipe yang tidak valid",
	})
	return 0, false
}

// GET /v1/cart
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var items []models.CartItem
	if err := config.DB.Preload("Product").Where("user_id = ?", userID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data keranjang: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
	})
}

// POST /v1/cart
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var req models.AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid: product_id dan quantity (min 1) diperlukan",
		})
		return
	}

	// Cek apakah produk ada
	var product models.Product
	if err := config.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Produk tidak ditemukan",
		})
		return
	}

	// Cari apakah item sudah ada di keranjang user
	var item models.CartItem
	err := config.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&item).Error

	if err == nil {
		// Jika ada, tambahkan quantity
		item.Quantity += req.Quantity
		if err := config.DB.Save(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal mengupdate item keranjang: " + err.Error(),
			})
			return
		}
	} else {
		// Jika tidak ada, buat baru
		item = models.CartItem{
			UserID:    userID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := config.DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal menambahkan ke keranjang: " + err.Error(),
			})
			return
		}
	}

	// Preload Product detail untuk response
	config.DB.Preload("Product").First(&item, item.ID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item berhasil ditambahkan ke keranjang",
		"data":    item,
	})
}

// PUT /v1/cart
func (h *CartHandler) UpdateQuantity(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid: product_id dan quantity diperlukan",
		})
		return
	}

	var item models.CartItem
	err := config.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&item).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Item tidak ditemukan di keranjang",
		})
		return
	}

	if req.Quantity <= 0 {
		// Hapus item jika quantity <= 0
		if err := config.DB.Delete(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal menghapus item dari keranjang: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Item berhasil dihapus dari keranjang karena quantity <= 0",
		})
		return
	}

	item.Quantity = req.Quantity
	if err := config.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengupdate quantity: " + err.Error(),
		})
		return
	}

	config.DB.Preload("Product").First(&item, item.ID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Quantity berhasil diupdate",
		"data":    item,
	})
}

// DELETE /v1/cart/:product_id
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Product ID tidak valid",
		})
		return
	}

	var item models.CartItem
	err = config.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Item tidak ditemukan di keranjang",
		})
		return
	}

	if err := config.DB.Unscoped().Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus item: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item berhasil dihapus dari keranjang",
	})
}

// DELETE /v1/cart
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	if err := config.DB.Unscoped().Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengosongkan keranjang: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Keranjang berhasil dikosongkan",
	})
}
