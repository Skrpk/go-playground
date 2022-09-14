package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	service "gitlab.com/saas-template1/go-monolith/internals/services"
	v1 "gitlab.com/saas-template1/go-monolith/internals/transport/http/v1"
)

type Handler struct {
	paymentService *service.PaymentsService
	authService    *service.AuthService
}

func NewHandler(paymentsService *service.PaymentsService, authService *service.AuthService) *Handler {
	return &Handler{
		paymentService: paymentsService,
		authService:    authService,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.Default(),
	)

	h.initApi(router)

	return router
}

func (h *Handler) initApi(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.paymentService, h.authService)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
