package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Putthakun/face-recognition-api-go/internal/usecase"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	uc usecase.EmployeeUsecase
}

func NewEmployeeHandler(uc usecase.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{uc: uc}
}

// POST /api/admin/employees  (multipart/form-data)
func (h *EmployeeHandler) Create(c *gin.Context) {
	empIDStr := c.PostForm("empId")
	empID, err := strconv.ParseInt(empIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "empId is required and must be a number"})
		return
	}

	photo, _ := c.FormFile("photo")
	input := usecase.CreateEmployeeInput{
		EmpID:    empID,
		Name:     c.PostForm("name"),
		Password: c.PostForm("password"),
		Role:     c.PostForm("role"),
		Photo:    photo,
	}

	result, err := h.uc.Create(input)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrEmployeeAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		case errors.Is(err, usecase.ErrFaceNotDetected):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "field": "photo"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GET /api/admin/employees
func (h *EmployeeHandler) GetAll(c *gin.Context) {
	employees, err := h.uc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

// PUT /api/admin/employees/:empId  (multipart/form-data)
func (h *EmployeeHandler) Update(c *gin.Context) {
	empID, err := strconv.ParseInt(c.Param("empId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid empId"})
		return
	}

	photo, _ := c.FormFile("photo")
	var isActive *bool
	if v := c.PostForm("isActive"); v != "" {
		b, _ := strconv.ParseBool(v)
		isActive = &b
	}

	input := usecase.UpdateEmployeeInput{
		Name:     c.PostForm("name"),
		Password: c.PostForm("password"),
		Role:     c.PostForm("role"),
		IsActive: isActive,
		Photo:    photo,
	}

	result, err := h.uc.Update(empID, input)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrEmployeeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		case errors.Is(err, usecase.ErrFaceNotDetected):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "field": "photo"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

// DELETE /api/admin/employees/:empId
func (h *EmployeeHandler) Delete(c *gin.Context) {
	empID, err := strconv.ParseInt(c.Param("empId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid empId"})
		return
	}

	if err := h.uc.Delete(empID); err != nil {
		if errors.Is(err, usecase.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
