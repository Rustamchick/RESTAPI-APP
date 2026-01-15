package handler

import (
	grpc_auth "restapi-app/pkg/client/auth/grpc"
	"restapi-app/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services  *service.Service
	grpc_auth *grpc_auth.Client
}

func NewHandler(services *service.Service, authClient *grpc_auth.Client) *Handler {
	return &Handler{
		services:  services,
		grpc_auth: authClient,
	}
}

// InitRoutes() инициализирует все наши эндпоинты
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListByID)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}
		items := api.Group("/items")
		{
			items.GET("/:item_id", h.getItemByID)
			items.PUT("/:item_id", h.updateItem)
			items.DELETE("/:item_id", h.deleteItem)
		}
	}

	return router
}
