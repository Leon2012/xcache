package cluster

type Cluster interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Del(key string) error
	Join(addr string) error
	Name() string
	IsLeader() bool
}
