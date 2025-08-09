package interfaces

import (
	"github.com/gin-gonic/gin"
	"marketplace-vendor-service/vendor-service/internal/account/interfaces/handlers"
)

func SetUpAccountsRouter(rg *gin.RouterGroup, h *handlers.AccountHandler) {
	accounts := rg.Group("account")
	{
		accounts.GET("/", h.GetAccountHandler)
		accounts.PUT("/", h.PutAccountHandler)
		accounts.PATCH("/", h.PatchAccountHandler)
	}

}
