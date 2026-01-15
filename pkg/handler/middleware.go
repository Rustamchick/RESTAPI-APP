package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

// userIdentity write userId in context for using in handlers by context
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid rest auth header")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	isAdmin, err := h.grpc_auth.IsAdmin(context.Background(), int64(userId))
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "invalid grpc auth header")
		return
	}

	if !isAdmin {
		logrus.Infof("grpc auth is OFF")
		newErrorResponse(c, http.StatusUnauthorized, "You are not admin")
		return
	}
	c.Set(userCtx, userId)

	c.Next()
}

func getUserId(c *gin.Context) (int, error) {
	id, exist := c.Get(userCtx)
	if !exist {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	userid, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found. type assertion")
		return 0, errors.New("user id not found")
	}

	return userid, nil
}
