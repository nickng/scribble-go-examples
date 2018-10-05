package main

import (
	"sync"

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
	wg.Add(N)
	for n := 1; n <= N; n++ {
		go alltoall.B(protocol, M, N, n, connAB, connAB.BasePort(), wg)
	}
	wg.Wait()
}
