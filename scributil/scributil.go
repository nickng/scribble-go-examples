package scributil

import (
	"flag"
	"log"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
)

func ParseFlags() (func(int) (transport2.ScribListener, error),
			func(host string, port int) (transport2.BinChannel, error),
			func() session2.ScribMessageFormatter,
			int, int) {
	trnsprt := flag.String("t", "tcp", "Transport: tcp/shm")  // N.B. cannot deref here -- gives default value before parsing
  port:= flag.Int("port", 33333, "base port for server sockets")
	K := flag.Int("K", 2, "K parameter value")
  flag.Parse()

	var listen func(int) (transport2.ScribListener, error)
	var dial func(host string, port int) (transport2.BinChannel, error)
	var fmtr func() session2.ScribMessageFormatter
	switch *trnsprt {
	case "tcp":
		listen = tcp.BListen
		dial = tcp.Dial
		fmtr = func() session2.ScribMessageFormatter { return new(session2.GobFormatter) } 
	case "shm":
		listen = shm.BListen
		dial = shm.Dial
		fmtr = func() session2.ScribMessageFormatter { return new(session2.PassByPointer) }
	default:
		log.Fatal("Unknown transport: ", trnsprt)
	}
	return listen, dial, fmtr, *port, *K
}
