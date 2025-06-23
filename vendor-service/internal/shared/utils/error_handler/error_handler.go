package error_handler

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func ErrorHandler(err error, message string) error {

	logrus.WithError(err).Error(message)
	return fmt.Errorf("%s: %w", message, err)

}
