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
auto_delete = true
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


func TestNewQueueError(t *testing.T) {
	c := newConfig("")
	_, err := newQueue(c)
	if err == nil {
		t.Error("Missing queue section should cause an error.")
	}
}

func TestNewQueueNoNameError(t *testing.T) {
	ini := `
[queue]
`
	c := newConfig(ini)
	_, err := newQueue(c)
	if err == nil {
		t.Error("Should fail on missing name.")
	}
}

func TestNewQueueDefaults(t *testing.T) {
	ini := `
[queue]
name = test
`
	c := newConfig(ini)
	q, err := newQueue(c)
	if err != nil {
		t.Error("Should not make an error")
	}
	if q.name != "test" {
		t.Error("name does not match")
	}
	if q.durable != true {
		t.Error("durable should default to true")
	}
	if q.autoDelete != false {
		t.Error("autoDelete should default to false")
	}
	if q.routingKey != "" {
		t.Error("routingKey should default to ''")
	}
}

func TestNewQueueValues(t *testing.T) {
	ini := `
[queue]
name = test
durable = false
auto_delete = true
routing_key = fire
exclusive = false
`
	c := newConfig(ini)
	q, err := newQueue(c)
	if err != nil {
		t.Error("Should not make an error")
	}
	if q.name != "test" {
		t.Error("name does not match")
	}
	if q.Name() != "test" {
		t.Error("name method is wrong.")
	}
	if q.durable != false {
		t.Error("durable is wrong")
	}
	if q.autoDelete != true {
		t.Error("autoDelete is wrong")
	}
	if q.exclusive != false {
		t.Error("exclusive is wrong")
	}
	if q.routingKey != "fire" {
		t.Error("routingKey is wrong")
	}
	if q.Tag() != "test-fire" {
		t.Error("tagname is wrong.")
	}
}

func TestNewConnectionErrorOnMissingSection(t *testing.T) {
	ini := ""
	conf := newConfig(ini)
	_, err := newConnection(conf)
	if err == nil {
		t.Error("Error should not be nil")
	}
}

func TestNewConnectionDefaults(t *testing.T) {
	ini := `
[connection]
name = test
`
	conf := newConfig(ini)
	c, err := newConnection(conf)
	if err != nil {
		t.Error("Should not make an error")
	}
	if c.host != "localhost" {
		t.Error("Invalid default")
	}
	if c.vhost != "/" {
		t.Error("Invalid default")
	}
	if c.user != "guest" {
		t.Error("Invalid default")
	}
	if c.password != "guest" {
		t.Error("Invalid default")
	}
	if c.port != 5672 {
		t.Error("Invalid default")
	}
}

func TestNewConnection(t *testing.T) {
	ini := `
[connection]
host = queue-server
vhost = /app
user = mark
password = sekret
port = 9000
`
	conf := newConfig(ini)
	c, err := newConnection(conf)
	if err != nil {
		t.Error("Should not make an error")
	}
	if c.host != "queue-server" {
		t.Error("Invalid value")
	}
	if c.vhost != "/app" {
		t.Error("Invalid value")
	}
	if c.user != "mark" {
		t.Error("Invalid value")
	}
	if c.password != "sekret" {
		t.Error("Invalid value")
	}
	if c.port != 9000 {
		t.Error("Invalid value")
	}
}

func TestConnectionUrl(t *testing.T) {
	ini := `
[connection]
name = test
`
	conf := newConfig(ini)
	c, _ := newConnection(conf)
	if c.Url() != "amqp://guest:guest@localhost:5672/" {
		t.Error("URL with defaults is bad")
	}
}
