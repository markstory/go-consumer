package consumer

import (
	"code.google.com/p/goconf/conf"
	"testing"
)

const configFile = `
[connection]
vhost = /my-domain
user = mark
password = sekret
`

// Test generating URL with defaults.
func TestMakeAmqpUrlWithDefaults(t *testing.T) {
	config, _ := conf.ReadConfigBytes([]byte(""))
	url := makeAmqpUrl(config)
	if url != "amqp://guest:guest@localhost:5672/" {
		t.Error("URL with defaults is bad")
	}
}

func TestMakeAmqpUrl(t *testing.T) {
	config, _ := conf.ReadConfigBytes([]byte(configFile))
	url := makeAmqpUrl(config)
	if url != "amqp://mark:sekret@localhost:5672/my-domain" {
		t.Error("URL with defaults is bad")
	}
}

func TestCreateMissingFile(t *testing.T) {
	consumer, _ := Create("")
	if consumer != nil {
		t.Error("Should fail no file")
	}
}


