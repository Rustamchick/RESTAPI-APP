package handler

import (
	"context"
	"net/http"
	"restapi-app"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// signUp
func (h *Handler) signUp(c *gin.Context) {
	var input restapi.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// registraion to grpc auth server
	userid, err := h.grpc_auth.Register(context.Background(), input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		logrus.Infof("No registration to auth grpc server")
		return
	}
	logrus.Infof("grpc userid: %d, rest userid: %d", userid, id)

	c.JSON(http.StatusOK, map[string]any{
		"id": id,
	})
}

type SignInUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// signIn
func (h *Handler) signIn(c *gin.Context) {
	var input SignInUser

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"token": token,
	})
}
