package handler

import (
	"avito/internal/domain"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) RegisterHandler(c *gin.Context) {
	var req domain.AuthStruct

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.AuthService.RegisterService(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered"})
}

func (h Handler) LoginHandler(c *gin.Context) {
	var req domain.AuthStruct

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.AuthService.LoginService(req)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect login or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot login user"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", token, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "user logged in"})

}

func (h Handler) RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
	}

	user, err := h.AuthService.CheckTokenService(tokenString)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
	}

	c.Set("user", user)

	c.Next()
}

func (h Handler) AdminCheck(c *gin.Context) {
	user, _ := c.Get("user")

	if !user.(domain.AuthStruct).IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden, you are not an admin"})
		c.Abort()
	}
	c.Next()
}
