package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/saas-template1/go-monolith/config"
	services "gitlab.com/saas-template1/go-monolith/internals/services"
	cognitoClient "gitlab.com/saas-template1/go-monolith/internals/services/cognito"
	httpTransport "gitlab.com/saas-template1/go-monolith/internals/transport/http"
)

func Run() {
	conf := config.NewConfig()
	paymentService, err := services.NewPaymentService(*conf.PaymentsConf)
	if err != nil {
		panic(err)
	}

	cognitoClient := cognitoClient.NewCognitoClient(*conf.AwsConfig)
	authService := services.NewAuthService(cognitoClient)

	httpHandler := httpTransport.NewHandler(paymentService, authService)

	server := httpHandler.Init()
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	server.Run()
}
