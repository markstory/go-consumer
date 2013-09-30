package consumer

import (
	"code.google.com/p/goconf/conf"
)

type topology struct {
	conn connection
	bindings []binding
}

func (t *topology) Connection() connection {
	return t.conn
}

func (t *topology) Bindings() []binding {
	return t.bindings
}

type binding struct {
	name     string
	exchange exchange
	queue    queue
}


/*
Create a queue topology from a config file

Once created, a topology can be used
to create an AMQP connection and declare
the relevant exchanges, queues, and bindings.

*/
func NewTopology(config *conf.ConfigFile) (t topology, err error) {
	conn, err := newConnection(config)
	if err != nil {
		return
	}
	var binds []binding

	t = topology{conn: conn, bindings: binds}
	return
}
