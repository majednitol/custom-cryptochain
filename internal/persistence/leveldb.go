package persistence

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
)

type DB struct {
	db *leveldb.DB
}

func Open(path string) (*DB, error) {
	db, err := leveldb.OpenFile(path, nil)
	return &DB{db: db}, err
}

// Put saves a struct or byte slice to DB
func (d *DB) Put(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return d.db.Put([]byte(key), bytes, nil)
}

// Get retrieves and unmarshals a struct from DB
func (d *DB) Get(key string, out interface{}) error {
	bytes, err := d.db.Get([]byte(key), nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}
