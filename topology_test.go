package consumer

import (
	"testing"
	"strings"
)

func TestNewTopologyEmptyConfig(t *testing.T) {
	ini := ""
	conf := newConfig(ini)
	_, err := NewTopology(conf)
	if err == nil {
		t.Error("should fail, config is empty")
	}
	if !strings.Contains(err.Error(), "Missing connection section") {
		t.Errorf("Error message is wrong. Got %s", err)
	}
}

const singleQueue = `
[connection]
host = localhost

[queue]
name = db_events
routing_key = events

[exchange]
name = events
`

func TestNewTopologyWithSingleQueue(t *testing.T) {
	conf := newConfig(singleQueue)
	top, err := NewTopology(conf)
	if err != nil {
		t.Error("Should not make an error")
	}
	if top.Connection().host != "localhost" {
		t.Error("host on connection does not match")
	}
	bindings := top.Bindings()
	if len(bindings) != 1 {
		t.Error("incorrect bindings made")
	}
	if bindings[0].queue.name != "db_events" {
		t.Error("Incorrect name for first queue.")
	}
	if bindings[0].queue.routingKey != "events" {
		t.Error("Incorrect routingKey for first queue.")
	}
	if bindings[0].exchange.name != "events" {
		t.Error("Incorrect name for first exchange.")
	}
}

const multiQueue = `
[connection]
host = localhost

[queue-front]
name = front-events
routing_key = events

[queue-back]
name = back-events
routing_key = events

[exchange-back]
name = be-events
type = direct

[exchange-front]
name = fe-events
type = direct
`
func TestNewTopologyWithMultipleQueue(t *testing.T) {
	conf := newConfig(multiQueue)
	top, err := NewTopology(conf)
	if err != nil {
		t.Error("Should not make an error")
	}
	if top.Connection().host != "localhost" {
		t.Error("host on connection does not match")
	}
	bindings := top.Bindings()
	if len(bindings) != 2 {
		t.Error("incorrect bindings made")
	}
	if bindings[0].queue.name != "back-events" {
		t.Error("Incorrect name for first queue.")
	}
	if bindings[0].queue.routingKey != "events" {
		t.Error("Incorrect routingKey for first queue.")
	}
	if bindings[0].exchange.name != "be-events" {
		t.Error("Incorrect name for first exchange.")
	}

	if bindings[1].queue.name != "front-events" {
		t.Error("Incorrect name for second queue.")
	}
	if bindings[1].queue.routingKey != "events" {
		t.Error("Incorrect routingKey for second queue.")
	}
	if bindings[1].exchange.name != "fe-events" {
		t.Error("Incorrect name for second exchange.")
	}
}


const unmatchedQueue = `
[connection]
host = localhost

[queue-front]
name = front-events
routing_key = events

[queue-back]
name = back-events
routing_key = events

[exchange-back]
name = be-events
type = direct
`

func TestNewTopologyUnmatchedQueue(t *testing.T) {
	conf := newConfig(unmatchedQueue)
	_, err := NewTopology(conf)
	if err == nil {
		t.Error("Should make an error")
	}
	if !strings.Contains(err.Error(), "Exchange and Queue lengths do not match") {
		t.Errorf("Error message is wrong. Got %s", err)
	}
}

const unmatchedExchange = `
[connection]
host = localhost

[exchange-front]
name = fe-events
type = direct

[queue-back]
name = back-events
routing_key = events

[exchange-back]
name = be-events
type = direct
`

func TestNewTopologyUnmatchedExchange(t *testing.T) {
	conf := newConfig(unmatchedExchange)
	_, err := NewTopology(conf)
	if err == nil {
		t.Error("Should make an error")
	}
	if !strings.Contains(err.Error(), "Exchange and Queue lengths do not match") {
		t.Errorf("Error message is wrong. Got %s", err)
	}
}

const missingBits = `
[connection]
host = localhost

[exchange-f]
name = fe-events
type = direct

[queue-b]
name = back-events
routing_key = events
`

func TestNewTopologyMisnamedBits(t *testing.T) {
	conf := newConfig(missingBits)
	_, err := NewTopology(conf)
	if err == nil {
		t.Error("Should make an error")
	}
	if !strings.Contains(err.Error(), "Exchange and Queue names do not match") {
		t.Errorf("Error message is wrong. Got %s", err)
	}
}
