package db

import "errors"

type MapDB struct {
	Storage map[string]string
}

func (m *MapDB) Get(k string) (string, error) {
	if v, found := m.Storage[k]; found {
		return v, nil
	} else {
		return "", errors.New("Key not found in Storage")
	}
}

func (m *MapDB) Put(k, v string) error {
	m.Storage[k] = v
	return nil
}

func NewMapDBClient() *MapDB {
	return &MapDB{Storage: make(map[string]string)}
}
