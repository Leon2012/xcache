package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/Leon2012/xcache/cluster/raft"
	"github.com/Leon2012/xcache/service"
	"github.com/Leon2012/xcache/store"
)

// Command line defaults
const (
	DefaultHTTPAddr = ":11000"
	DefaultRaftAddr = ":12000"
)

// Command line parameters
var httpAddr string
var raftAddr string
var joinAddr string

func init() {
	flag.StringVar(&httpAddr, "haddr", DefaultHTTPAddr, "Set the HTTP bind address")
	flag.StringVar(&raftAddr, "raddr", DefaultRaftAddr, "Set Raft bind address")
	flag.StringVar(&joinAddr, "join", "", "Set join address, if any")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <raft-data-path> \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "No Raft storage directory specified\n")
		os.Exit(1)
	}

	// Ensure Raft storage exists.
	raftDir := flag.Arg(0)
	if raftDir == "" {
		fmt.Fprintf(os.Stderr, "No Raft storage directory specified\n")
		os.Exit(1)
	}
	os.MkdirAll(raftDir, 0700)

	//s := store.NewStoreMem()
	//s, err := store.NewStoreLeveldb(filepath.Join(raftDir, "my.db"))
	s := store.NewRbTreeStore()
	// if err != nil {
	// 	log.Fatalf("failed to create store: %s", err.Error())
	// 	os.Exit(-1)
	// }
	//defer s.Close()

	r := raft.NewRaft(raftDir, raftAddr)
	if err := r.Init(joinAddr == "", s); err != nil {
		log.Fatalf("failed to open store: %s", err.Error())
	}

	//h := service.NewHttpd(httpAddr, r)
	pTCPAddr, _ := net.ResolveTCPAddr("tcp4", httpAddr)
	l, _ := net.ListenTCP("tcp", pTCPAddr)
	h := service.NewMemcached(l, r)
	if err := h.Start(); err != nil {
		log.Fatalf("failed to start HTTP service: %s", err.Error())
	}

	//If join was specified, make the join request.
	if joinAddr != "" {
		if err := join(joinAddr, raftAddr); err != nil {
			log.Fatalf("failed to join node at %s: %s", joinAddr, err.Error())
		}
	}

	log.Println("hraft started successfully")

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("hraftd exiting")
}

func join(joinAddr, raftAddr string) error {
	command := "join " + raftAddr + "\r\n"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", joinAddr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(command))
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

// func join(joinAddr, raftAddr string) error {
// 	b, err := json.Marshal(map[string]string{"addr": raftAddr})
// 	if err != nil {
// 		return err
// 	}
// 	resp, err := http.Post(fmt.Sprintf("http://%s/join", joinAddr), "application-type/json", bytes.NewReader(b))
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	return nil
// }
