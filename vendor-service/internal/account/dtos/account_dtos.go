package dtos

import (
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/account/domain/account_models"
	"time"
)

type AccountResponse struct {
	Id          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Logo        string    `json:"logo"`
	Address     string    `json:"address"`
	Website     string    `json:"website"`
}

func AccountToDto(account *account_models.VendorAccount) AccountResponse {
	return AccountResponse{
		Id:          account.Id,
		Email:       account.Email,
		Name:        account.Name,
		Description: account.Description,
		Logo:        account.Logo,
		Address:     account.Address,
		Website:     account.Website,
	}
}

type PutRequestDto struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Address     string `json:"address"`
	Website     string `json:"website"`
}

func UpdateAccountWithDto(accountReq PutRequestDto, existingAccount *account_models.VendorAccount) account_models.VendorAccount {
	return account_models.VendorAccount{
		Id:           existingAccount.Id,
		Email:        accountReq.Email,
		PasswordHash: existingAccount.PasswordHash,
		Name:         accountReq.Name,
		Description:  accountReq.Description,
		Logo:         accountReq.Logo,
		Address:      accountReq.Address,
		Website:      accountReq.Website,
		CreatedAt:    existingAccount.CreatedAt,
		UpdatedAt:    time.Now(),
	}
}

type AccountPatchRequest struct {
	Email       *string `json:"email"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Logo        *string `json:"logo"`
	Address     *string `json:"address"`
	Website     *string `json:"website"`
}

func PatchDtoToAccount(existingAccount *account_models.VendorAccount, accountReq AccountPatchRequest) *account_models.VendorAccount {

	if accountReq.Email != nil {
		existingAccount.Email = *accountReq.Email
	}

	if accountReq.Name != nil {
		existingAccount.Name = *accountReq.Name
	}

	if accountReq.Description != nil {
		existingAccount.Description = *accountReq.Description
	}

	if accountReq.Logo != nil {
		existingAccount.Logo = *accountReq.Logo
	}

	if accountReq.Address != nil {
		existingAccount.Address = *accountReq.Address
	}

	if accountReq.Website != nil {
		existingAccount.Website = *accountReq.Website
	}

	existingAccount.UpdatedAt = time.Now()

	return existingAccount
}

type LoginRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
