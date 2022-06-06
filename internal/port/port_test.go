package port_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/bjornaer/sailor/internal/port"
	"github.com/stretchr/testify/assert"
)

func TestHandlePorts(t *testing.T) {
	content := []byte(`{
		"PORTKEY": {
			"name": "test",
			"city": "test",
			"country": "test",
			"alias": [],
			"regions": [],
			"coordinates": [
			  55.5136433,
			  25.4052165
			],
			"province": "test",
			"timezone": "test",
			"unlocs": [
			  "test"
			],
			"code": "00000"
		  }
	}`)
	tmpfile, err := ioutil.TempFile("", "test_file")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	var filePortKey string
	err = port.HandlePorts(tmpfile.Name(), func(portKey string, portData port.PortData) {
		filePortKey = portKey
	})
	assert.Equal(t, "PORTKEY", filePortKey)
}
