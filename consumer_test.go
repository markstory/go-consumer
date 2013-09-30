package consumer

import (
	"testing"
)

func TestCreateMissingFile(t *testing.T) {
	consumer, _ := Create("")
	if consumer != nil {
		t.Error("Should fail no file")
	}
}
