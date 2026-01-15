package handler

import (
	"net/http"
	"restapi-app"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		return
	}

	var input restapi.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listid, err := h.services.TodoList.Create(userid, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": listid,
	})

}

type getAllListsResponse struct {
	Data []restapi.TodoList `json:"data"`
	Err  error              `json:"error"`
}

func (h *Handler) getAllLists(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAllLists(userid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllListsResponse{ // здесь как бы не обязательно дополнительно структуру придумывать вроде бы, так как просто листс тоже можно вставить и всё работает
		Data: lists,
		Err:  nil,
	})
}

func (h *Handler) getListByID(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		return
	}

	listid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	list, err := h.services.TodoList.GetListById(userid, listid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"list": list,
	})
}

func (h *Handler) updateList(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		return
	}

	listid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input restapi.UpdateListInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err = h.services.TodoList.UpdateList(userid, listid, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "list updated",
	})
}

func (h *Handler) deleteList(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		return
	}

	listid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.TodoList.Delete(userid, listid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "list deleted",
	})
}
