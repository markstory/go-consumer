package consumer

import (
	"code.google.com/p/goconf/conf"
	"fmt"
)

type exchange struct {
	name       string
	kind       string
	durable    bool
	autoDelete bool
}

type binding struct {
	exchange   string
	queue      string
	routingKey string
}

type queue struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
}

/*
Convert the configuration file into domain objects
*/
func readConfigFile(config *conf.ConfigFile) (ex exchange, q queue, bind binding, err error) {
	ex, err = newExchange(config)
	if err != nil {
		return
	}
	q, err = newQueue(config)
	if err != nil {
		return
	}
	bind, err = newBinding(config)
	if err != nil {
		return
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
	if config.HasOption("exchange", "autoDelete") {
		autoDelete, _ = config.GetBool("exchange", "autoDelete")
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
	if config.HasOption("queue", "autoDelete") {
		autoDelete, _ = config.GetBool("queue", "autoDelete")
	}
	exclusive := true
	if config.HasOption("queue", "exclusive") {
		durable, _ = config.GetBool("queue", "durable")
	}
	q = queue{
		name:       name,
		durable:    durable,
		autoDelete: autoDelete,
		exclusive:  exclusive,
	}
	return
}

/*
Create a binding from the config file.
*/
func newBinding(config *conf.ConfigFile) (bind binding, err error) {
	return
}
