package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/3_many-to-many/ManyToMany/Alltoall"
	"github.com/nickng/scribble-go-examples/3_many-to-many/alltoall"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	connAB, K := scributil.ParseFlags()
	M, N := alltoall.SplitMN(K)
	scributil.Debugf("[info] K=%d: M=%d N=%d\n", K, M, N)
	protocol := Alltoall.New()

	wg := new(sync.WaitGroup)
	wg.Add(M + N)
	for m := 1; m <= M; m++ {
		go func(m int) {
			time.Sleep(100 * time.Millisecond)
			alltoall.A(protocol, M, N, m, connAB, "localhost", connAB.BasePort(), wg)
		}(m)
	}
	for n := 1; n <= N; n++ {
		go alltoall.B(protocol, M, N, n, connAB, connAB.BasePort(), wg)
	}
	wg.Wait()
}
