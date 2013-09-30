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
	return fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		c.user,
		c.password,
		c.host,
		c.port,
		c.vhost)
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
	ex, err = newExchange(config, "exchange")
	if err != nil {
		return
	}
	q, err = newQueue(config, "queue")
	if err != nil {
		return
	}
	return
}

/*
Create a new connection struct from the config file data.
*/
func newConnection(config *conf.ConfigFile) (c connection, err error) {
	if !config.HasSection("connection") {
		return c, fmt.Errorf("Missing connection section in configuration file.")
	}
	c = connection{
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
func newExchange(config *conf.ConfigFile, section string) (ex exchange, err error) {
	if !config.HasSection(section) {
		return ex, fmt.Errorf("Missing exchange section in configuration file.")
	}
	if _, err := config.GetString(section, "name"); err != nil {
		return ex, fmt.Errorf("Missing name from exchange section.")
	}
	ex = exchange{
		name: "",
		kind: "direct",
		durable: true,
		autoDelete: false,
	}
	ex.name, _ = config.GetString(section, "name")
	if config.HasOption("exchange", "type") {
		ex.kind, _ = config.GetString(section, "type")
	}
	if config.HasOption(section, "durable") {
		ex.durable, _ = config.GetBool(section, "durable")
	}
	if config.HasOption(section, "auto_delete") {
		ex.autoDelete, _ = config.GetBool(section, "auto_delete")
	}
	return
}

/*
Create a queue from the config file.
*/
func newQueue(config *conf.ConfigFile, section string) (q queue, err error) {
	if !config.HasSection(section) {
		return q, fmt.Errorf("Missing queue section in configuration file.")
	}
	if _, err := config.GetString(section, "name"); err != nil {
		return q, fmt.Errorf("Missing name from queue section.")
	}
	name, _ := config.GetString(section, "name")
	q = queue{
		name:       name,
		durable:    true,
		autoDelete: false,
		exclusive:  true,
		routingKey: "",
	}
	if config.HasOption(section, "durable") {
		q.durable, _ = config.GetBool(section, "durable")
	}
	if config.HasOption(section, "auto_delete") {
		q.autoDelete, _ = config.GetBool(section, "auto_delete")
	}
	if config.HasOption(section, "exclusive") {
		q.exclusive, _ = config.GetBool(section, "exclusive")
	}
	if config.HasOption(section, "routing_key") {
		q.routingKey, _ = config.GetString(section, "routing_key")
	}
	return
}
