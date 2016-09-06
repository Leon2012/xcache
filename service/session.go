package service

import (
	"bufio"
	"net"
	"strconv"
	"strings"

	logger "github.com/Leon2012/xcache/log"
	"github.com/Leon2012/xcache/service/memcache"
)

type Session struct {
	memcached *Memcached
	conn      net.Conn
	reader    *bufio.Reader
	writer    *bufio.Writer
}

func NewSession(conn net.Conn, memcached *Memcached) *Session {
	return &Session{
		memcached: memcached,
		conn:      conn,
		reader:    bufio.NewReader(conn),
		writer:    bufio.NewWriter(conn),
	}
}

func (m *Session) handleMessage() {
	defer m.close()
	for {
		req, err := memcache.Read(m.reader)
		if err != nil {
			if mcerr, ok := err.(memcache.MCError); ok {
				m.writeError(mcerr)
				continue
			} else {
				logger.Error("found err: %s", err.Error())
				return
			}
		}
		switch req.Command {
		case "join":
			logger.Info("node %s will join", req.Key)
			err := m.memcached.cluster.Join(req.Key)
			if err != nil {
				logger.Error("node join error : %s", err)
			}
			return
			break
		case "quit":
			return
			break
		case "version":
			res := memcache.MCRes{
				Response: "VERSION " + version,
			}
			// m.writer.WriteString(res.Protocol())
			// m.writer.Flush()
			m.writeResponse(res)
			break
		default:
			res, err := m.do(req)
			if err != nil {
				if mcerr, ok := err.(memcache.MCError); ok {
					m.writeError(mcerr)
					continue
				} else {
					logger.Error("found err: %s", err.Error())
					return
				}
			} else {
				if !req.Noreply {
					//m.writer.WriteString(res.Protocol())
					//m.writer.Flush()
					m.writeResponse(res)
				}
			}
		}
	}
}

func (m *Session) do(req *memcache.MCReq) (memcache.MCRes, error) {
	// m.memcached.mu.Lock()
	// defer m.memcached.mu.Lock()

	var (
		err error          = nil
		res memcache.MCRes = memcache.MCRes{}
	)

	switch req.Command {
	case "set":
		err = m.memcached.cluster.Set(req.Key, req.Data, 0, int(req.ExpireAt))
		if err != nil {
			res.Response = "NOT_STORED"
			//mcerr := memcache.NewMCError(memcache.SERVER_ERROR, err.Error())
			//return nil, mcerr
		} else {
			res.Response = "STORED"
		}
		break
	case "add":
		err = m.memcached.cluster.Add(req.Key, req.Data, 0, int(req.ExpireAt))
		if err != nil {
			res.Response = "NOT_STORED"
		} else {
			res.Response = "STORED"
		}
		break
	case "replace":
		err = m.memcached.cluster.Replace(req.Key, req.Data, 0, int(req.ExpireAt))
		if err != nil {
			res.Response = "NOT_STORED"
		} else {
			res.Response = "STORED"
		}
		break
	case "delete":
		err = m.memcached.cluster.Del(req.Key)
		if err != nil {
			res.Response = "NOT_FOUND"
		} else {
			res.Response = "DELETED"
		}
		break
	case "get":
		data, err := m.memcached.cluster.Get(req.Key)
		//fmt.Println("get data:" + string(data))
		res.Response = "END"
		if err == nil {
			mcValue := memcache.MCValue{}
			mcValue.Data = data
			mcValue.Key = req.Key
			mcValue.Flags = "0"
			res.Values = append(res.Values, mcValue)
		}
		break
	case "incr":
		value, err := m.memcached.cluster.Incr(req.Key, req.Offset)
		if err != nil {
			res.Response = "NOT_FOUND"
		} else {
			res.Response = strconv.Itoa(int(value))
		}
		break
	case "decr":
		value, err := m.memcached.cluster.Decr(req.Key, req.Offset)
		if err != nil {
			res.Response = "NOT_FOUND"
		} else {
			res.Response = strconv.Itoa(int(value))
		}
		break
	}
	return res, nil
}

func (m *Session) writeError(err memcache.MCError) {
	logger.Error("found err: %s", err.Error())
	m.writer.WriteString(err.Error())
	m.writer.Flush()
}

func (m *Session) writeResponse(res memcache.MCRes) {
	protocol := res.Protocol()
	logger.Debug("response : %s", strings.Replace(protocol, "\r\n", " ", -1))
	m.writer.WriteString(protocol)
	m.writer.Flush()
}

func (m *Session) close() {
	m.conn.Close()
}
