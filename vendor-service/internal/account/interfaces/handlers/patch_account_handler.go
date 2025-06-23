package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/account/application/services"
	"marketplace-vendor-service/vendor-service/internal/account/dtos"
	"net/http"
)

// PatchAccountHandler godoc
// @Summary      Partially update account
// @Description  Partially updates the account for the given vendor.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer access token"
// @Param        accountReq body dtos.AccountPatchRequest true "Account patch data"
// @Success      200 {object} dtos.AccountResponse
// @Failure      400 {object} map[string]interface{} "Invalid vendorId or request body"
// @Failure      404 {object} map[string]interface{} "Account not found"
// @Failure      500 {object} map[string]interface{}
// @Router       /account [patch]
func (h *AccountHandler) PatchAccountHandler(c *gin.Context) {

	v, _ := c.Get("vendorId")
	vendorId, ok := v.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendorId"})
		return
	}

	var accountReq dtos.AccountPatchRequest

	if err := c.ShouldBindJSON(&accountReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := services.PatchAccount(c.Request.Context(), h.AccountRepo, accountReq, vendorId)

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
