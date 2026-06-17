package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Putthakun/face-recognition-api-go/internal/usecase"
	"github.com/gin-gonic/gin"
)

type CameraHandler struct {
	uc usecase.CameraUsecase
}

func NewCameraHandler(uc usecase.CameraUsecase) *CameraHandler {
	return &CameraHandler{uc: uc}
}

// GET /api/admin/cameras
func (h *CameraHandler) GetAll(c *gin.Context) {
	cameras, err := h.uc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cameras)
}

// POST /api/admin/cameras
func (h *CameraHandler) Create(c *gin.Context) {
	var body struct {
		Location string `json:"location" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	cam, err := h.uc.Create(body.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cam)
}

// PUT /api/admin/cameras/:cameraId
func (h *CameraHandler) Update(c *gin.Context) {
	cameraID, err := strconv.ParseInt(c.Param("cameraId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid cameraId"})
		return
	}

	var body struct {
		Location string `json:"location" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	cam, err := h.uc.Update(cameraID, body.Location)
	if err != nil {
		if errors.Is(err, usecase.ErrCameraNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cam)
}

// DELETE /api/admin/cameras/:cameraId
func (h *CameraHandler) Delete(c *gin.Context) {
	cameraID, err := strconv.ParseInt(c.Param("cameraId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid cameraId"})
		return
	}

	if err := h.uc.Delete(cameraID); err != nil {
		if errors.Is(err, usecase.ErrCameraNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
