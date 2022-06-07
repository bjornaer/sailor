package db

import (
	"time"
)

// DBClient interface just handles strings *but* I would totally use generics here
type DBClient interface {
	Put(k, v string) error
	Get(k string) (string, error)
}

// InitDBClient takes variadic args, if none is given it uses an in memory Map
// if "redis" is passed, a second argument for the address is required
func InitDBClient(clientType ...string) (DBClient, error) {
	if len(clientType) == 0 {
		// default use in memory map
		return NewMapDBClient()
	}
	// if not default use whatever came from variable first
	dbClientType := clientType[0]
	if dbClientType == "redis" {
		address := clientType[1]
		return NewRedisClient(address)
	} else {
		// for anything else we go back to the in memory map
		// Ideally we mighjjt have more DB type of options and handle those
		return NewMapDBClient()
	}
}

func InitDBClientWithRetry(retry int, clientType ...string) (DBClient, error) {
	i := 0
	var err error
	var client DBClient
	for i < retry {
		client, err = InitDBClient(clientType...)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
		i++
	}
	return client, err
}
