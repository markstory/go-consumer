package consumer

import (
	"code.google.com/p/goconf/conf"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"syscall"
)

/*
Declare the exchange based on the config file.
*/
func bind(conn *amqp.Connection, top topology) (err error) {
	channel, err := conn.Channel()
	if err != nil {
		return
	}
	for _, bind := range top.Bindings() {
		err := declare(channel, bind)
		if err != nil {
			break
		}
	}
	return
}

func declare(channel *amqp.Channel, bind binding) (err error) {
	ex := bind.Exchange()
	log.Printf("Declaring Exchange %s", ex)
	err = channel.ExchangeDeclare(ex.name, ex.kind, ex.durable, ex.autoDelete, false, false, nil)
	if err != nil {
		return
	}

	q := bind.Queue()
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

	topology, err := NewTopology(config)
	if err != nil {
		return
	}

	c = &Consumer{
		conf: config,
		topology: topology,
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
	topology  topology
	connected bool
}

/*
Get the topology struct for the chosen config file.

The topology struct contains the parsed config file data as simple
structs that can later be traversed to inspect or bind to AMQP.
*/
func (c *Consumer) Topology() topology {
	return c.topology
}

/*
Connect to the AMQP server.

Will do the following work:

Create the connection. Declare the exchange.
Declare the queue. Bind the queue + exchange together.
*/
func (c *Consumer) Connect() (err error) {
	if c.connected {
		return err
	}

	connData := c.topology.Connection()
	conn, err := amqp.Dial(connData.Url())
	if err != nil {
		return
	}

	err = bind(conn, c.topology)
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
func (c *Consumer) Consume(handler worker) (err error) {
	err = c.Connect()
	if err != nil {
		return
	}
	channel, err := c.conn.Channel()
	for _, binding := range c.topology.Bindings() {
		queue := binding.Queue()
		log.Printf("Consuming from queue: %s", queue.Name())

		messages, err := channel.Consume(queue.Name(), queue.Tag(), false, queue.Exclusive(), false, false, nil)
		if err != nil {
			return err
		}
		go c.process(handler, messages)
	}

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
		err := c.Stop()
		if err != nil {
			log.Fatalf("Could not close channel.")
		}
		log.Print("Channel closed.")
	}
}

/*
Disconnect from the AMQP server and stop consuming messages.
*/
func (c *Consumer) Stop() error {
	channel, _ := c.conn.Channel()
	for _, binding := range c.topology.Bindings() {
		queue := binding.Queue()
		err := channel.Cancel(queue.Tag(), false)
		if err != nil {
			return err
		}
	}
	c.conn.Close()
	return nil
}

/*
Simple message type so users of this library don't have to import amqp as well
*/
type Message struct {
	amqp.Delivery
}
