//go:generate scribblec-param.sh ../../Fib.scr -d ../../ -param Fibonacci github.com/nickng/scribble-go-examples/17_fibonacci/Fib -param-api Fib

package fibonacci

import (
	"fmt"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/scributil"

	"github.com/nickng/scribble-go-examples/17_fibonacci/Fib/Fibonacci"
	"github.com/nickng/scribble-go-examples/17_fibonacci/Fib/Fibonacci/family_1/Fib_1toKsub2_not_2toKsub1and3toK"
	"github.com/nickng/scribble-go-examples/17_fibonacci/Fib/Fibonacci/family_1/Fib_1toKsub2and2toKsub1_not_3toK"
	"github.com/nickng/scribble-go-examples/17_fibonacci/Fib/Fibonacci/family_1/Fib_1toKsub2and2toKsub1and3toK"
	"github.com/nickng/scribble-go-examples/17_fibonacci/Fib/Fibonacci/family_1/Fib_2toKsub1and3toK_not_1toKsub2"
	"github.com/nickng/scribble-go-examples/17_fibonacci/Fib/Fibonacci/family_1/Fib_3toK_not_1toKsub2and2toKsub1"
)

// F1 implements F[1].
func F1(p *Fibonacci.Fibonacci, K, self int, kplus2 scributil.ClientConn, kplus2Host string, kplus2Port int, wg *sync.WaitGroup) {
	F1 := p.New_family_1_Fib_1toKsub2_not_2toKsub1and3toK(K, self)

	scributil.Debugf("[connection] Fib[%d]: dialling to Fib[%d] at %s:%d.\n", self, self+2, kplus2Host, kplus2Port)
	// F[3..K-1]
	if err := F1.Fib_1toKsub2and2toKsub1and3toK_Dial(self+2, kplus2Host, kplus2Port, kplus2.Dial, kplus2.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("Fib[%d]: Ready.\n", self)

	F1.Run(func(s *Fib_1toKsub2_not_2toKsub1and3toK.Init) Fib_1toKsub2_not_2toKsub1and3toK.End {
		sEnd := s.Fib_selfplus2_Scatter_T([]int{0})
		return *sEnd
	})
	wg.Done()
}

// F2 implements F[2].
func F2(p *Fibonacci.Fibonacci, K, self int, kplus1 scributil.ClientConn, kplus1Host string, kplus1Port int, kplus2 scributil.ClientConn, kplus2Host string, kplus2Port int, wg *sync.WaitGroup) {
	F2 := p.New_family_1_Fib_1toKsub2and2toKsub1_not_3toK(K, self)

	scributil.Debugf("[connection] Fib[%d]: dialling to Fib[%d] at %s:%d.\n", self, self+2, kplus2Host, kplus2Port)
	// F[3..K-1]
	if err := F2.Fib_2toKsub1and3toK_not_1toKsub2_Dial(self+2, kplus2Host, kplus2Port, kplus2.Dial, kplus2.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("[connection] Fib[%d]: dialling to Fib[%d] at %s:%d.\n", self, self+1, kplus1Host, kplus1Port)
	// F[3..K-1]
	if err := F2.Fib_2toKsub1and3toK_not_1toKsub2_Dial(self+1, kplus1Host, kplus1Port, kplus1.Dial, kplus1.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("Fib[%d]: Ready.\n", self)

	F2.Run(func(s *Fib_1toKsub2and2toKsub1_not_3toK.Init) Fib_1toKsub2and2toKsub1_not_3toK.End {
		s0 := s.Fib_selfplus1_Scatter_T([]int{1})
		sEnd := s0.Fib_selfplus2_Scatter_T([]int{1})
		return *sEnd
	})
	wg.Done()
}

// Fi implements F[3..K-2].
func Fi(p *Fibonacci.Fibonacci, K, self int, ksub2 scributil.ServerConn, ksub2Port int, ksub1 scributil.ServerConn, ksub1Port int, kplus1 scributil.ClientConn, kplus1Host string, kplus1Port int, kplus2 scributil.ClientConn, kplus2Host string, kplus2Port int, wg *sync.WaitGroup) {
	Fi := p.New_family_1_Fib_1toKsub2and2toKsub1and3toK(K, self)

	scributil.Debugf("[connection] Fib[%d]: dialling to Fib[%d] at %s:%d.\n", self, self+2, kplus2Host, kplus2Port)
	// F[3..K-2]
	if err := Fi.Fib_1toKsub2and2toKsub1and3toK_Dial(self+2, kplus2Host, kplus2Port, kplus2.Dial, kplus2.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("[connection] Fib[%d]: dialling to Fib[%d] at %s:%d.\n", self, self+1, kplus1Host, kplus1Port)
	// F[3..K-2]
	if err := Fi.Fib_1toKsub2and2toKsub1and3toK_Dial(self+1, kplus1Host, kplus1Port, kplus1.Dial, kplus1.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	ln1, err := ksub1.Listen(ksub1Port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln1.Close()
	scributil.Debugf("[connection] Fib[%d]: listening for Fib[%d] at :%d.\n", self, self-1, ksub1Port)
	if self == 3 {
		// Fib[2]
		if err := Fi.Fib_1toKsub2and2toKsub1_not_3toK_Accept(self-1, ln1, ksub1.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	} else {
		// Fib[3..K-1]
		if err := Fi.Fib_1toKsub2and2toKsub1and3toK_Accept(self-1, ln1, ksub1.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	}
	ln2, err := ksub2.Listen(ksub2Port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln2.Close()
	scributil.Debugf("[connection] Fib[%d]: listening for Fib[%d] at :%d.\n", self, self-2, ksub2Port)
	if self == 3 {
		// Fib[1]
		if err := Fi.Fib_1toKsub2_not_2toKsub1and3toK_Accept(self-2, ln2, ksub2.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	} else {
		// Fib[3..K-1]
		if err := Fi.Fib_1toKsub2and2toKsub1and3toK_Accept(self-2, ln2, ksub2.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	}
	scributil.Debugf("Fib[%d]: Ready.\n", self)

	Fi.Run(func(s *Fib_1toKsub2and2toKsub1and3toK.Init) Fib_1toKsub2and2toKsub1and3toK.End {
		v := make([]int, 2)
		s0 := s.Fib_selfsub2_Gather_T(v)
		s1 := s0.Fib_selfsub1_Gather_T(v[1:])
		s2 := s1.Fib_selfplus1_Scatter_T([]int{v[0] + v[1]})
		sEnd := s2.Fib_selfplus2_Scatter_T([]int{v[0] + v[1]})
		return *sEnd
	})
	wg.Done()
}

// Fksub1 implements F[K-1].
func Fksub1(p *Fibonacci.Fibonacci, K, self int, ksub2 scributil.ServerConn, ksub2Port int, ksub1 scributil.ServerConn, ksub1Port int, kplus1 scributil.ClientConn, kplus1Host string, kplus1Port int, wg *sync.WaitGroup) {
	Fksub1 := p.New_family_1_Fib_2toKsub1and3toK_not_1toKsub2(K, self)

	scributil.Debugf("[connection] Fib[%d]: dialling to Fib[%d] at %s:%d.\n", self, self+1, kplus1Host, kplus1Port)
	// F[K]
	if err := Fksub1.Fib_3toK_not_1toKsub2and2toKsub1_Dial(self+1, kplus1Host, kplus1Port, kplus1.Dial, kplus1.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	ln1, err := ksub1.Listen(ksub1Port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln1.Close()
	scributil.Debugf("[connection] Fib[%d]: listening for Fib[%d] at :%d.\n", self, self-1, ksub1Port)
	// Fib[3..K-2]
	if err := Fksub1.Fib_1toKsub2and2toKsub1and3toK_Accept(self-1, ln1, ksub1.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	ln2, err := ksub2.Listen(ksub2Port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln2.Close()
	scributil.Debugf("[connection] Fib[%d]: listening for Fib[%d] at :%d.\n", self, self-2, ksub2Port)
	// Fib[3..K-2]
	if err := Fksub1.Fib_1toKsub2and2toKsub1and3toK_Accept(self-2, ln2, ksub2.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("Fib[%d]: Ready.\n", self)

	Fksub1.Run(func(s *Fib_2toKsub1and3toK_not_1toKsub2.Init) Fib_2toKsub1and3toK_not_1toKsub2.End {
		v := make([]int, 2)
		s0 := s.Fib_selfsub2_Gather_T(v)
		s1 := s0.Fib_selfsub1_Gather_T(v[1:])
		sEnd := s1.Fib_selfplus1_Scatter_T([]int{v[0] + v[1]})
		return *sEnd
	})
	wg.Done()
}

// Fk implements F[K].
func Fk(p *Fibonacci.Fibonacci, K, self int, ksub2 scributil.ServerConn, ksub2Port int, ksub1 scributil.ServerConn, ksub1Port int, wg *sync.WaitGroup) {
	Fk := p.New_family_1_Fib_3toK_not_1toKsub2and2toKsub1(K, self)

	// Establish connection in reverse: K-1 then K-2
	ln1, err := ksub1.Listen(ksub1Port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln1.Close()
	scributil.Debugf("[connection] Fib[%d]: listening for Fib[%d] at :%d.\n", self, self-1, ksub1Port)
	// Fib[K-1]
	if err := Fk.Fib_2toKsub1and3toK_not_1toKsub2_Accept(self-1, ln1, ksub1.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	ln2, err := ksub2.Listen(ksub2Port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln2.Close()
	scributil.Debugf("[connection] Fib[%d]: listening for Fib[%d] at :%d.\n", self, self-2, ksub2Port)
	// Fib[3..K-2]
	if err := Fk.Fib_1toKsub2and2toKsub1and3toK_Accept(self-2, ln2, ksub2.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("Fib[%d]: Ready.\n", self)

	Fk.Run(func(s *Fib_3toK_not_1toKsub2and2toKsub1.Init) Fib_3toK_not_1toKsub2and2toKsub1.End {
		v := make([]int, 2)
		s0 := s.Fib_selfsub2_Gather_T(v)
		sEnd := s0.Fib_selfsub1_Gather_T(v[1:])
		fmt.Printf("The Fib(%d) is %d.\n", self, v[0]+v[1])
		return *sEnd
	})
	wg.Done()
}
