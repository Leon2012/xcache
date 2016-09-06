package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Leon2012/xcache/cluster"
	logger "github.com/Leon2012/xcache/log"
	"github.com/valyala/fasthttp"
)

func init() {
	logger.SetModule("httpd")
}

type Httpd struct {
	addr    string
	cluster cluster.Cluster
}

func NewHttpd(addr string, c cluster.Cluster) *Httpd {
	return &Httpd{
		addr:    addr,
		cluster: c,
	}
}

func (s *Httpd) Start() error {
	go func() {
		fasthttp.ListenAndServe(s.addr, s.handleFastHTTP)
	}()
	return nil
}

func (s *Httpd) handleFastHTTP(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	if strings.HasPrefix(path, "/key") {
		s.handleKeyRequest(ctx)
	} else if path == "/join" {
		s.handleJoin(ctx)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}

func (s *Httpd) handleJoin(ctx *fasthttp.RequestCtx) {
	m := map[string]string{}
	buf := bytes.NewReader(ctx.PostBody())
	if err := json.NewDecoder(buf).Decode(&m); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if len(m) != 1 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	remoteAddr, ok := m["addr"]
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if err := s.cluster.Join(remoteAddr); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}

func (s *Httpd) handleKeyRequest(ctx *fasthttp.RequestCtx) {
	getKey := func() string {
		parts := strings.Split(string(ctx.Path()), "/")
		if len(parts) != 3 {
			return ""
		}
		return parts[2]
	}
	switch string(ctx.Method()) {
	case "GET":
		k := getKey()
		if k == "" {
			ctx.SetStatusCode(http.StatusBadRequest)
		}
		v, err := s.cluster.Get(k)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(map[string]string{k: string(v)})
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		//io.WriteString(w, string(b))
		ctx.WriteString(string(b))

	case "POST":
		// Read the value from the POST body.
		m := map[string]string{}
		buf := bytes.NewReader(ctx.PostBody())
		if err := json.NewDecoder(buf).Decode(&m); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		for k, v := range m {
			if err := s.cluster.Set(k, []byte(v), 0, 60); err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}
		}

	case "DELETE":
		k := getKey()
		if k == "" {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		if err := s.cluster.Del(k); err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		s.cluster.Del(k)

	default:
		ctx.SetStatusCode(http.StatusMethodNotAllowed)
	}
	return
}
