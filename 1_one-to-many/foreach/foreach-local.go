package foreach

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Foreach"
	"github.com/nickng/scribble-go-examples/scributil"
)

func Local() {
	connAB, K := scributil.ParseFlags()
	wg := new(sync.WaitGroup)
	wg.Add(K)

	p := Foreach.New() // FIXME: K should be param here?
	for i := 1; i <= K; i++ {
		go Server_gather(p, K, i, connAB, connAB.Port(i), wg)
	}
	time.Sleep(100 * time.Millisecond)
	Client_scatter(p, K, 1, connAB, "localhost", connAB.BasePort())
	wg.Wait()
}
