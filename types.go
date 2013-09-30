package consumer

import (
	"code.google.com/p/goconf/conf"
	"fmt"
)

type connection struct {
	host     string
	vhost    string
	user     string
	password string
	port     int
}

// Get the AMQP connection URL
func (c *connection) Url() string {
	return ""
}

func (c connection) String() string {
	return fmt.Sprintf("%#v", c)
}


type exchange struct {
	name       string
	kind       string
	durable    bool
	autoDelete bool
}

func (e exchange) String() string {
	return fmt.Sprintf("%#v", e)
}


type queue struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	routingKey string
}

func (q *queue) Name() string {
	return q.name
}

func (q *queue) Tag() string {
	return q.name + "-" + q.routingKey
}

func (q *queue) Exclusive() bool {
	return q.exclusive
}

func (q queue) String() string {
	return fmt.Sprintf("%#v", q)
}


/*
Convert the configuration file into domain objects
*/
func readConfigFile(config *conf.ConfigFile) (ex exchange, q queue, err error) {
	ex, err = newExchange(config)
	if err != nil {
		return
	}
	q, err = newQueue(config)
	if err != nil {
		return
	}
	return
}

/*
Create a new connection struct from the config file data.
*/
func newConnection(config *conf.ConfigFile) (c *connection, err error) {
	if !config.HasSection("connection") {
		return c, fmt.Errorf("Missing connection section in configuration file.")
	}
	c = &connection{
		host:     "localhost",
		vhost:    "/",
		user:     "guest",
		password: "guest",
		port:     5672,
	}
	if config.HasOption("connection", "host") {
		c.host, _ = config.GetString("connection", "host")
	}
	if config.HasOption("connection", "vhost") {
		c.vhost, _ = config.GetString("connection", "vhost")
	}
	if config.HasOption("connection", "user") {
		c.user, _ = config.GetString("connection", "user")
	}
	if config.HasOption("connection", "password") {
		c.password, _ = config.GetString("connection", "password")
	}
	if config.HasOption("connection", "port") {
		c.port, _ = config.GetInt("connection", "port")
	}
	return
}

/*
Create an exchange from the config file.
*/
func newExchange(config *conf.ConfigFile) (ex exchange, err error) {
	if !config.HasSection("exchange") {
		return ex, fmt.Errorf("Missing exchange section in configuration file.")
	}
	if _, err := config.GetString("exchange", "name"); err != nil {
		return ex, fmt.Errorf("Missing name from exchange section.")
	}
	name, _ := config.GetString("exchange", "name")
	kind := "direct"
	if config.HasOption("exchange", "type") {
		kind, _ = config.GetString("exchange", "type")
	}
	durable := true
	if config.HasOption("exchange", "durable") {
		durable, _ = config.GetBool("exchange", "durable")
	}
	autoDelete := false
	if config.HasOption("exchange", "auto_delete") {
		autoDelete, _ = config.GetBool("exchange", "auto_delete")
	}
	ex = exchange{
		name:       name,
		kind:       kind,
		durable:    durable,
		autoDelete: autoDelete,
	}
	return
}

/*
Create a queue from the config file.
*/
func newQueue(config *conf.ConfigFile) (q queue, err error) {
	if !config.HasSection("queue") {
		return q, fmt.Errorf("Missing queue section in configuration file.")
	}
	if _, err := config.GetString("queue", "name"); err != nil {
		return q, fmt.Errorf("Missing name from queue section.")
	}
	name, _ := config.GetString("queue", "name")
	durable := true
	if config.HasOption("queue", "durable") {
		durable, _ = config.GetBool("queue", "durable")
	}
	autoDelete := false
	if config.HasOption("queue", "auto_delete") {
		autoDelete, _ = config.GetBool("queue", "auto_delete")
	}
	exclusive := true
	if config.HasOption("queue", "exclusive") {
		exclusive, _ = config.GetBool("queue", "exclusive")
	}
	routingKey := ""
	if config.HasOption("queue", "routing_key") {
		routingKey, _ = config.GetString("queue", "routing_key")
	}
	q = queue{
		name:       name,
		durable:    durable,
		autoDelete: autoDelete,
		exclusive:  exclusive,
		routingKey: routingKey,
	}
	return
}
