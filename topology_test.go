package consumer

import (
	"testing"
)

func TestNewTopologyEmptyConfig(t *testing.T) {
	ini := ""
	conf := newConfig(ini)
	_, err := NewTopology(conf)
	if err == nil {
		t.Error("should fail, config is empty")
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

func TestNewTopologyWithMultipleQueue(t *testing.T) {
	/* conf := newConfig(multiQueue)*/
	t.Log("not done")
}

func TestNewTopologyUnmatchedQueue(t *testing.T) {
	t.Log("not done")
}

func TestNewTopologyUnmatchedExchange(t *testing.T) {
	t.Log("not done")
}
