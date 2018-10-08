package scatter

import (
	"log"
	"sync"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter/A_1to1"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter/B_1toK"
	"github.com/nickng/scribble-go-examples/1_one-to-many/messages"
)

func Server_gather(listen func(int) (transport2.ScribListener, error), fmtr func() session2.ScribMessageFormatter,
			port int,
			p *Scatter.Scatter, K int, self int, wg *sync.WaitGroup) {
	ss, err := listen(port)
	if err != nil {
		log.Fatalf("Cannot listen: %v", err)
	}
	B := p.New_B_1toK(K, self)
	if err := B.A_1to1_Accept(1, ss, fmtr()); err != nil {
		log.Fatal(err)
	}
	B.Run(func(s *B_1toK.Init) B_1toK.End {
		d := make([]messages.Data, 1)
		end := s.A_1to1_Gather_Data(d)  // CHECKME: check values correct?
		return *end
	})
	wg.Done()
}

func Client_scatter(dial func(host string, port int) (transport2.BinChannel, error), fmtr func() session2.ScribMessageFormatter,
			host string, port int,
			p *Scatter.Scatter, K int, self int) {
	A := p.New_A_1to1(K, 1)
	for i := 1; i <= K; i++ {
		if err := A.B_1toK_Dial(i, host, port+i, dial, fmtr()); err != nil {  // FIXME: nil pointer error if no server
			log.Fatal(err)
		}
	}
	A.Run(func(s *A_1to1.Init) A_1to1.End {
		var d []messages.Data
		for i := 0; i < K; i++ {
			d = append(d, messages.Data{V: i})
		}
		end := s.B_1toK_Scatter_Data(d)
		return *end
	})
}
