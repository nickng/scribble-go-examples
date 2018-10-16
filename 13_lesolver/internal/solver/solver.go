//go:generate scribblec-param.sh ../../Solver.scr -d ../../ -param Solver github.com/nickng/scribble-go-examples/13_lesolver/Solver -param-api W
//go:generate scribblec-param.sh ../../Solver.scr -d ../../ -param Sync github.com/nickng/scribble-go-examples/13_lesolver/Solver -param-api W

package solver

import (
	"encoding/gob"
	"log"
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Solver"
	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Solver/family_1/W_l1r1plusl1r0toK_not_l1r1toKsubl1r0"
	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Solver/family_1/W_l1r1plusl1r0toKandl1r1toKsubl1r0"
	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Solver/family_1/W_l1r1toKsubl1r0_not_l1r1plusl1r0toK"
	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Sync"
	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Sync/family_1/W_l1r1toKsubl0r1_not_l1r2toK"
	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Sync/family_1/W_l1r1toKsubl0r1andl1r2toK"
	"github.com/nickng/scribble-go-examples/13_lesolver/Solver/Sync/family_1/W_l1r2toK_not_l1r1toKsubl0r1"
	"github.com/nickng/scribble-go-examples/13_lesolver/message"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func init() {
	gob.Register(new(message.Data))
}

var p2s = new(pipetosync)

// W1i implements W[1][i] (leftmost column).
// It also implements the sync protocol (as the lead columne).
func W1i(p *Solver.Solver, p2 *Sync.Sync, N, self session2.Pair, next scributil.ClientConn, nextHost string, nextPort int, syncConn scributil.ConnParam, syncBasePort int, wg *sync.WaitGroup) {
	p2s.init(N.Y)

	// Normal W[1][i] process.
	implW1i(p, N, self, next, nextHost, nextPort)

	// Sync phase.
	implW1Sync(p2, N, self, syncConn, nextHost, syncBasePort)
	wg.Done()
}

func implW1i(p *Solver.Solver, N, self session2.Pair, next scributil.ClientConn, nextHost string, nextPort int) {
	W1i := p.New_family_1_W_l1r1toKsubl1r0_not_l1r1plusl1r0toK(N, self)

	selfnext := self.Plus(session2.XY(1, 0))
	scributil.Debugf("[connection] W[%s]: dialling to W[%s] at %s:%d.\n", self, selfnext, nextHost, nextPort)
	if err := W1i.W_l1r1plusl1r0toKandl1r1toKsubl1r0_Dial(selfnext, nextHost, nextPort, next.Dial, next.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("W[%s]: Ready.\n", self)

	W1i.Run(func(s *W_l1r1toKsubl1r0_not_l1r1plusl1r0toK.Init) W_l1r1toKsubl1r0_not_l1r1plusl1r0toK.End {
		d := []message.Data{message.Data{V: self.X}}
		sEnd := s.W_selfplusl1r0_Scatter_Data(d)
		scributil.Debugf("W[%s]: sent %v.\n", self, d)
		return *sEnd
	})
}

// implW1Sync implements column 1 sync step.
func implW1Sync(p2 *Sync.Sync, N, self session2.Pair, syncConn scributil.ConnParam, syncHost string, syncBasePort int) {
	switch self.Y {
	case 1: // top left
		W11sync := p2.New_family_1_W_l1r1toKsubl0r1_not_l1r2toK(session2.XY(1, N.Y), self)
		selfnext := self.Plus(session2.XY(0, 1))
		time.Sleep(time.Duration((N.Y-self.Y)*100) * time.Millisecond)
		scributil.Debugf("W[%s]: dialling to W[%s] at %s:%d\n", self, selfnext, syncHost, syncBasePort+self.Y+1)
		if err := W11sync.W_l1r1toKsubl0r1andl1r2toK_Dial(selfnext, syncHost, syncBasePort+self.Y+1, syncConn.Dial, syncConn.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
		scributil.Debugf("W[%s]: Ready (sync).\n", self)
		W11sync.Run(func(s *W_l1r1toKsubl0r1_not_l1r2toK.Init) W_l1r1toKsubl0r1_not_l1r2toK.End {
			var d []message.Data
			d = []message.Data{message.Data{V: 1}}
			//d = <-p2s.get(self.Y)
			sEnd := s.W_selfplusl0r1_Scatter_Data(d)
			scributil.Debugf("W[%s]: sent %v (sync).\n", self, d)
			return *sEnd
		})
	case N.Y: // bottom left
		W1Ksync := p2.New_family_1_W_l1r2toK_not_l1r1toKsubl0r1(session2.XY(1, N.Y), self)
		lnSync, err := syncConn.Listen(syncBasePort + self.Y)
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer lnSync.Close()
		selfprev := self.Sub(session2.XY(0, 1))
		time.Sleep(time.Duration((N.Y-self.Y)*100) * time.Millisecond)
		scributil.Debugf("W[%s]: listening for W[%s] at :%d\n", self, selfprev, syncBasePort+self.Y)
		if err := W1Ksync.W_l1r1toKsubl0r1andl1r2toK_Accept(selfprev, lnSync, syncConn.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
		scributil.Debugf("W[%s]: Ready (sync).\n", self)
		W1Ksync.Run(func(s *W_l1r2toK_not_l1r1toKsubl0r1.Init) W_l1r2toK_not_l1r1toKsubl0r1.End {
			d := make([]message.Data, 1)
			sEnd := s.W_selfplusl0rneg1_Gather_Data(d)
			scributil.Debugf("W[%s]: received %v (sync).\n", self, d)
			return *sEnd
		})
	default: // in between
		W1isync := p2.New_family_1_W_l1r1toKsubl0r1andl1r2toK(session2.XY(1, N.Y), self)
		selfnext := self.Plus(session2.XY(0, 1))
		time.Sleep(time.Duration((N.Y-self.Y)*100) * time.Millisecond)
		scributil.Debugf("W[%s]: dialling to W[%s] at %s:%d\n", self, selfnext, syncHost, syncBasePort+self.Y+1)
		if selfnext.Y == N.Y {
			if err := W1isync.W_l1r2toK_not_l1r1toKsubl0r1_Dial(selfnext, syncHost, syncBasePort+self.Y+1, syncConn.Dial, syncConn.Formatter()); err != nil {
				log.Fatalf("cannot dial: %v", err)
			}
		} else {
			if err := W1isync.W_l1r1toKsubl0r1andl1r2toK_Dial(selfnext, syncHost, syncBasePort+self.Y+1, syncConn.Dial, syncConn.Formatter()); err != nil {
				log.Fatalf("cannot dial: %v", err)
			}
		}
		lnSync, err := syncConn.Listen(syncBasePort + self.Y)
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer lnSync.Close()
		selfprev := self.Sub(session2.XY(0, 1))
		scributil.Debugf("W[%s]: listening for W[%s] at :%d\n", self, selfprev, syncBasePort+self.Y)
		if selfprev.Y == 1 {
			if err := W1isync.W_l1r1toKsubl0r1_not_l1r2toK_Accept(selfprev, lnSync, syncConn.Formatter()); err != nil {
				log.Fatalf("cannot accept: %v", err)
			}
		} else {
			if err := W1isync.W_l1r1toKsubl0r1andl1r2toK_Accept(selfprev, lnSync, syncConn.Formatter()); err != nil {
				log.Fatalf("cannot accept: %v", err)
			}
		}
		scributil.Debugf("W[%s]: Ready (sync).\n", self)
		W1isync.Run(func(s *W_l1r1toKsubl0r1andl1r2toK.Init) W_l1r1toKsubl0r1andl1r2toK.End {
			d := make([]message.Data, 1)
			s0 := s.W_selfplusl0rneg1_Gather_Data(d)
			scributil.Debugf("W[%s]: received %v (sync).\n", self, d)
			sEnd := s0.W_selfplusl0r1_Scatter_Data(d)
			scributil.Debugf("W[%s]: sent %v (sync).\n", self, d)
			return *sEnd
		})
	}
}

// Wii implements W[i][i].
func Wii(p *Solver.Solver, N, self session2.Pair, prev scributil.ServerConn, prevPort int, next scributil.ClientConn, nextHost string, nextPort int, wg *sync.WaitGroup) {
	Wii := p.New_family_1_W_l1r1plusl1r0toKandl1r1toKsubl1r0(N, self)

	selfnext := self.Plus(session2.XY(1, 0))
	scributil.Debugf("[connection] W[%s]: dialling to W[%s] at %s:%d.\n", self, selfnext, nextHost, nextPort)
	if selfnext.X == N.X {
		if err := Wii.W_l1r1plusl1r0toK_not_l1r1toKsubl1r0_Dial(selfnext, nextHost, nextPort, next.Dial, next.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	} else {
		if err := Wii.W_l1r1plusl1r0toKandl1r1toKsubl1r0_Dial(selfnext, nextHost, nextPort, next.Dial, next.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	}
	ln, err := prev.Listen(prevPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	selfprev := self.Sub(session2.XY(1, 0))
	scributil.Debugf("[connection] W[%s]: listening for W[%s] at :%d.\n", self, selfprev, prevPort)
	if selfprev.X == 1 {
		if err := Wii.W_l1r1toKsubl1r0_not_l1r1plusl1r0toK_Accept(selfprev, ln, prev.Formatter()); err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
	} else {
		if err := Wii.W_l1r1plusl1r0toKandl1r1toKsubl1r0_Accept(selfprev, ln, prev.Formatter()); err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
	}
	scributil.Debugf("W[%s]: Ready.\n", self)

	Wii.Run(func(s *W_l1r1plusl1r0toKandl1r1toKsubl1r0.Init) W_l1r1plusl1r0toKandl1r1toKsubl1r0.End {
		d := make([]message.Data, 1)
		s0 := s.W_selfpluslneg1r0_Gather_Data(d)
		scributil.Debugf("W[%s]: received %v.\n", self, d)
		d = compute(d)
		sEnd := s0.W_selfplusl1r0_Scatter_Data(d)
		scributil.Debugf("W[%s]: sent %v.\n", self, d)
		return *sEnd
	})
	wg.Done()
}

// WKi implements W[K][i] (rightmost column).
func WKi(p *Solver.Solver, N, self session2.Pair, prev scributil.ServerConn, prevPort int, wg *sync.WaitGroup) {
	implWKi(p, N, self, prev, prevPort, nil)
	wg.Done()
}

func implWKi(p *Solver.Solver, N, self session2.Pair, prev scributil.ServerConn, prevPort int, data chan<- []message.Data) {
	WKK := p.New_family_1_W_l1r1plusl1r0toK_not_l1r1toKsubl1r0(N, self)

	ln, err := prev.Listen(prevPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	selfprev := self.Sub(session2.XY(1, 0))
	scributil.Debugf("[connection] W[%s]: listening for W[%s] at :%d.\n", self, selfprev, prevPort)
	if err := WKK.W_l1r1plusl1r0toKandl1r1toKsubl1r0_Accept(selfprev, ln, prev.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("W[%s]: Ready.\n", self)

	WKK.Run(func(s *W_l1r1plusl1r0toK_not_l1r1toKsubl1r0.Init) W_l1r1plusl1r0toK_not_l1r1toKsubl1r0.End {
		d := make([]message.Data, 1)
		sEnd := s.W_selfpluslneg1r0_Gather_Data(d)
		scributil.Debugf("W[%s]: received %v.\n", self, d)
		return *sEnd
	})
}

func compute(d []message.Data) []message.Data {
	return d
}

type pipetosync struct {
	mu    sync.Mutex
	pipes []chan []message.Data
}

func (p *pipetosync) get(i int) chan []message.Data {
	p.mu.Lock()
	ch := p.pipes[i-1]
	p.mu.Unlock()
	return ch
}

func (p *pipetosync) init(n int) {
	p.mu.Lock()
	if p.pipes == nil {
		p.pipes = make([]chan []message.Data, n)
		for i := range p.pipes {
			p.pipes[i] = make(chan []message.Data, 1)
		}
	}
	p.mu.Unlock()
}
