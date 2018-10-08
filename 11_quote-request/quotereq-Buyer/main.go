package main

import (
	"sync"

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
	wg.Add(1)
	// Must wait until all Supplier/Manufacturer are ready.
	quotereq.Buyer(protocol, S, 1,
		conn, "localhost", conn.BasePort(), // ->Supplier[1..S]
		wg)
	wg.Wait()
}
