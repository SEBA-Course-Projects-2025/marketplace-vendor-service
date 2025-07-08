package amqp

import (
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"marketplace-vendor-service/vendor-service/internal/shared/utils/error_handler"
	"os"
)

type AMQPConfig struct {
	Channel             *amqp.Channel
	ConfirmationChannel <-chan amqp.Confirmation
}

func SetUpExchange(channel *amqp.Channel) error {

	if err := channel.ExchangeDeclare("vendor.order.events", "direct", true, false, false, false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if _, err := channel.QueueDeclare("order.created.vendor", true, false, false, false, amqp.Table{"x-dead-letter-exchange": "order.vendor.dlx", "x-dead-letter-routing-key": "order.created.vendor.dlq"}); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if err := channel.QueueBind("order.created.vendor", "order.created.vendor", "vendor.order.events", false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if _, err := channel.QueueDeclare("order.canceled.vendor", true, false, false, false, amqp.Table{"x-dead-letter-exchange": "order.vendor.dlx", "x-dead-letter-routing-key": "order.canceled.vendor.dlq"}); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if err := channel.QueueBind("order.canceled.vendor", "order.canceled.vendor", "vendor.order.events", false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if _, err := channel.QueueDeclare("vendor.updated.order", true, false, false, false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if err := channel.QueueBind("vendor.updated.order", "vendor.updated.order", "vendor.order.events", false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if err := channel.ExchangeDeclare("vendor.product.events", "direct", true, false, false, false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if _, err := channel.QueueDeclare("vendor.product.quantity.checked", true, false, false, false, amqp.Table{"x-dead-letter-exchange": "vendor.product.dlx", "x-dead-letter-routing-key": "vendor.product.quantity.checked.dlq"}); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if err := channel.QueueBind("vendor.product.quantity.checked", "vendor.product.quantity.checked", "vendor.product.events", false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if err := channel.QueueBind("vendor.check.product.quantity", "vendor.check.product.quantity", "vendor.product.events", false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}
	
	if err := channel.QueueBind("vendor.cancel.product.order", "vendor.cancel.product.order", "vendor.product.events", false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	return nil

}

func SetUpDlq(channel *amqp.Channel) error {
	if err := channel.ExchangeDeclare("order.vendor.dlx", "direct", true, false, false, false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if _, err := channel.QueueDeclare("order.created.vendor.dlq", true, false, false, false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if err := channel.QueueBind("order.created.vendor.dlq", "order.created.vendor.dlq", "order.vendor.dlx", false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if _, err := channel.QueueDeclare("order.canceled.vendor.dlq", true, false, false, false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if err := channel.QueueBind("order.canceled.vendor.dlq", "order.canceled.vendor.dlq", "order.vendor.dlx", false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if err := channel.ExchangeDeclare("vendor.product.dlx", "direct", true, false, false, false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if _, err := channel.QueueDeclare("vendor.product.quantity.checked.dlq", true, false, false, false, amqp.Table{"x-dead-letter-exchange": "vendor.product.dlx", "x-dead-letter-routing-key": "vendor.product.quantity.checked.dlq"}); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	if err := channel.QueueBind("vendor.product.quantity.checked.dlq", "vendor.product.quantity.checked.dlq", "vendor.product.dlx", false, nil); err != nil {
		return error_handler.ErrorHandler(err, err.Error())
	}

	return nil

}

func ConnectAMQP() (*AMQPConfig, error) {

	if err := godotenv.Load(); err != nil {
		return nil, error_handler.ErrorHandler(err, err.Error())
	}

	amqpUrl := os.Getenv("AMQP_URL")

	connection, err := amqp.Dial(amqpUrl)

	if err != nil {
		return nil, error_handler.ErrorHandler(err, "Error connecting to the CloudAMQP")
	}

	channel, err := connection.Channel()

	if err != nil {
		return nil, error_handler.ErrorHandler(err, "Error opening the channel")
	}

	if err := SetUpExchange(channel); err != nil {
		return nil, error_handler.ErrorHandler(err, err.Error())
	}

	if err := SetUpDlq(channel); err != nil {
		return nil, error_handler.ErrorHandler(err, err.Error())
	}

	if err := channel.Confirm(false); err != nil {
		return nil, error_handler.ErrorHandler(err, err.Error())
	}

	confirmations := channel.NotifyPublish(make(chan amqp.Confirmation, 100))

	return &AMQPConfig{Channel: channel, ConfirmationChannel: confirmations}, nil

}
