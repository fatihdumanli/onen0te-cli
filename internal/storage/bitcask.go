package storage

import (
	"encoding/json"
	"fmt"

	"git.mills.io/prologic/bitcask"
	errors "github.com/pkg/errors"
)

type Bitcask struct{ Path string }

//Get the value of the given key.
func (b *Bitcask) Get(key string, obj interface{}) error {
	//TODO: return an error if the obj is not a ptr type
	var db, closer, err = b.openBitcaskDb()
	defer closer()
	if err != nil {
		return errors.Wrap(err, "couldn't open bitcask db")
	}

	bytes, err := db.Get([]byte(key))
	if err != nil {
		return errors.Wrap(err, "couldn't get the key")
	}

	err = json.Unmarshal(bytes, obj)
	if err != nil {
		return errors.Wrap(err, "couldn't unmarshal the value")
	}

	return nil
}

//Save a key-value pair to the local storage
func (b *Bitcask) Set(key string, val interface{}) error {

	var db, closer, err = b.openBitcaskDb()
	defer closer()
	if err != nil {
		return errors.Wrap(err, "couldn't open bitcask db")
	}

	bytes, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "couldn't serialize the data")
	}

	err = db.Put([]byte(key), []byte(string(bytes)))
	if err != nil {
		return errors.Wrap(err, "couldn't save the data")
	}
	return nil
}

//Remove the given key
func (b *Bitcask) Remove(key string) error {
	var db, closer, err = b.openBitcaskDb()
	defer closer()
	if err != nil {
		return errors.Wrap(err, "couldn't open bitcask db")
	}

	//If the key has not found
	if !db.Has([]byte(key)) {
		return fmt.Errorf("the key %s does not exist", key)
	}

	err = db.Delete([]byte(key))
	if err != nil {
		return errors.Wrap(err, "couldn't delete the key")
	}
	return err
}

func (b *Bitcask) GetKeys() (*[]string, error) {
	var db, closer, err = b.openBitcaskDb()
	var keys []string
	defer closer()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open bitcask db")
	}

	var keysChan = db.Keys()
	for i := 0; i < db.Len(); i++ {
		key := <-keysChan
		keys = append(keys, string(key))
	}

	return &keys, nil
}

func (b *Bitcask) openBitcaskDb() (*bitcask.Bitcask, func() error, error) {
	db, err := bitcask.Open(b.Path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't open bitcask db")
	}
	return db, db.Close, nil
}

func init() {
	//Make sure that Bitcask implements the Storer interface
	var x Storer
	x = &Bitcask{}
	_ = x
}
