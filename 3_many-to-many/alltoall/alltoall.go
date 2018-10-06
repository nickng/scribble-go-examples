//go:generate scribblec-param.sh ../ManyToMany.scr -d ../ -param Alltoall github.com/nickng/scribble-go-examples/3_many-to-many/ManyToMany -param-api A -param-api B

package alltoall

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/3_many-to-many/ManyToMany/Alltoall"
	"github.com/nickng/scribble-go-examples/3_many-to-many/ManyToMany/Alltoall/A_1toM"
	"github.com/nickng/scribble-go-examples/3_many-to-many/ManyToMany/Alltoall/B_1toN"
	"github.com/nickng/scribble-go-examples/3_many-to-many/message"
	"github.com/nickng/scribble-go-examples/scributil"

	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
)

func init() {
	gob.Register(new(message.Data))
}

// A implements A[1]...A[M]
func A(p *Alltoall.Alltoall, M, N, selfM int, cc scributil.ClientConn, host string, port int, wg *sync.WaitGroup) {
	A := p.New_A_1toM(M, N, selfM)

	wgCli := new(sync.WaitGroup)
	wgCli.Add(N)
	mu := new(sync.Mutex) // lock is needed for concurrent access to connection map
	for n := 1; n <= N; n++ {
		go func(n int) {
			mu.Lock()
			if err := A.B_1toN_Dial(n, host, port+portOffset(selfM, M, n, N), cc.Dial, cc.Formatter()); err != nil {
				log.Fatalf("cannot connect: %v", err)
			}
			mu.Unlock()
			wgCli.Done()
		}(n)
	}
	wgCli.Wait()
	A.Run(func(s *A_1toM.Init) A_1toM.End {
		var d []message.Data
		for i := 0; i < N; i++ {
			d = append(d, message.Data{V: selfM})
		}
		sEnd := s.B_1toN_Scatter_Data(d)
		scributil.Debugf("A[%d]: sent %v.\n", selfM, d)
		return *sEnd
	})
	wg.Done()
}

// B implements B[1]...B[N]
func B(p *Alltoall.Alltoall, M, N, selfN int, sc scributil.ServerConn, port int, wg *sync.WaitGroup) {
	B := p.New_B_1toN(M, N, selfN)

	ss := make([]transport2.ScribListener, M)
	for m := 1; m <= M; m++ {
		ln, err := sc.Listen(port + portOffset(m, M, selfN, N))
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		ss[m-1] = ln
	}
	wgSvr := new(sync.WaitGroup)
	wgSvr.Add(M)
	mu := new(sync.Mutex) // lock is needed for concurrent access to connection map
	for m := 1; m <= M; m++ {
		go func(m int) {
			mu.Lock()
			if err := B.A_1toM_Accept(m, ss[m-1], sc.Formatter()); err != nil {
				log.Fatalf("cannot accept: %v", err)
			}
			mu.Unlock()
			wgSvr.Done()
		}(m)
	}
	wgSvr.Wait()
	B.Run(func(s *B_1toN.Init) B_1toN.End {
		d := make([]message.Data, M)
		sEnd := s.A_1toM_Gather_Data(d)
		scributil.Debugf("B[%d]: received %v.\n", selfN, d)
		return *sEnd
	})
	wg.Done()
}

func portOffset(m, M, n, N int) int {
	return n*M + m
}

// SplitMN splits parameter K to M and N as balanced as possible.
// It ensures M+N = K.
func SplitMN(K int) (M, N int) {
	return K / 2, K - K/2
}
