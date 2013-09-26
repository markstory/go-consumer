# Go Consumer

GoConsumer is a wrapper around amqp that allows you to easily write
applications that consume AMQP queues. Exchange, Queue and bindings
are defined in a simple configuration file. GoConsumer handles plumbing
those together allowing you to focus on writing great consumer code.

This project borrows many ideas from [sparkplug](https://pypi.python.org/pypi/sparkplug/1.4)

## Configuration

Configuration is done through an ini style file. In this file you define the connection,
exchange and queue(s) you want to bind into. Once your file has been loaded and parsed you
can attach your consumer function in. The following is an example config file:

	[connection]
	host = localhost
	virtual_host = '/'
	user = guest
	password = guest
	ssl = False

	[exchange]
	name = grand
	type = direct
	durable = True
	auto_delete = False

	[queue]
	name = hose
	durable = True
	auto_delete = True
	exclusive = False
	routing_key = fire

Consuming functions need to have the following signature:

	func(*consumer.Message)

A simple example application would look like:

	import (
		"github.com/markstory/go-consumer"
		"log"
	)

	c, err = consumer.LoadConfig("./consumer.ini")
	if err != nil {
		log.Fatalf("Unable to create consumer. Error: %v", err)
	}
	err = c.Consume(func(msg *consumer.Message) {
		log.Print("Got a message")
		msg.Ack(true)
	})

Your consumer function will receive message types that can be acked
or nacked as you see fit.

## Signals

GoConsumer handles SIGINT, SIGTERM and SIGQUIT. In all cases the it attempts to shutdown
the AMQP connection and finish consuming any buffered messages.
