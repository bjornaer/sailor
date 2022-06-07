package port

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type PortData struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

func (p *PortData) Serialize() (string, error) {
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err := e.Encode(p)
	if err != nil {
		return string(b.Bytes()[:]), errors.New("serialization failed")
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

func DeserializePort(str string) (*PortData, error) {
	var p PortData
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, errors.New("deserialization failed")
	}
	b := new(bytes.Buffer)
	b.Write(by)
	dec := gob.NewDecoder(b)
	if err := dec.Decode(&p); err != nil {
		return nil, errors.New("deserialization failed")
	}
	return &p, nil
}

type PortHandler func(portKey string, portData PortData)

func HandlePorts(fileName string, portHandler PortHandler) error {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error to read [file=%v]: %v", fileName, err.Error())
	}
	defer f.Close()
	r := bufio.NewReader(f)

	return decodeStream(r, portHandler)
}

func decodeStream(r io.Reader, portHandler PortHandler) error {
	dec := json.NewDecoder(r)

	// Expect start of object as the first token.
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('{') {
		return fmt.Errorf("expected {, got %v", t)
	}

	// While there are more tokens in the JSON stream...
	for dec.More() {

		// Read the port key.
		t, err := dec.Token()
		if err != nil {
			return err
		}
		portKey := t.(string) // type assert token to string.

		// Decode the value.
		var port PortData
		if err := dec.Decode(&port); err != nil {
			return err
		}
		// Execute callback function over port data
		portHandler(portKey, port)
	}
	return nil
}
