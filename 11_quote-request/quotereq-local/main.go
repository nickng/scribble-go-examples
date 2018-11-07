package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/11_quote-request/QuoteRequest/WebService"
	"github.com/nickng/scribble-go-examples/11_quote-request/internal/quotereq"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	conn, K := scributil.ParseFlags()
	protocol := WebService.New()
	M, S := (K - K/2), K/2
	scributil.Debugf("[info] M=%d S=%d.\n", M, S)

	// Port offset
	// 1..S            : Buyer -> Supplier[1],Supplier[2..S]
	// S+2..S+S        : Supplier[1]->Supplier[2..S] scatter channel
	// 2S+M+1..2S+M+M  : Supplier[1]->Manufacturer[1..M]
	// 2S+iM+1..2S+iM+M: Supplier[i]->Manufacturer[1..M]
	// 2S+SM+1..2S+SM+M: Supplier[S]->Manufacturer[1..M]

	wg := new(sync.WaitGroup)
	wg.Add(M + S + 1)
	go func() {
		// Must wait until all Supplier/Manufacturer are ready.
		time.Sleep(2 * time.Second)
		quotereq.Buyer(protocol, S, 1,
			conn, "localhost", conn.BasePort(), // ->Supplier[1..S]
			wg)
	}()
	// Spawn in reverse order of Supplier.
	for s := 1; s <= S; s++ {
		if s == 1 {
			go func(s int) {
				time.Sleep(time.Duration((S-s+1)*100) * time.Millisecond)
				quotereq.Supplier1(protocol, M, S, s,
					conn, conn.Port(s), // Buyer->
					conn, "localhost", conn.Port(S), // ->baseport Supplier[1..S]
					conn, "localhost", conn.Port(2*S+s), // ->baseport Manufacturer[1..M]
					wg)
			}(s)
		} else {
			go func(s int) {
				time.Sleep(time.Duration((S-s+1)*100) * time.Millisecond)
				quotereq.Supplier2toS(protocol, M, S, s,
					conn, conn.Port(s), // Buyer->
					conn, conn.Port(S+s), // Supplier[1]->
					conn, "localhost", conn.Port(2*S+s), // ->baseport Manufacturer[1..M]
					wg)
			}(s)
		}
	}
	// Must spawn first
	for m := 1; m <= M; m++ {
		go func(m int) {
			quotereq.Manufacturer(protocol, M, S, m,
				conn, conn.Port(2*S+(m-1)*S), // Supplier[1..S]->Manufacturer[m]
				wg)
		}(m)
	}
	wg.Wait()
}
