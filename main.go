package main

import (
	"log"
	"log/slog"
	"net"
)

const defaultListenAddr = ":5000"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitChan  chan struct{}
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr

	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitChan:  make(chan struct{}),
	}
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		case <-s.quitChan:
			return
		}
	}
}

func (s *Server) acceptLoop() error{
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error: %v", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.loop()

	slog.Info("server started on ", "listenAddr", s.ListenAddr)
	// go s.acceptLoop() // should it be a goroutine?
	// return nil
	return s.acceptLoop()
}


func (s *Server) handleConn(conn net.Conn) {
	// defer conn.Close()
	peer := NewPeer(conn)
	s.addPeerCh <- peer

	peer.readLoop()
}

func main() {
	server := NewServer(Config{})
	log.Fatal(server.Start())
}
