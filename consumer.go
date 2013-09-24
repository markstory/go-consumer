package consumer

import (
	"code.google.com/p/goconf/conf"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// Create the amqp:// url from the config file.
func makeAmqpUrl(config *conf.ConfigFile) string {
	options := map[string]string{
		"host":     "localhost",
		"vhost":    "/",
		"user":     "guest",
		"password": "guest",
		"port":     "5672",
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

/*
Create the amqp.Connection based on the config file.
*/
func connect(config *conf.ConfigFile) (*amqp.Connection, error) {
	amqpUrl := makeAmqpUrl(config)
	return amqp.Dial(amqpUrl)
}

/*
Declare the exchange based on the config file.
*/
func bind(config *conf.ConfigFile, conn *amqp.Connection) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	ex, q, err := readConfigFile(config)
	err = channel.ExchangeDeclare(ex.name, ex.kind, ex.durable, ex.autoDelete, false, false, nil)
	if err != nil {
		return err
	}
	_, err = channel.QueueDeclare(q.name, q.durable, q.autoDelete, q.exclusive, false, nil)
	if err != nil {
		return err
	}
	err = channel.QueueBind(q.name, q.routingKey, ex.name, false, nil)
	if err != nil {
		return err
	}
	return nil
}

/*
Create a new consumer using the connection, exchange
binding and queue configurations in the provide configuration
file. Once created you can bind consumers to start handling messages
*/
func Create(configFile string) (c *Consumer, err error) {
	log.Printf("Creating new consumer for config file: %s", configFile)

	config, err := conf.ReadConfigFile(configFile)
	if err != nil {
		return
	}

	c = &Consumer{
		conf: config,
	}

	return
}

type worker func(amqp.Delivery)

/*
A consumer that applications use to register
functions to act as consumers.

Consumers will connect to the AMQP server when the Consume
method is called. You can manualy connect using the Connect
method as well.

*/
type Consumer struct {
	conf      *conf.ConfigFile
	conn      *amqp.Connection
	connected bool
}

/*
Connect to the AMQP server.

Will do the following work:

- Create the connection
- Declare the exhange
- Declare the queue
- Bind the queue + exchange together.
*/
func (c *Consumer) Connect() (ok bool, err error) {
	if c.connected {
		return true, err
	}

	conn, err := connect(c.conf)
	if err != nil {
		return
	}

	err = bind(c.conf, c.conn)
	if err != nil {
		return
	}
	c.conn = conn
	return
}

/*
Takes a function that accepts amqp.Delivery and binds
it to the configured queue.

The provided function will be called each time a message is
received and the function is expected to Ack or Nack the message.
*/
func (c *Consumer) Consume(worker) (ok bool, err error) {
	return
}
