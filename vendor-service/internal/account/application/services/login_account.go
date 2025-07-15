package services

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"marketplace-vendor-service/vendor-service/internal/account/domain"
	"marketplace-vendor-service/vendor-service/internal/account/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
	"marketplace-vendor-service/vendor-service/internal/shared/utils/error_handler"
	"marketplace-vendor-service/vendor-service/internal/shared/utils/jwt_helper"
)

func LoginAccount(ctx context.Context, accountRepo domain.AccountRepository, loginReq dtos.LoginRequestDto) (string, error) {

	logrus.Info("Starting LoginAccount application service")

	ctx, span := tracer.Tracer.Start(ctx, "LoginAccount")
	defer span.End()

	account, err := accountRepo.FindByEmail(ctx, loginReq.Email)

	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(loginReq.Password)); err != nil {
		return "", error_handler.ErrorHandler(err, err.Error())
	}

	logrus.Info("Successfully logged in account")

	return jwt_helper.GenerateVendorJwt(account.Id)

}
