package storage

import (
	"encoding/json"

	"git.mills.io/prologic/bitcask"
)

type Bitcask struct{}

func (b *Bitcask) Get(key string) (interface{}, error) {
	var db, closer, err = openBitcaskDb()
	defer closer()
	if err != nil {
		return nil, err
	}

	var obj interface{}
	bytes, err := db.Get([]byte(key))
	if err != nil {
		return nil, KeyNotFound
	}
	err = json.Unmarshal(bytes, &obj)
	if err != nil {
		return nil, CouldntUnmarshalTheKey
	}

	return obj, nil
}

func (b *Bitcask) Set(key string, val interface{}) error {

	var db, closer, err = openBitcaskDb()
	defer closer()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(val)
	err = db.Put([]byte(key), bytes)
	if err != nil {
		return CouldntSaveTheKey
	}
	return nil
}

func openBitcaskDb() (*bitcask.Bitcask, func() error, error) {
	db, err := bitcask.Open("/tmp/cnote")
	if err != nil {
		return nil, nil, CannotOpenDatabase
	}
	return db, db.Close, nil
}

func init() {
	//Make sure that Bitcask implements the Storer interface
	var x Storer
	x = &Bitcask{}
	_ = x
}
