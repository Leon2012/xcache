package store

type Store interface {
	Set(key string, value []byte, flag, expire int) error
	Get(key string) ([]byte, error)
	Del(key string) error
	Exist(key string) bool

	Serialize() (map[string][]byte, error)
	Unserialize(map[string][]byte) (Store, error)
}
