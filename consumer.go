package consumer

import (
	"code.google.com/p/goconf/conf"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"syscall"
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
func bind(config *conf.ConfigFile, conn *amqp.Connection) (q queue, err error) {
	channel, err := conn.Channel()
	if err != nil {
		return
	}
	ex, q, err := readConfigFile(config)
	log.Printf("Declaring Exchange %s", ex)
	err = channel.ExchangeDeclare(ex.name, ex.kind, ex.durable, ex.autoDelete, false, false, nil)
	if err != nil {
		return
	}
	log.Printf("Declaring Queue %s", q)
	_, err = channel.QueueDeclare(q.name, q.durable, q.autoDelete, q.exclusive, false, nil)
	if err != nil {
		return
	}
	log.Printf("Declaring Binding %s routingkey=%s", q.name, q.routingKey)
	err = channel.QueueBind(q.name, q.routingKey, ex.name, false, nil)
	if err != nil {
		return
	}
	return
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

type worker func(*Message)

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
	channel   *amqp.Channel
	queue     queue
	connected bool
}

func (c *Consumer) Queue() queue {
	return c.queue
}

/*
Connect to the AMQP server.

Will do the following work:

- Create the connection.
- Declare the exchange.
- Declare the queue.
- Bind the queue + exchange together.
*/
func (c *Consumer) Connect() (err error) {
	if c.connected {
		return err
	}

	conn, err := connect(c.conf)
	if err != nil {
		return
	}

	q, err := bind(c.conf, conn)
	if err != nil {
		return
	}
	c.conn = conn
	c.queue = q
	return
}

/*
Takes a function that accepts amqp.Delivery and binds
it to the configured queue.

The provided function will be called each time a message is
received and the function is expected to Ack or Nack the message.
*/
func (c *Consumer) Consume(handler worker) (err error) {
	err = c.Connect()
	if err != nil {
		return
	}
	channel, err := c.conn.Channel()
	queue := c.Queue()

	log.Printf("Consuming from queue: %s", queue.Name())
	messages, err := channel.Consume(queue.Name(), queue.Tag(), false, queue.Exclusive(), false, false, nil)
	if err != nil {
		return
	}

	go c.process(handler, messages)
	c.StartLoop()
	return
}

/*
Consumer from the channel - run inside a separate goroutine
*/
func (c *Consumer) process(handler worker, messages <-chan amqp.Delivery) {
	for rawMsg := range messages {
		msg := &Message{rawMsg}
		handler(msg)
	}
}

/*
Start the loop that keeps the process alive.

Registers signal handlers to cancel consumers, on
signals.
*/
func (c *Consumer) StartLoop() {
	kill := make(chan os.Signal, 1)

	// Listen for common kill types
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	select {
	case s := <-kill:
		log.Printf("Caught signal %s Stopping consumer.", s)
		channel, _ := c.conn.Channel()
		err := channel.Cancel(c.queue.Tag(), false)
		if err != nil {
			log.Fatalf("Could not close channel.")
		}
		c.conn.Close()
		log.Print("Channel closed.")
	}
}


/*
Simple message type so users of this library don't have to import amqp as well
*/
type Message struct {
	amqp.Delivery
}
