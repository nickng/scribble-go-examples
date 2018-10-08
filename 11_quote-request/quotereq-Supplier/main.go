package main

import (
	"flag"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/11_quote-request/QuoteRequest/WebService"
	"github.com/nickng/scribble-go-examples/11_quote-request/internal/quotereq"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	self := flag.Int("self", 1, "Supplier ID")
	conn, K := scributil.ParseFlags()
	protocol := WebService.New()
	M, S := (K - K/2), K/2
	scributil.Debugf("[info] M=%d S=%d.\n", M, S)
	if *self < 1 || S < *self {
		log.Fatalf("Supplier ID must be between 1 and %d (got %d).", S, *self)
	}

	// Port offset
	// 1..S            : Buyer -> Supplier[1],Supplier[2..S]
	// S+2..S+S        : Supplier[1]->Supplier[2..S] scatter channel
	// 2S+M+1..2S+M+M  : Supplier[1]->Manufacturer[1..M]
	// 2S+iM+1..2S+iM+M: Supplier[i]->Manufacturer[1..M]
	// 2S+SM+1..2S+SM+M: Supplier[S]->Manufacturer[1..M]

	wg := new(sync.WaitGroup)
	wg.Add(1)
	if *self == 1 {
		func(s int) {
			quotereq.Supplier1(protocol, M, S, s,
				conn, conn.Port(s), // Buyer->
				conn, "localhost", conn.Port(S), // ->baseport Supplier[1..S]
				conn, "localhost", conn.Port(2*S+s), // ->baseport Manufacturer[1..M]
				wg)
		}(*self)
	} else {
		func(s int) {
			quotereq.Supplier2toS(protocol, M, S, s,
				conn, conn.Port(s), // Buyer->
				conn, conn.Port(S+s), // Supplier[1]->
				conn, "localhost", conn.Port(2*S+s), // ->baseport Manufacturer[1..M]
				wg)
		}(*self)
	}
	wg.Wait()
}
