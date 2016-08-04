package store

type Store interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Del(key string) error

	Serialize() (map[string][]byte, error)
	Unserialize(map[string][]byte) (Store, error)
}
