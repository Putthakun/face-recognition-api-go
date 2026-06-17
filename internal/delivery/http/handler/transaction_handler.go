package handler

import (
	"net/http"
	"strconv"

	"github.com/Putthakun/face-recognition-api-go/internal/usecase"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	uc usecase.TransactionUsecase
}

func NewTransactionHandler(uc usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{uc: uc}
}

// POST /api/transactions  (internal — no auth)
func (h *TransactionHandler) Create(c *gin.Context) {
	var body struct {
		EmpID    *int64 `json:"empId"`
		CameraID *int64 `json:"cameraId"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tx, err := h.uc.Create(body.EmpID, body.CameraID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"transactionId": tx.TransactionID,
		"createdAt":     tx.CreatedAt,
	})
}

// GET /api/transactions  (Admin, Supervisor)
func (h *TransactionHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	sortDesc := c.DefaultQuery("sortDesc", "true") != "false"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	items, total, err := h.uc.GetPaginated(page, pageSize, sortDesc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"items":    items,
	})
}
