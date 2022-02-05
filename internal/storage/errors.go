package storage

import "errors"

var (
	CannotOpenDatabase     = errors.New("Cannot open the local database")
	KeyNotFound            = errors.New("key is not found")
	CouldntUnmarshalTheKey = errors.New("Couldn't unmarshal the key")
	CouldntSaveTheKey      = errors.New("Couldn't save the key")
)
