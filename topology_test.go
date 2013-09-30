package consumer

import (
	"testing"
)

const singleQueue = `
[connection]
host = localhost

[queue]
name = db_events
routing_key = events

[exchange]
name = events
`

func TestNewTopologyEmptyConfig(t *testing.T) {
	ini := ""
	conf := newConfig(ini)
	_, err := NewTopology(conf)
	if err == nil {
		t.Error("should fail, config is empty")
	}
}

func TestNewTopologyWithSingleQueue(t *testing.T) {
	conf := newConfig(singleQueue)
	top, err := NewTopology(conf)
	if err != nil {
		t.Error("Should not make an error")
	}
	if top.Connection().host != "localhost" {
		t.Error("host on connection does not match")
	}

}

func TestCreateTopologyWithMultipleQueue(t *testing.T) {
	/* conf := newConfig(multiQueue)*/
	t.Log("not done")
}

func TestCreateTopologyBindings(t *testing.T) {
	t.Log("not done")
}
