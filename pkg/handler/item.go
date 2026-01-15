package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"restapi-app"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createItem(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		return
	}

	listid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input restapi.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoItem.CreateItem(userid, listid, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllItems(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		return
	}

	listid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.services.TodoItem.GetAllItems(userid, listid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}

func (h *Handler) getItemByID(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		return
	}

	// listid, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	itemid, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.services.TodoItem.GetItemByID(userid, itemid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newErrorResponse(c, http.StatusBadRequest, errors.New("Item not exists").Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"item": item,
	})
}

func (h *Handler) updateItem(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		return
	}

	// listid, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	itemid, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input restapi.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.TodoItem.UpdateItem(userid, itemid, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "item updated",
	})
}

func (h *Handler) deleteItem(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		return
	}

	// listid, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	itemid, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.TodoItem.DeleteItem(userid, itemid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "item deleted",
	})
}
