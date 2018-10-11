// Package scributil contains helper functions
// for using the scribble-go runtime API.
package scributil

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
)

// ConnParam is a connection parameter object.
// Use methods provided to create new connections.
type ConnParam interface {
	Dial(host string, port int) (transport2.BinChannel, error)
	Listen(port int) (transport2.ScribListener, error)
	Formatter() session2.ScribMessageFormatter

	BasePort() int
	Port(offset int) int
	Param() int
}

type shmConn struct {
	basePort int
	paramVal int
}

func newShmConn(basePort int, K int) shmConn {
	return shmConn{basePort: basePort, paramVal: K}
}

func (c shmConn) Dial(host string, port int) (transport2.BinChannel, error) {
	return shm.Dial(host, port)
}

func (c shmConn) Listen(port int) (transport2.ScribListener, error) {
	return shm.Listen(port)
}

func (c shmConn) Formatter() session2.ScribMessageFormatter {
	return new(session2.PassByPointer)
}

func (c shmConn) BasePort() int       { return c.basePort }
func (c shmConn) Port(offset int) int { return c.basePort + offset }
func (c shmConn) Param() int          { return c.paramVal }

type tcpConn struct {
	basePort int
	paramVal int
}

func newTcpConn(basePort int, K int) tcpConn {
	return tcpConn{basePort: basePort, paramVal: K}
}

func (c tcpConn) Dial(host string, port int) (transport2.BinChannel, error) {
	return tcp.Dial(host, port)
}
func (c tcpConn) Listen(port int) (transport2.ScribListener, error) {
	return tcp.Listen(port)
}

func (c tcpConn) Formatter() session2.ScribMessageFormatter {
	return new(session2.GobFormatter)
}

func (c tcpConn) BasePort() int       { return c.basePort }
func (c tcpConn) Port(offset int) int { return c.basePort + offset }
func (c tcpConn) Param() int          { return c.paramVal }

// ServerConn expose only the server-side parameter of the ConnParam.
type ServerConn interface {
	Listen(port int) (transport2.ScribListener, error)
	Formatter() session2.ScribMessageFormatter
}

// ClientConn expose only the client-side parameters of the ConnParam.
type ClientConn interface {
	Dial(host string, port int) (transport2.BinChannel, error)
	Formatter() session2.ScribMessageFormatter
}

// ParseFlags parses command line flags and returns
// a connection parameter object.
func ParseFlags() (cparam ConnParam, K int, I int) {
	const (
		tcpTran     = "tcp"
		shmTran     = "shm"
		defaultPort = 33333
	)
	var (
		tran  string // transport
		port  int    // base port
		param int    // parameter value
		self  int    // self ID
	)
	flag.StringVar(&tran, "t", tcpTran, "transport: tcp/shm")
	flag.IntVar(&port, "port", defaultPort, "base port for sever sockets")
	flag.IntVar(&param, "K", 2, "K parameter value")
	flag.IntVar(&self, "I", 1, "self parameter value")
	flag.Parse()

	Debugf("[info] Transport selected: %s\n", tran)
	switch tran {
	case tcpTran:
		return newTcpConn(port, param), param, self
	case shmTran:
		return newShmConn(port, param), param, self
	default:
		log.Fatalf("unrecognised transport: %s", tran)
	}
	return nil, param, self
}

// Debugf prints a message if debug mode is on.
func Debugf(format string, args ...interface{}) {
	if os.Getenv("DEBUG") == "1" {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
