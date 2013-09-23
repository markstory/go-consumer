package consumer

import (
	"fmt"
	"code.google.com/p/goconf/conf"
)

type exchange struct {
	name string
	kind string
	durable bool
	autoDelete bool
}

type binding struct {
	exchange string
	queue string
	routingKey string
}

type queue struct {
	name string
	durable bool
	autoDelete bool
	exclusive bool
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
		name: name,
		kind: kind,
		durable: durable,
		autoDelete: autoDelete,
	}
	return
}

/*
Create a queue from the config file.
*/
func newQueue(config *conf.ConfigFile) (q queue, err error) {
	return
}

/*
Create a binding from the config file.
*/
func newBinding(config *conf.ConfigFile) (bind binding, err error) {
	return
}
