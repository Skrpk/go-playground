package v1

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/saas-template1/go-monolith/internals/services"
	"io/ioutil"
	"log"
	"net/http"
)

type PaymentsHandler struct {
	paymentService *services.PaymentsService
}

type CheckoutPayload struct {
	Product      string `json:"product" form:"product" binding:"required"`
	BillingCycle string `json:"billing_cycle" form:"billing_cycle" binding:"required"`
}

func NewPaymentsHandler(paymentService *services.PaymentsService) *PaymentsHandler {
	return &PaymentsHandler{
		paymentService: paymentService,
	}
}

func (pH *PaymentsHandler) Checkout(c *gin.Context) {
	var payload CheckoutPayload
	if c.BindJSON(&payload) == nil {
		url, err := pH.paymentService.Checkout(payload.Product, payload.BillingCycle)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"url": url})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": "validation failed"})
}

func (pH *PaymentsHandler) HandleWebhook(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Printf("ioutil.ReadAll: %v", err)
		return
	}

	err = pH.paymentService.HandleWebhook(b, c.GetHeader("Stripe-Signature"))

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Printf("ioutil.ReadAll: %v", err)
		return
	}
}

func (pH *PaymentsHandler) Success(c *gin.Context) {
	sessionId := c.Query("session_id")
	if sessionId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Session id is missing"})
	}
	status := pH.paymentService.Success(sessionId)

	c.JSON(http.StatusOK, status)
}
