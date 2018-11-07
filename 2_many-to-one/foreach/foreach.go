//go:generate scribblec-param.sh ../ManyToOne.scr -d ../ -param Foreach github.com/nickng/scribble-go-examples/2_many-to-one/ManyToOne -param-api A -param-api B

package foreach

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/2_many-to-one/ManyToOne/Foreach"
	"github.com/nickng/scribble-go-examples/2_many-to-one/ManyToOne/Foreach/A_1toK"
	"github.com/nickng/scribble-go-examples/2_many-to-one/ManyToOne/Foreach/B_1to1"
	"github.com/nickng/scribble-go-examples/2_many-to-one/message"
	"github.com/nickng/scribble-go-examples/scributil"

	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
)

func init() {
	gob.Register(new(message.Data))
}

var WaitConn chan struct{}

// A1toK is the gather sender.
func A1toK(p *Foreach.Foreach, K, self int, cc scributil.ClientConn, host string, port int, wg *sync.WaitGroup) {
	A := p.New_A_1toK(K, self)
	if WaitConn != nil {
		<-WaitConn
	}
	if err := A.B_1to1_Dial(1, host, port, cc.Dial, cc.Formatter()); err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	A.Run(func(s *A_1toK.Init) A_1toK.End {
		d := []message.Data{message.Data{V: self}}
		sEnd := s.B_1_Scatter_Data(d)
		scributil.Debugf("A[%d] sent %v.\n", self, d)
		return *sEnd
	})
	wg.Done()
}

// B is the gather receiver.
func B(p *Foreach.Foreach, K, self int, sc scributil.ServerConn, port int, wg *sync.WaitGroup) {
	B := p.New_B_1to1(K, self)

	ss := make([]transport2.ScribListener, K)
	for i := 1; i <= K; i++ {
		ln, err := sc.Listen(port + i)
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		ss[i-1] = ln
	}
	if WaitConn != nil {
		close(WaitConn)
	}
	wgSvr := new(sync.WaitGroup)
	wgSvr.Add(K)
	for i := 1; i <= K; i++ {
		go func(i int) {
			if err := B.A_1toK_Accept(i, ss[i-1], sc.Formatter()); err != nil {
				log.Fatalf("cannot accept: %v", err)
			}
			wgSvr.Done()
		}(i)
	}
	wgSvr.Wait()

	B.Run(func(s *B_1to1.Init) B_1to1.End {
		dK := make([]message.Data, K)
		i := 0
		sEnd := s.Foreach(func(s0 *B_1to1.Init_11) B_1to1.End {
			sEnd := s0.A_I_Gather_Data(dK[i:])
			i++
			return *sEnd
		})
		scributil.Debugf("B[%d] received %v.\n", self, dK)
		return *sEnd
	})
	wg.Done()
}
