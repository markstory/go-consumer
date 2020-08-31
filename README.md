我是光年实验室高级招聘经理。
我在github上访问了你的开源项目，你的代码超赞。你最近有没有在看工作机会，我们在招软件开发工程师，拉钩和BOSS等招聘网站也发布了相关岗位，有公司和职位的详细信息。
我们公司在杭州，业务主要做流量增长，是很多大型互联网公司的流量顾问。公司弹性工作制，福利齐全，发展潜力大，良好的办公环境和学习氛围。
公司官网是http://www.gnlab.com,公司地址是杭州市西湖区古墩路紫金广场B座，若你感兴趣，欢迎与我联系，
电话是0571-88839161，手机号：18668131388，微信号：echo 'bGhsaGxoMTEyNAo='|base64 -D ,静待佳音。如有打扰，还请见谅，祝生活愉快工作顺利。

# Go Consumer

GoConsumer is a wrapper around amqp that allows you to easily write
applications that consume AMQP queues. Exchange, Queue and bindings
are defined in a simple configuration file. GoConsumer handles plumbing
those together allowing you to focus on writing great consumer code.

This project borrows many ideas from [sparkplug](https://pypi.python.org/pypi/sparkplug/1.4)

## Configuration

Configuration is done through an ini style file. In this file you define the connection,
exchange and queue(s) you want to bind into. Once your file has been loaded and parsed you
can attach your consumer function in. The following is an simple configuration file:

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

You can also configure your consumer to draw from multiple queues as well:

	[connection]
	host = localhost
	virtual_host = '/'
	user = guest
	password = guest
	ssl = False

	[exchange-app]
	name = app
	type = direct
	durable = True
	auto_delete = False

	[queue-app]
	name = app
	durable = True
	auto_delete = True
	exclusive = False
	routing_key = app-event

	[exchange-backend]
	name = backend
	type = direct
	durable = True
	auto_delete = False

	[queue-backend]
	name = backend
	durable = True
	auto_delete = True
	exclusive = False
	routing_key = backend-event

When multiple queues are being bound, each `exchange` and `queue` section should be suffixed
with the same value. This defines the binding between the exchange and queue.

## Consumers

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
