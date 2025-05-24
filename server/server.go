package server

import (
	"bufio"
	"github.com/Novicehood/woodis/errors"
	"github.com/Novicehood/woodis/proto"
	"net"
	"strings"
	"sync"
)

type CmdFunc func(p *Peer, cmd string, args []string)

type Hook func(*Peer, string, ...string) bool

type Peer struct {
	writer *bufio.Writer
	closed bool
	mu     sync.Mutex
	wg     sync.WaitGroup
}

func (p *Peer) Close() {
	p.mu.Lock()
	p.closed = true
	p.mu.Unlock()
}

func (p *Peer) Flush() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.writer.Flush()
}

type Server struct {
	l         net.Listener
	mu        sync.Mutex
	wg        sync.WaitGroup
	CmdMap    map[string]CmdFunc
	peers     map[net.Conn]struct{}
	Ctx       interface{}
	preHook   Hook
	countConn int
}

func NewServer(addr string) *Server {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(errors.INIT_SERVER_ERROR)
	}
	s := Server{
		l:      l,
		CmdMap: map[string]CmdFunc{},
		peers:  map[net.Conn]struct{}{},
	}
	s.Start()
	return &s
}

func (s *Server) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.serve(s.l)
		s.mu.Lock()
		for peer, _ := range s.peers {
			peer.Close()
			delete(s.peers, peer)
		}
		s.mu.Unlock()
	}()
}

func (s *Server) serve(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		s.HandleConn(conn)
	}
}

func (s *Server) HandleConn(conn net.Conn) {
	s.wg.Add(1)
	s.mu.Lock()
	s.countConn++
	s.peers[conn] = struct{}{}
	s.mu.Unlock()
	go func() {
		defer s.wg.Done()
		defer conn.Close()
		s.ServePeer(conn)
	}()
}

func (s *Server) ServePeer(conn net.Conn) {
	reader := bufio.NewReader(conn)

	peer := &Peer{
		writer: bufio.NewWriter(conn),
	}
	readCh := make(chan []string)
	go func() {
		defer close(readCh)
		for {
			args, err := proto.ReadArgs(reader)
			if err != nil {
				peer.Close()
				return
			}
			readCh <- args
		}
	}()

	for args := range readCh {
		s.Dispatch(peer, args)
		peer.Flush()
		if peer.closed {
			return
		}
	}

}

func (s *Server) Dispatch(p *Peer, args []string) {

}

func (s *Server) Close() {
	s.mu.Lock()
	if s.l != nil {
		s.l.Close()
	}
	s.l = nil
	s.mu.Unlock()
	s.wg.Wait()
}

func (s *Server) Register(cmd string, f CmdFunc) {
	s.mu.Lock()
	cmd = strings.ToUpper(cmd)
	defer s.mu.Unlock()
	if _, ok := s.CmdMap[cmd]; ok {
		panic(errors.INIT_SERVER_ERROR)
	}
	s.CmdMap[cmd] = f
}
