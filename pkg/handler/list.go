package handler

import (
	"net/http"
	"strconv"

	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/gin-gonic/gin"
)

// @Summary Create todo list
// @Tags lists
// @Description create todo list
// @ID create-list
// @Accept  json
// @Produce  json
// @Param input body entities.ToDoList true "list info"
// @Success 200 {object} map[string]int "id of the created list"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [post]

func (h *Handler) createList(c *gin.Context) {
	var input entities.ToDoList

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.ToDoList.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllListResponse struct {
	Data []entities.ToDoList `json:"data"`
}

// @Summary Get all todo lists
// @Tags lists
// @Description get all todo lists
// @ID get-all-lists
// @Produce json
// @Success 200 {object} getAllListResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [get]

func (h *Handler) getAllList(c *gin.Context) {
	lists, err := h.services.ToDoList.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllListResponse{
		Data: lists,
	})
}

// @Summary Get todo list
// @Tags lists
// @Description get  todo list by id
// @ID get-list
// @Produce json
// @Param id path int true "list id"
// @Success 200 {object} entities.ToDoList
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/:id [get]

func (h *Handler) getListById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}

	list, err := h.services.ToDoList.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

// @Summary Update todo list
// @Tags lists
// @Description update todo list by id
// @ID update-list
// @Produce json
// @Param id path int true "list id"
// @Param input body entities.UpdateListInput true "list info"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/:id [put]

func (h *Handler) updateListById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}

	var input entities.UpdateListInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.ToDoList.UpdateById(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "success"})
}

// @Summary Delete todo list
// @Tags lists
// @Description delete todo list by id
// @ID delete-list
// @Produce json
// @Param id path int true "list id"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/:id [delete]

func (h *Handler) deleteListById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}

	err = h.services.ToDoList.DeleteById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "success",
	})

}
