package handler

import (
	"net/http"
	"strconv"

	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createEmployee(c *gin.Context) {
	var input entities.EmployeeInputAndResponse

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

func (h *Handler) getAllEmployees(c *gin.Context) {

	// извлечение company_id из запроса
	companyQuery := c.Query("company_id")
	var companyId *int
	if companyQuery != "" {
		id, err := strconv.Atoi(companyQuery)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "invalid company_id parameter")
			return
		}
		companyId = &id
	}

	// извлечение department_id из запроса
	departmentQuery := c.Query("department_id")
	var departmentId *int
	if departmentQuery != "" {
		id, err := strconv.Atoi(departmentQuery)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "invalid department_id parameter")
			return
		}
		departmentId = &id
	}

	// извлечение параметров пагинации из запроса
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	employees, err := h.services.Employees.GetAll(companyId, departmentId, offset, limit)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, employees)
}

func (h *Handler) getEmployeeById(c *gin.Context) {}

func (h *Handler) updateEmployeeById(c *gin.Context) {}

func (h *Handler) deleteEmployeeById(c *gin.Context) {}
