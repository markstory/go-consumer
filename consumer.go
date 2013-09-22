package consumer

import (
	"github.com/streadway/amqp"
	"code.google.com/p/goconf/conf"
	"log"
	"fmt"
)

// Create the amqp:// url from the config file.
func makeAmqpUrl(config *conf.ConfigFile) string {
	options := map[string]string{
		"host": "localhost",
		"vhost": "/",
		"user": "guest",
		"password": "guest",
		"port": "5672",
	}
	for key, _ := range options {
		if config.HasOption("connection", key) {
			options[key], _ = config.GetString("connection", key)
		}
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%s%s",
		options["user"],
		options["password"],
		options["host"],
		options["port"],
		options["vhost"])
}


// Create a new consumer using the connection, exchange
// binding and queue configurations in the provide configuration
// file. Once created you can bind consumers to start handling messages
func Create(configFile string) (c *Consumer, err error) {
	log.Printf("Creating new consumer for config file: %s", configFile)

	config, err := conf.ReadConfigFile(configFile)
	if err != nil {
		return
	}

	amqpUrl := makeAmqpUrl(config)
	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		return
	}

	c = &Consumer{
		conn: conn,
		conf: config,
	}

	return
}

type worker func(amqp.Delivery)

// The main type that users of this package interact with
//
type Consumer struct {
	conn *amqp.Connection
	conf *conf.ConfigFile
}

// Takes a function that accepts amqp.Delivery and binds
// it to the configured queue.
//
// The provided function will be called each time a message is
// received and the function is expected to Ack or Nack the message.
//
func (c *Consumer) Consume(worker) (ok bool, err error) {
	return
}
