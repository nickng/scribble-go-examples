//go:generate scribblec-param.sh ../../Mesh.scr -d ../../ -param Mesh1 github.com/nickng/scribble-go-examples/7_mesh/Mesh -param-api W
//go:generate scribblec-param.sh ../../Mesh.scr -d ../../ -param Mesh3 github.com/nickng/scribble-go-examples/7_mesh/Mesh -param-api W

package mesh

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Mesh1"
	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Mesh1/family_1/W_l1r1plusl1r0toKhw_not_l1r1toKhwsubl1r0"
	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Mesh1/family_1/W_l1r1plusl1r0toKhwandl1r1toKhwsubl1r0"
	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Mesh1/family_1/W_l1r1toKhwsubl1r0_not_l1r1plusl1r0toKhw"
	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Mesh3"
	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Mesh3/family_1/W_l1r1toK1wsubl0r1_not_l1r2toK1w"
	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Mesh3/family_1/W_l1r1toK1wsubl0r1andl1r2toK1w"
	"github.com/nickng/scribble-go-examples/7_mesh/Mesh/Mesh3/family_1/W_l1r2toK1w_not_l1r1toK1wsubl0r1"
	"github.com/nickng/scribble-go-examples/7_mesh/message"
	"github.com/nickng/scribble-go-examples/scributil"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func init() {
	gob.Register(new(message.Data))
}

// W11 implements W[1][1].
func W11(p *Mesh1.Mesh1, N, self session2.Pair, next scributil.ClientConn, nextHost string, nextPort int, wg *sync.WaitGroup) {
	W11 := p.New_family_1_W_l1r1toKhwsubl1r0_not_l1r1plusl1r0toKhw(N, self)

	selfnext := self.Plus(session2.XY(1, 0))
	scributil.Debugf("[connection] W[%s]: dialling to W[%s] at %s:%d.\n", self, selfnext, nextHost, nextPort)
	if selfnext.Eq(N) {
		if err := W11.W_l1r1plusl1r0toKhw_not_l1r1toKhwsubl1r0_Dial(selfnext, nextHost, nextPort, next.Dial, next.Formatter()); err != nil {
			scributil.Debugf("cannot dial: %v", err)
		}
	} else {
		if err := W11.W_l1r1plusl1r0toKhwandl1r1toKhwsubl1r0_Dial(selfnext, nextHost, nextPort, next.Dial, next.Formatter()); err != nil {
			scributil.Debugf("cannot dial: %v", err)
		}
	}
	scributil.Debugf("W[%s]: Ready.\n", self)

	W11.Run(func(s *W_l1r1toKhwsubl1r0_not_l1r1plusl1r0toKhw.Init) W_l1r1toKhwsubl1r0_not_l1r1plusl1r0toKhw.End {
		d := []message.Data{message.Data{V: self.X}}
		sEnd := s.W_selfplusl1r0_Scatter_Data(d)
		scributil.Debugf("W[%s]: sent %v.\n", self, d)
		return *sEnd
	})
	wg.Done()
}

// Wi1 implements W[i][1].
func Wi1(p *Mesh1.Mesh1, N, self session2.Pair, prev scributil.ServerConn, prevPort int, next scributil.ClientConn, nextHost string, nextPort int, wg *sync.WaitGroup) {
	W1i := p.New_family_1_W_l1r1plusl1r0toKhwandl1r1toKhwsubl1r0(N, self)

	selfnext := self.Plus(session2.XY(1, 0))
	scributil.Debugf("[connection] W[%s]: dialling to W[%s] at %s:%d.\n", self, selfnext, nextHost, nextPort)
	if selfnext.Eq(N) {
		if err := W1i.W_l1r1plusl1r0toKhw_not_l1r1toKhwsubl1r0_Dial(selfnext, nextHost, nextPort, next.Dial, next.Formatter()); err != nil {
			scributil.Debugf("cannot dial: %v", err)
		}
	} else {
		if err := W1i.W_l1r1plusl1r0toKhwandl1r1toKhwsubl1r0_Dial(selfnext, nextHost, nextPort, next.Dial, next.Formatter()); err != nil {
			scributil.Debugf("cannot dial: %v", err)
		}
	}
	ln, err := prev.Listen(prevPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	selfprev := self.Sub(session2.XY(1, 0))
	scributil.Debugf("[connection] W[%s]: listening for W[%s] at :%d.\n", self, selfprev, prevPort)
	if selfprev.Eq(session2.XY(1, 1)) {
		if err := W1i.W_l1r1plusl1r0toKhw_not_l1r1toKhwsubl1r0_Accept(selfprev, ln, prev.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	} else {
		if err := W1i.W_l1r1plusl1r0toKhwandl1r1toKhwsubl1r0_Accept(selfprev, ln, prev.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	}
	scributil.Debugf("W[%s]: Ready.\n", self)

	W1i.Run(func(s *W_l1r1plusl1r0toKhwandl1r1toKhwsubl1r0.Init) W_l1r1plusl1r0toKhwandl1r1toKhwsubl1r0.End {
		d := make([]message.Data, 1)
		s0 := s.W_selfpluslneg1r0_Gather_Data(d)
		scributil.Debugf("W[%s]: received %v.\n", self, d)
		sEnd := s0.W_selfplusl1r0_Scatter_Data(d)
		scributil.Debugf("W[%s]: sent %v.\n", self, d)
		return *sEnd
	})
	wg.Done()
}

// WK1 implements W[K][1].
func WK1(p *Mesh1.Mesh1, N, self session2.Pair, prev scributil.ServerConn, prevPort int, wg *sync.WaitGroup) {
	WK1 := p.New_family_1_W_l1r1plusl1r0toKhw_not_l1r1toKhwsubl1r0(N, self)

	ln, err := prev.Listen(prevPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	selfprev := self.Sub(session2.XY(1, 0))
	scributil.Debugf("[connection] W[%s]: listening for W[%s] at :%d.\n", self, selfprev, prevPort)
	if selfprev.Eq(session2.XY(1, 1)) {
		if err := WK1.W_l1r1toKhwsubl1r0_not_l1r1plusl1r0toKhw_Accept(selfprev, ln, prev.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	} else {
		if err := WK1.W_l1r1plusl1r0toKhwandl1r1toKhwsubl1r0_Accept(selfprev, ln, prev.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	}
	scributil.Debugf("W[%s]: Ready.\n", self)

	WK1.Run(func(s *W_l1r1plusl1r0toKhw_not_l1r1toKhwsubl1r0.Init) W_l1r1plusl1r0toKhw_not_l1r1toKhwsubl1r0.End {
		d := make([]message.Data, 1)
		sEnd := s.W_selfpluslneg1r0_Gather_Data(d)
		scributil.Debugf("W[%s]: received %v.\n", self, d)
		return *sEnd
	})
	wg.Done()
}

// W11v implements W[1][1] (vertical loop).
func W11v(p *Mesh3.Mesh3, N, self session2.Pair, next scributil.ClientConn, nextHost string, nextPort int, wg *sync.WaitGroup) {
	W11v := p.New_family_1_W_l1r1toK1wsubl0r1_not_l1r2toK1w(N, self)

	selfnext := self.Plus(session2.XY(0, 1))
	scributil.Debugf("[connection] W[%s]: dialling to W[%s] at %s:%d.\n", self, selfnext, nextHost, nextPort)
	if selfnext.Eq(N) {
		if err := W11v.W_l1r2toK1w_not_l1r1toK1wsubl0r1_Dial(selfnext, nextHost, nextPort, next.Dial, next.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	} else {
		if err := W11v.W_l1r1toK1wsubl0r1andl1r2toK1w_Dial(selfnext, nextHost, nextPort, next.Dial, next.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	}
	scributil.Debugf("W[%s]: Ready.\n", self)

	W11v.Run(func(s *W_l1r1toK1wsubl0r1_not_l1r2toK1w.Init) W_l1r1toK1wsubl0r1_not_l1r2toK1w.End {
		d := []message.Data{message.Data{V: self.Y}}
		sEnd := s.W_selfplusl0r1_Scatter_Data(d)
		scributil.Debugf("W[%s]: sent %v.\n", self, d)
		return *sEnd
	})
	wg.Done()
}

// W1iv implements W[1][i] (vertical loop).
func W1iv(p *Mesh3.Mesh3, N, self session2.Pair, prev scributil.ServerConn, prevPort int, next scributil.ClientConn, nextHost string, nextPort int, wg *sync.WaitGroup) {
	W1iv := p.New_family_1_W_l1r1toK1wsubl0r1andl1r2toK1w(N, self)

	selfnext := self.Plus(session2.XY(0, 1))
	scributil.Debugf("[connection] W[%s]: dialling to W[%s] at %s:%d.\n", self, selfnext, nextHost, nextPort)
	if selfnext.Eq(N) {
		if err := W1iv.W_l1r2toK1w_not_l1r1toK1wsubl0r1_Dial(selfnext, nextHost, nextPort, next.Dial, next.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	} else {
		if err := W1iv.W_l1r1toK1wsubl0r1andl1r2toK1w_Dial(selfnext, nextHost, nextPort, next.Dial, next.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	}
	ln, err := prev.Listen(prevPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	selfprev := self.Sub(session2.XY(0, 1))
	scributil.Debugf("[connection] W[%s]: listening for W[%s] at :%d.\n", self, selfprev, prevPort)
	if selfprev.Eq(session2.XY(1, 1)) {
		if err := W1iv.W_l1r1toK1wsubl0r1_not_l1r2toK1w_Accept(selfprev, ln, prev.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	} else {
		if err := W1iv.W_l1r1toK1wsubl0r1andl1r2toK1w_Accept(selfprev, ln, prev.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	}
	scributil.Debugf("W[%s]: Ready.\n", self)

	W1iv.Run(func(s *W_l1r1toK1wsubl0r1andl1r2toK1w.Init) W_l1r1toK1wsubl0r1andl1r2toK1w.End {
		d := make([]message.Data, 1)
		s0 := s.W_selfplusl0rneg1_Gather_Data(d)
		scributil.Debugf("W[%s]: received %v.\n", self, d)
		sEnd := s0.W_selfplusl0r1_Scatter_Data(d)
		scributil.Debugf("W[%s]: sent %v.\n", self, d)
		return *sEnd
	})
	wg.Done()
}

// W1Kv implements W[1][K] (vertical loop).
func W1Kv(p *Mesh3.Mesh3, N, self session2.Pair, prev scributil.ServerConn, prevPort int, wg *sync.WaitGroup) {
	W1Kv := p.New_family_1_W_l1r2toK1w_not_l1r1toK1wsubl0r1(N, self)

	ln, err := prev.Listen(prevPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	selfprev := self.Sub(session2.XY(0, 1))
	scributil.Debugf("[connection] W[%s]: listening for W[%s] at :%d.\n", self, selfprev, prevPort)
	if selfprev.Eq(session2.XY(1, 1)) {
		if err := W1Kv.W_l1r1toK1wsubl0r1_not_l1r2toK1w_Accept(selfprev, ln, prev.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	} else {
		if err := W1Kv.W_l1r1toK1wsubl0r1andl1r2toK1w_Accept(selfprev, ln, prev.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	}
	scributil.Debugf("W[%s]: Ready.\n", self)

	W1Kv.Run(func(s *W_l1r2toK1w_not_l1r1toK1wsubl0r1.Init) W_l1r2toK1w_not_l1r1toK1wsubl0r1.End {
		d := make([]message.Data, 1)
		sEnd := s.W_selfplusl0rneg1_Gather_Data(d)
		scributil.Debugf("W[%s]: received %v.\n", self, d)
		return *sEnd
	})
	wg.Done()
}
