package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/account/application/services"
	"marketplace-vendor-service/vendor-service/internal/account/dtos"
	"net/http"
)

// LoginAccountHandler godoc
// @Summary      Login account
// @Description  Authenticates an account and returns a JWT token.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        loginReq body dtos.LoginRequestDto true "Login credentials"
// @Success      200 {object} map[string]string "JWT token"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      401 {object} map[string]interface{} "Invalid account credentials"
// @Failure      500 {object} map[string]interface{}
// @Router       /auth/login [post]
func (h *AccountHandler) LoginAccountHandler(c *gin.Context) {

	var loginReq dtos.LoginRequestDto

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := services.LoginAccount(c.Request.Context(), h.AccountRepo, loginReq)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid account credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
