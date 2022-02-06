package storage

import (
	"encoding/json"
	"log"

	"git.mills.io/prologic/bitcask"
)

type Bitcask struct{}

//Get the value of the given key.
func (b *Bitcask) Get(key string, obj interface{}) error {
	//TODO: return an error if the obj is not a ptr type
	var db, closer, err = openBitcaskDb()
	defer closer()
	if err != nil {
		return nil
	}

	bytes, err := db.Get([]byte(key))

	if err != nil {
		return KeyNotFound
	}

	err = json.Unmarshal(bytes, obj)
	if err != nil {
		log.Fatal(err)
		return CouldntUnmarshalTheKey
	}

	return nil
}

//Save a key-value pair to the local storage
func (b *Bitcask) Set(key string, val interface{}) error {

	var db, closer, err = openBitcaskDb()
	defer closer()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(val)

	err = db.Put([]byte(key), []byte(string(bytes)))
	if err != nil {
		return CouldntSaveTheKey
	}
	return nil
}

//Remove the given key
func (b *Bitcask) Remove(key string) error {
	var db, closer, err = openBitcaskDb()
	defer closer()
	if err != nil {
		return err
	}

	//If the key has not found
	if !db.Has([]byte(key)) {
		return KeyNotFound
	}

	err = db.Delete([]byte(key))
	return err
}

func (b *Bitcask) GetKeys() (*[]string, error) {
	var db, closer, err = openBitcaskDb()
	var keys []string

	defer closer()
	if err != nil {
		return nil, err
	}

	var keysChan = db.Keys()
	for i := 0; i < db.Len(); i++ {
		key := <-keysChan
		keys = append(keys, string(key))
	}

	return &keys, nil
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
