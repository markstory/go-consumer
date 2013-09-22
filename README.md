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

	[binding]
	queue = hose
	exchange = grand
	routing_key = fire

	[queue]
	name = hose
	durable = True
	auto_delete = True
	exclusive = False

Once you've started your program you can attach your consumer:

	import (
		"github.com/markstory/go-consumer"
	)

	c = consumer.LoadConfig("./consumer.ini")
	c.BindConsumer(MyFunc)

Your consumer function will receive message types that can be acked
or nacked as you see fit.

## Signals

GoConsumer handles SIGKILL and SIGQUIT. In both cases the process will be terminated
once the client buffers have drained and message processing is complete.
