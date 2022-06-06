package port

import (
	"bufio"
	"encoding/json"
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
