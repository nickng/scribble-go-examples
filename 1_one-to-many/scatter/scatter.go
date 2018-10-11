//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Scatter github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api A -param-api B

package scatter

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter/A_1to1"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter/B_1toK"
	"github.com/nickng/scribble-go-examples/1_one-to-many/messages"
	"github.com/nickng/scribble-go-examples/scributil"
)

func init() {
	gob.Register(new(messages.Data))
}

func Server_gather(p *Scatter.Scatter, K int, self int, sc scributil.ServerConn, port int, wg *sync.WaitGroup) {
	ss, err := sc.Listen(port)
	if err != nil {
		log.Fatalf("Cannot listen: %v", err)
	}
	B := p.New_B_1toK(K, self)
	if err := B.A_1to1_Accept(1, ss, sc.Formatter()); err != nil {
		log.Fatal(err)
	}
	B.Run(func(s *B_1toK.Init) B_1toK.End {
		d := make([]messages.Data, 1)
		end := s.A_1_Gather_Data(d)
		scributil.Debugf("B[%d]: received %v.\n", self, d)
		return *end
	})
	wg.Done()
}

func Client_scatter(p *Scatter.Scatter, K int, self int, cc scributil.ClientConn, host string, port int) {
	A := p.New_A_1to1(K, 1)
	for i := 1; i <= K; i++ {
		if err := A.B_1toK_Dial(i, host, port+i, cc.Dial, cc.Formatter()); err != nil { // FIXME: nil pointer error if no server
			log.Fatal(err)
		}
	}
	A.Run(func(s *A_1to1.Init) A_1to1.End {
		var d []messages.Data
		for i := 0; i < K; i++ {
			d = append(d, messages.Data{V: i})
		}

		scributil.Debugf("A[%d]: sending %v\n", self, d)
		scributil.Delay(1500)

		end := s.B_1toK_Scatter_Data(d)

		scributil.Debugf("A[%d]: sent %v\n", self, d)
		return *end
	})
}
