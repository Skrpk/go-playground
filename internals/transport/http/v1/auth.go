package v1

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/saas-template1/go-monolith/internals/services"
	"net/http"
)

type AuthHandler struct {
	authService *services.AuthService
}

type SignInPayload struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type RefreshPayload struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

type SignUpPayload struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (aH *AuthHandler) SignIn(c *gin.Context) {
	var payload SignInPayload
	if c.BindJSON(&payload) == nil {
		err, res := aH.authService.SignIn(payload.Email, payload.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err.Error()})
			return
		}
		c.SetCookie("refresh_token", res.RefreshToken, 300000, "/", "localhost", true, true)
		c.JSON(http.StatusOK, gin.H{"token": res.AccessToken})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": "validation failed"})
}
func (aH *AuthHandler) Refresh(c *gin.Context) {
	var payload RefreshPayload
	if c.BindJSON(&payload) == nil {
		err, res := aH.authService.Refresh(payload.RefreshToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": "validation failed"})
}

func (aH *AuthHandler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, "hello")
}

func (aH *AuthHandler) SignUp(c *gin.Context) {
	var payload SignUpPayload
	if c.BindJSON(&payload) == nil {
		err, _ := aH.authService.SignUp(payload.Email, payload.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": "validation failed"})
}
