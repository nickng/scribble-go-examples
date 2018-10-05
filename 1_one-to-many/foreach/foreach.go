//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Foreach github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api A -param-api B

package foreach

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Foreach"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Foreach/A_1to1"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Foreach/B_1toK"
	"github.com/nickng/scribble-go-examples/1_one-to-many/messages"
	"github.com/nickng/scribble-go-examples/scributil"
)

func init() {
	gob.Register(new(messages.Data))
}

func Server_gather(p *Foreach.Foreach, K int, self int, sc scributil.ServerConn, port int, wg *sync.WaitGroup) {
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

func Client_scatter(p *Foreach.Foreach, K int, self int, cc scributil.ClientConn, host string, port int) {
	A := p.New_A_1to1(K, self)
	for i := 1; i <= K; i++ {
		if err := A.B_1toK_Dial(i, host, port+i, cc.Dial, cc.Formatter()); err != nil { // FIXME: nil pointer error if no server
			log.Fatal(err)
		}
	}
	var d []messages.Data
	for i := 1; i <= K; i++ {
		d = append(d, messages.Data{V: i})
	}
	A.Run(func(s *A_1to1.Init) A_1to1.End {
		scributil.Debugf("A: sent %v.\n", d)
		end := s.Foreach(func(s *A_1to1.Init_7) A_1to1.End {
			end := s.B_I_Scatter_Data(d)
			d = d[1:]
			return *end
		})
		return *end
	})
}
