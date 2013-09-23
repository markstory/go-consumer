package consumer

import (
	"code.google.com/p/goconf/conf"
	"testing"
)

func newConfig(content string) (*conf.ConfigFile) {
	config, _ := conf.ReadConfigBytes([]byte(content))
	return config
}

func TestNewExchangeError(t *testing.T) {
	c := newConfig("")
	_, err := newExchange(c)
	if err == nil {
		t.Error("Missing exchange section should cause an error.")
	}
}

func TestNewExchangeNoNameError(t *testing.T) {
	ini := `
[exchange]
`
	c := newConfig(ini)
	_, err := newExchange(c)
	if err == nil {
		t.Error("Should fail on missing name.")
	}
}

func TestNewExchangeDefaults(t *testing.T) {
	ini := `
[exchange]
name = test
`
	c := newConfig(ini)
	ex, err := newExchange(c)
	if err != nil {
		t.Error("Should not make an error")
	}
	if ex.name != "test" {
		t.Error("name does not match")
	}
	if ex.kind != "direct" {
		t.Error("default value for kind is wrong")
	}
	if ex.durable != true {
		t.Error("durable should default to true")
	}
	if ex.autoDelete != false {
		t.Error("autoDelete should default to false")
	}
}

func TestNewExchangeValues(t *testing.T) {
	ini := `
[exchange]
name = test
type = fanout
durable = false
autoDelete = true
`
	c := newConfig(ini)
	ex, err := newExchange(c)
	if err != nil {
		t.Error("Should not make an error")
	}
	if ex.name != "test" {
		t.Error("name does not match")
	}
	if ex.kind != "fanout" {
		t.Error("kind is wrong")
	}
	if ex.durable != false {
		t.Error("durable is wrong")
	}
	if ex.autoDelete != true {
		t.Error("autoDelete is wrong")
	}
}