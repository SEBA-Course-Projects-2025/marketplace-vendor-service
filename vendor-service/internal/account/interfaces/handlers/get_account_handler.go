package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/account/application/services"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
	"net/http"
)

// GetAccountHandler godoc
// @Summary      Get account by vendor Id
// @Description  Returns the account for the given vendor.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer access token"
// @Success      200 {object} dtos.AccountResponse
// @Failure      400 {object} map[string]interface{} "Invalid vendorId"
// @Failure      404 {object} map[string]interface{} "Account not found"
// @Failure      500 {object} map[string]interface{}
// @Router       /account [get]
func (h *AccountHandler) GetAccountHandler(c *gin.Context) {

	ctx, span := tracer.Tracer.Start(c.Request.Context(), "GetAccountHandler")
	defer span.End()

	v, _ := c.Get("vendorId")
	vendorId, ok := v.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendorId"})
		return
	}

	account, err := services.GetAccount(ctx, h.AccountRepo, vendorId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)

}
