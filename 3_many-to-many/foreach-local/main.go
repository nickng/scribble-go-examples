package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/3_many-to-many/ManyToMany/Foreach"
	"github.com/nickng/scribble-go-examples/3_many-to-many/foreach"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	connAB, K := scributil.ParseFlags()
	M, N := foreach.SplitMN(K)
	scributil.Debugf("[info] K=%d: M=%d N=%d\n", K, M, N)
	protocol := Foreach.New()

	wg := new(sync.WaitGroup)
	wg.Add(M + N)
	for m := 1; m <= M; m++ {
		go func(m int) {
			time.Sleep(100 * time.Millisecond)
			foreach.A(protocol, M, N, m, connAB, "localhost", connAB.BasePort(), wg)
		}(m)
	}
	for n := 1; n <= N; n++ {
		go foreach.B(protocol, M, N, n, connAB, connAB.BasePort(), wg)
	}
	wg.Wait()
}
