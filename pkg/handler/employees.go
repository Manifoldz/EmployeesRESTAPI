package handler

import (
	"net/http"

	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createEmployee(c *gin.Context) {
	var input entities.CreateEmployeeInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Employees.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllEmployees(c *gin.Context) {}

func (h *Handler) getEmployeeById(c *gin.Context) {}

func (h *Handler) updateEmployeeById(c *gin.Context) {}

func (h *Handler) deleteEmployeeById(c *gin.Context) {}
