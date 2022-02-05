package storage

type Storer interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
}
