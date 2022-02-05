package storage

type Storer interface {
	Get(key string, obj interface{}) error
	Set(key string, value interface{}) error
	Remove(key string) error
}
