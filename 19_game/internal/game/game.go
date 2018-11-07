//go:generate scribblec-param.sh ../../Game.scr -d ../../ -param Game github.com/nickng/scribble-go-examples/19_game/Game -param-api A -param-api B -param-api C
//go:generate scribblec-param.sh ../../Game.scr -d ../../ -param Proto1 github.com/nickng/scribble-go-examples/19_game/Game -param-api P -param-api Q

package game

import (
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/19_game/Game/Game"
	"github.com/nickng/scribble-go-examples/19_game/Game/Game/A_1to1"
	"github.com/nickng/scribble-go-examples/19_game/Game/Game/B_1to1"
	"github.com/nickng/scribble-go-examples/19_game/Game/Game/C_1to1"
	"github.com/nickng/scribble-go-examples/19_game/Game/Proto1"
	"github.com/nickng/scribble-go-examples/19_game/Game/Proto1/P_1toK"
	"github.com/nickng/scribble-go-examples/19_game/Game/Proto1/Q_1to1"
	"github.com/nickng/scribble-go-examples/scributil"
)

// Q implements Q[1].
func Q(p *Proto1.Proto1, N, self int, pConn scributil.ClientConn, pHost string, pBasePort int, bConn scributil.ClientConn, bHost string, bBasePort int, cConn scributil.ClientConn, cHost string, cBasePort int, wg *sync.WaitGroup) {
	Q := p.New_Q_1to1(N, self)
	for i := 1; i <= N; i++ {
		scributil.Debugf("Q: dialling to P[%d] at %s:%d.\n", i, pHost, pBasePort+i)
		if err := Q.P_1toK_Dial(i, pHost, pBasePort+i, pConn.Dial, pConn.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	}
	scributil.Debugf("Q: Ready.\n")

	Q.Run(func(s *Q_1to1.Init) Q_1to1.End {
		var As []*A_1to1.Init
		for i := 1; i <= N; i++ {
			Game := Game.New()
			A := Game.New_A_1to1(1)
			scributil.Debugf("[connection] A: dialling to B at %s:%d.\n", bHost, bBasePort+i)
			if err := A.B_1to1_Dial(1, bHost, bBasePort+i, bConn.Dial, bConn.Formatter()); err != nil {
				log.Fatalf("cannot dial: %v", err)
			}
			scributil.Debugf("[connection] A: dialling to C at %s:%d.\n", cHost, cBasePort+i)
			if err := A.C_1to1_Dial(1, cHost, cBasePort+i, cConn.Dial, cConn.Formatter()); err != nil {
				log.Fatalf("cannot dial: %v", err)
			}
			As = append(As, A.Init())
		}
		sEnd := s.P_1toK_Scatter_Play(As)
		return *sEnd
	})
	wg.Done()
}

// P1K implements P[1..K].
func P1K(p *Proto1.Proto1, N, self int, pConn scributil.ServerConn, pPort int, wg *sync.WaitGroup) {
	P := p.New_P_1toK(N, self)

	ln, err := pConn.Listen(pPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	scributil.Debugf("[connection] P[%d]: listening for Q at :%d\n", self, pPort)
	if err := P.Q_1to1_Accept(1, ln, pConn.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("P[%d]: Ready.\n", self)

	P.Run(func(s *P_1toK.Init) P_1toK.End {
		As := make([]*A_1to1.Init, 1)
		sEnd := s.Q_1_Gather_Play(As)
		scributil.Debugf("P[%d]: received endpoint A.\n", self)
		A(As[0])
		return *sEnd
	})
	wg.Done()
}

// A implements A[1].
func A(s *A_1to1.Init) A_1to1.End {
	s0 := s.B_1_Scatter_Bar()
	scributil.Debugf("A: sent Bar.\n")
	sEnd := s0.C_1_Scatter_Bar()
	scributil.Debugf("A: sent Bar.\n")
	return *sEnd
}

// B implements B[1].
func B(p *Game.Game, N, self int, bConn scributil.ServerConn, bPort int, wg *sync.WaitGroup) {
	B := p.New_B_1to1(self)
	ln, err := bConn.Listen(bPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	scributil.Debugf("B: listening for A at :%d.\n", bPort)
	if err := B.A_1to1_Accept(1, ln, bConn.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("B: Ready.\n")

	B.Run(func(s *B_1to1.Init) B_1to1.End {
		for {
			switch s0 := s.A_1_Branch().(type) {
			case *B_1to1.Foo:
				s = s0.Recv_Foo()
				scributil.Debugf("B: received Foo.\n")
			case *B_1to1.Bar:
				sEnd := s0.Recv_Bar()
				scributil.Debugf("B: received Bar.\n")
				return *sEnd
			}
		}
	})
	wg.Done()
}

// C implements C[1].
func C(p *Game.Game, N, self int, cConn scributil.ServerConn, cPort int, wg *sync.WaitGroup) {
	C := p.New_C_1to1(self)
	ln, err := cConn.Listen(cPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	scributil.Debugf("C: listening for A at :%d.\n", cPort)
	if err := C.A_1to1_Accept(1, ln, cConn.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("C: Ready.\n")

	C.Run(func(s *C_1to1.Init) C_1to1.End {
		for {
			switch s0 := s.A_1_Branch().(type) {
			case *C_1to1.Foo_C_Init:
				s = s0.Recv_Foo()
				scributil.Debugf("C: received Foo.\n")
			case *C_1to1.Bar_C_Init:
				sEnd := s0.Recv_Bar()
				scributil.Debugf("C: received Bar.\n")
				return *sEnd
			}
		}
	})
	wg.Done()
}
