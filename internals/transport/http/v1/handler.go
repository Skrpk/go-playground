package v1

import (
	"github.com/gin-gonic/gin"
	services "gitlab.com/saas-template1/go-monolith/internals/services"
	m "gitlab.com/saas-template1/go-monolith/internals/transport/http/middlewares"
)

type Handler struct {
	paymentHandler *PaymentsHandler
	authHandler    *AuthHandler
}

func NewHandler(pS *services.PaymentsService, aS *services.AuthService) *Handler {
	paymentsHandler := NewPaymentsHandler(pS)
	authHandler := NewAuthHandler(aS)
	return &Handler{
		paymentHandler: paymentsHandler,
		authHandler:    authHandler,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	h.initPaymentsRoutes(v1)
	h.initAuthRoutes(v1)
}

func (h *Handler) initPaymentsRoutes(api *gin.RouterGroup) {
	newApi := api.Group("/payments")
	newApi.POST("/checkout", h.paymentHandler.Checkout)
	newApi.POST("/webhook", h.paymentHandler.HandleWebhook)
	newApi.GET("/success", h.paymentHandler.Success)
}

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	jwtMiddleware := m.MakeJwtVerificationMiddleware()
	newApi := api.Group("/auth")
	newApi.POST("/signin", h.authHandler.SignIn)
	newApi.POST("/signup", h.authHandler.SignUp)
	newApi.POST("/refresh", h.authHandler.Refresh)
	newApi.GET("/test", jwtMiddleware, h.authHandler.Test)
}
