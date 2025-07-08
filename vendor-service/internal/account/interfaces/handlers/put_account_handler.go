package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/account/application/services"
	"marketplace-vendor-service/vendor-service/internal/account/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
	"net/http"
)

// PutAccountHandler godoc
// @Summary      Fully update account
// @Description  Fully updates the account for the given vendor.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer access token"
// @Param        accountReq body dtos.PutRequestDto true "Account updated data"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{} "Invalid vendorId or account data"
// @Failure      404 {object} map[string]interface{} "Account not found"
// @Failure      500 {object} map[string]interface{}
// @Router       /account [put]
func (h *AccountHandler) PutAccountHandler(c *gin.Context) {

	ctx, span := tracer.Tracer.Start(c.Request.Context(), "PutAccountHandler")
	defer span.End()

	v, _ := c.Get("vendorId")
	vendorId, ok := v.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendorId"})
		return
	}

	var accountReq dtos.PutRequestDto

	if err := c.ShouldBindJSON(&accountReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account data"})
		return
	}

	if err := services.PutAccount(ctx, h.AccountRepo, accountReq, vendorId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}
