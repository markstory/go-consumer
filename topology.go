package consumer

import (
	"code.google.com/p/goconf/conf"
	"fmt"
	"sort"
	"strings"
)

type topology struct {
	conn     connection
	bindings []binding
}

func (t *topology) Connection() connection {
	return t.conn
}

func (t *topology) Bindings() []binding {
	return t.bindings
}


type binding struct {
	exchange exchange
	queue    queue
}

func (b *binding) Queue() queue {
	return b.queue
}

func (b *binding) Exchange() exchange {
	return b.exchange
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
	sections := config.GetSections()
	sort.Strings(sections)

	count, err := checkSections(sections)
	if err != nil {
		return
	}

	var (
		binds []binding
		queues []queue
		exchanges []exchange
	)

	for _, section := range sections {
		if strings.HasPrefix(section, "queue") {
			q, _ := newQueue(config, section)
			queues = append(queues, q)
		}
		if strings.HasPrefix(section, "exchange") {
			ex, _ := newExchange(config, section)
			exchanges = append(exchanges, ex)
		}
	}

	for i := 0; i < count; i++ {
		bind := binding{
			exchange: exchanges[i],
			queue: queues[i],
		}
		binds = append(binds, bind)
	}

	t = topology{conn: conn, bindings: binds}
	return
}

func checkSections(sections []string) (int, error) {
	qCount := 0
	exCount := 0
	for _, section := range sections {
		if strings.HasPrefix(section, "queue") {
			qCount += 1
		}
		if strings.HasPrefix(section, "exchange") {
			exCount += 1
		}
	}
	var err error
	if qCount != exCount {
		err = fmt.Errorf("Exchange and Queue sections do not match.")
	}
	return qCount, err
}
