package cluster

type Cluster interface {
	Get(key string) ([]byte, error)
	Has(key string) bool
	Set(key string, value []byte, flag, expire int) error
	Add(key string, value []byte, flag, expire int) error
	Replace(key string, value []byte, flag, expire int) error
	Incr(key string, offset int64) (int64, error)
	Decr(key string, offset int64) (int64, error)
	Del(key string) error
	Join(addr string) error
	Name() string
	IsLeader() bool
}
