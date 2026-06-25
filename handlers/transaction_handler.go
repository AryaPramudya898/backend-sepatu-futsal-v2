package handlers

import (
	"net/http"

	"github.com/AryaPramudya898/backend-sepatu-futsal.git/config"
	"github.com/AryaPramudya898/backend-sepatu-futsal.git/models"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct{}

func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{}
}

func (h *TransactionHandler) getUserID(c *gin.Context) (uint, bool) {
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

// POST /v1/transactions
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid: reference, amount, dan description diperlukan",
		})
		return
	}

	// Ambil semua item di keranjang user
	var cartItems []models.CartItem
	if err := config.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil item keranjang: " + err.Error(),
		})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Keranjang belanja kosong",
		})
		return
	}

	// Mulai database transaction untuk konsistensi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transaction := models.Transaction{
		UserID:      userID,
		Reference:   req.Reference,
		Amount:      req.Amount,
		Description: req.Description,
		Status:      "pending",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat transaksi: " + err.Error(),
		})
		return
	}

	// Pindahkan cart items ke transaction items
	for _, cartItem := range cartItems {
		txnItem := models.TransactionItem{
			TransactionID: transaction.ID,
			ProductID:     cartItem.ProductID,
			ProductName:   cartItem.Product.Name,
			Price:         cartItem.Product.Price,
			Quantity:      cartItem.Quantity,
		}
		if err := tx.Create(&txnItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal membuat detail transaksi: " + err.Error(),
			})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menyimpan transaksi: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Transaksi berhasil dibuat (pending)",
		"data":    transaction,
	})
}

// GET /v1/transactions
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var transactions []models.Transaction
	if err := config.DB.Preload("Items").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil riwayat transaksi: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transactions,
	})
}

// PUT /v1/transactions/:reference/status
func (h *TransactionHandler) UpdateTransactionStatus(c *gin.Context) {
	reference := c.Param("reference")
	if reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Reference tidak boleh kosong",
		})
		return
	}

	var req models.UpdateTransactionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid: status diperlukan",
		})
		return
	}

	var transaction models.Transaction
	if err := config.DB.Where("reference = ?", reference).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Transaksi tidak ditemukan",
		})
		return
	}

	transaction.Status = req.Status
	if req.TransactionID != "" {
		transaction.TransactionID = req.TransactionID
	}

	if err := config.DB.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengupdate status transaksi: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status transaksi berhasil diupdate ke " + req.Status,
		"data":    transaction,
	})
}
