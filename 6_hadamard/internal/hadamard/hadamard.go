//go:generate scribblec-param.sh ../../Hadamard.scr -d ../../ -param Hadamard github.com/nickng/scribble-go-examples/6_hadamard/Hadamard -param-api A -param-api B -param-api C

package hadamard

import (
	"encoding/gob"
	"fmt"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/scributil"

	"github.com/nickng/scribble-go-examples/6_hadamard/Hadamard/Hadamard"
	"github.com/nickng/scribble-go-examples/6_hadamard/Hadamard/Hadamard/A_l1r1toK"
	"github.com/nickng/scribble-go-examples/6_hadamard/Hadamard/Hadamard/B_l1r1toK"
	"github.com/nickng/scribble-go-examples/6_hadamard/Hadamard/Hadamard/C_l1r1toK"
	"github.com/nickng/scribble-go-examples/6_hadamard/message"
	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
)

func init() {
	gob.Register(new(message.Val))
}

// A implements A[].
func A(p *Hadamard.Hadamard, N, self session2.Pair, cconn scributil.ClientConn, cHost string, cPort int, wg *sync.WaitGroup) {
	A := p.New_A_l1r1toK(N, self)

	scributil.Debugf("[connection] A[%s]: dialling to C[%s] at %s:%d.\n", self, self, cHost, cPort)
	if err := A.C_l1r1toK_Dial(self, cHost, cPort, cconn.Dial, cconn.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("A[%s]: Ready.\n", self)

	A.Run(func(s *A_l1r1toK.Init) A_l1r1toK.End {
		v := []message.Val{message.Val{V: self.Flatten(N)}}
		sEnd := s.C_selfplusl0r0_Scatter_Val(v)
		scributil.Debugf("A[%s]: sent %v.\n", self, v)
		return *sEnd
	})
	wg.Done()
}

// B implements B[].
func B(p *Hadamard.Hadamard, N, self session2.Pair, cconn scributil.ClientConn, cHost string, cPort int, wg *sync.WaitGroup) {
	B := p.New_B_l1r1toK(N, self)

	scributil.Debugf("[connection] B[%s]: dialling to C[%s] at %s:%d.\n", self, self, cHost, cPort)
	if err := B.C_l1r1toK_Dial(self, cHost, cPort, cconn.Dial, cconn.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("B[%s]: Ready.\n", self)

	B.Run(func(s *B_l1r1toK.Init) B_l1r1toK.End {
		v := []message.Val{message.Val{V: self.Flatten(N)}}
		sEnd := s.C_selfplusl0r0_Scatter_Val(v)
		scributil.Debugf("B[%s]: sent %v.\n", self, v)
		return *sEnd
	})
	wg.Done()
}

// C implements C[].
func C(p *Hadamard.Hadamard, N, self session2.Pair, aconn scributil.ServerConn, aconnPort int, bconn scributil.ServerConn, bconnPort int, wg *sync.WaitGroup) {
	C := p.New_C_l1r1toK(N, self)

	wgSvr := new(sync.WaitGroup)
	wgSvr.Add(2)
	go func() {
		lnA, err := aconn.Listen(aconnPort)
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer lnA.Close()
		scributil.Debugf("[connection] C[%s]: listening for A[%s] at :%d.\n", self, self, aconnPort)
		if err := C.A_l1r1toK_Accept(self, lnA, aconn.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
		wgSvr.Done()
	}()
	go func() {
		lnB, err := bconn.Listen(bconnPort)
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer lnB.Close()
		scributil.Debugf("[connection] C[%s]: listening for B[%s] at :%d.\n", self, self, bconnPort)
		if err := C.B_l1r1toK_Accept(self, lnB, bconn.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
		wgSvr.Done()
	}()
	wgSvr.Wait()
	scributil.Debugf("C[%s]: Ready.\n", self)

	C.Run(func(s *C_l1r1toK.Init) C_l1r1toK.End {
		v1 := make([]message.Val, 1)
		s0 := s.A_selfplusl0r0_Gather_Val(v1)
		scributil.Debugf("C[%s]: received %v.\n", self, v1)
		v2 := make([]message.Val, 1)
		sEnd := s0.B_selfplusl0r0_Gather_Val(v2)
		scributil.Debugf("C[%s]: received %v.\n", self, v2)
		fmt.Printf("The product at %s is %d.\n", self.String(), v1[0].V*v2[0].V)
		return *sEnd
	})
	wg.Done()
}
