package storage

import "errors"

var (
	CannotOpenDatabase     = errors.New("Cannot open the local database")
	CouldntGetTheKey       = errors.New("Cannot get the key")
	CouldntUnmarshalTheKey = errors.New("Couldn't unmarshal the key")
	CouldntSaveTheKey      = errors.New("Couldn't save the key")
)
