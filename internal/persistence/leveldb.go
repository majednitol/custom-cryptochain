package persistence

import "github.com/syndtr/goleveldb/leveldb"

type DB struct {
	db *leveldb.DB
}

func Open(path string) (*DB, error) {
	db, err := leveldb.OpenFile(path, nil)
	return &DB{db: db}, err
}

func (d *DB) Put(key, value []byte) error {
	return d.db.Put(key, value, nil)
}

func (d *DB) Get(key []byte) ([]byte, error) {
	return d.db.Get(key, nil)
}
