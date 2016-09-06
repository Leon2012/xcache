package service

import (
	"net"
	"sync"

	"github.com/Leon2012/xcache/cluster"
	logger "github.com/Leon2012/xcache/log"
)

var version string

func init() {
	logger.SetModule("memcached")
	version = "0.0.1"
}

type Memcached struct {
	listener  *net.TCPListener
	cluster   cluster.Cluster
	keepalive bool
	rcvBuf    int
	sendBuf   int
	sessions  []*Session
	mu        sync.Mutex
}

func NewMemcached(lis *net.TCPListener, cluster cluster.Cluster) *Memcached {
	return &Memcached{
		listener:  lis,
		cluster:   cluster,
		keepalive: true,
		rcvBuf:    1024,
		sendBuf:   1024,
	}
}

func (m *Memcached) Start() error {
	go func() {
		m.handleAccept()
	}()
	return nil
}

func (m *Memcached) handleAccept() {
	var (
		conn *net.TCPConn
		err  error
	)
	for {
		if conn, err = m.listener.AcceptTCP(); err != nil {
			logger.Error("listener.Accept(\"%s\") error(%v)", m.listener.Addr().String(), err)
			return
		}
		if err = conn.SetKeepAlive(m.keepalive); err != nil {
			logger.Error("conn.SetKeepAlive() error(%v)", err)
			return
		}
		if err = conn.SetReadBuffer(m.rcvBuf); err != nil {
			logger.Error("conn.SetReadBuffer() error(%v)", err)
			return
		}
		if err = conn.SetWriteBuffer(m.sendBuf); err != nil {
			logger.Error("conn.SetWriteBuffer() error(%v)", err)
			return
		}

		session := NewSession(conn, m)
		m.sessions = append(m.sessions, session)
		go session.handleMessage()
	}
}
