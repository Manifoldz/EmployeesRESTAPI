package handler

import (
	"net/http"
	"strconv"

	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/gin-gonic/gin"
)

// @Summary Create todo item
// @Tags items
// @Description create todo item to the list
// @ID create-item
// @Accept  json
// @Produce  json
// @Param id path int true "list id"
// @Param input body entities.ToDoItem true "item info"
// @Success 200 {object} map[string]int "id of the created item"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/:id/items [post]

func (h *Handler) createItem(c *gin.Context) {
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id parameter")
		return
	}

	var input entities.ToDoItem
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.ToDoItem.Create(listId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Get todo list-items
// @Tags items
// @Description get all todo items from the list
// @ID get-all-items
// @Accept json
// @Produce json
// @Param id path int true "list id"
// @Param done query string false "filter by done status"
// @Param limit query int false "limit number of items returned" default(10)
// @Param offset query int false "offset for items returned" default(0)
// @Success 200 {object} []entities.ToDoItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/:id/items [get]

func (h *Handler) getAllItem(c *gin.Context) {
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id parameter")
		return
	}

	// Извлечение параметра фильтрации по статусу из запроса
	statusQuery := c.Query("done")
	var status *bool
	if statusQuery != "" {
		statusVal := statusQuery == "true"
		status = &statusVal
	}
	// Извлечение параметров пагинации из запроса
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	items, err := h.services.ToDoItem.GetAll(listId, status, offset, limit)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

// @Summary Get todo item
// @Tags items
// @Description get todo item by id
// @ID get-item
// @Produce json
// @Param id path int true "item id"
// @Success 200 {object} entities.ToDoItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/:id [get]

func (h *Handler) getItemById(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id parameter")
		return
	}
	item, err := h.services.ToDoItem.GetById(itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary Update todo item
// @Tags items
// @Description update todo item by id
// @ID update-item
// @Produce json
// @Param id path int true "item id"
// @Param input body entities.UpdateItemInput true "item info"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/:id [put]

func (h *Handler) updateItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}

	var input entities.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.ToDoItem.UpdateById(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"success"})
}

// @Summary Delete todo item
// @Tags items
// @Description delete todo item by id
// @ID delete-item
// @Produce json
// @Param id path int true "item id"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/:id [delete]

func (h *Handler) deleteItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id parameter")
		return
	}
	err = h.services.ToDoItem.DeleteById(itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"success"})
}
