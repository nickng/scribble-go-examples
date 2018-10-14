//go:generate scribblec-param.sh ../../Jacobi.scr -d ../../ -param Jacobi github.com/nickng/scribble-go-examples/18_jacobi/Jacobi -param-api W

package jacobi

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/scributil"

	"github.com/nickng/scribble-go-examples/18_jacobi/Jacobi/Jacobi"
	"github.com/nickng/scribble-go-examples/18_jacobi/Jacobi/Jacobi/family_1/W_1to1_not_2to2and2toKsub1and3toK"
	"github.com/nickng/scribble-go-examples/18_jacobi/Jacobi/Jacobi/family_1/W_2to2and2toKsub1_not_1to1and3toK"
	"github.com/nickng/scribble-go-examples/18_jacobi/Jacobi/Jacobi/family_1/W_2toKsub1and3toK_not_1to1and2to2"
	"github.com/nickng/scribble-go-examples/18_jacobi/Jacobi/Jacobi/family_1/W_3toK_not_1to1and2to2and2toKsub1"
	"github.com/nickng/scribble-go-examples/18_jacobi/message"
)

func init() {
	gob.Register(new(message.Bound))
	gob.Register(new(message.Converged))
	gob.Register(new(message.Dimen))
}

// W1 implements W[1].
func W1(p *Jacobi.Jacobi, K, self int, w2 scributil.ClientConn, w2Host string, w2Port int, wg *sync.WaitGroup) {
	// W[1]\[2,2..K-1,3..K]
	W1 := p.New_family_1_W_1to1_not_2to2and2toKsub1and3toK(K, self)

	scributil.Debugf("[connection] W[%d]: dialling to W[%d] at %s:%d.\n", self, self+1, w2Host, w2Port)
	// W[2]
	if err := W1.W_2to2and2toKsub1_not_1to1and3toK_Dial(self+1, w2Host, w2Port, w2.Dial, w2.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("W[%d]: Ready.\n", self)

	W1.Run(func(s *W_1to1_not_2to2and2toKsub1and3toK.Init) W_1to1_not_2to2and2toKsub1and3toK.End {
		dimen := []message.Dimen{message.Dimen{}}
		nextDimen := make([]message.Dimen, 1)
		populate(dimen, 1, 1)     // populate subgrid[1]
		populate(nextDimen, 2, 2) // populate subgrid[2]
		scributil.Delay(1500)
		s0 := s.W_2_Scatter_Dimen(nextDimen)
		var (
			prevLeft    message.Bound
			prevRight   message.Bound
			left, right message.Bound
		)
		// calculate initial prevLeft and prevRight
		prevLeft, prevRight = initialBounds(dimen[0])
		for !converged(prevLeft, prevRight) {
			left, right = calculate(dimen[0], prevLeft, prevRight)

			scributil.Delay(1500)
			s1 := s0.W_2_Scatter_Bound([]message.Bound{right})
			scributil.Debugf("W[%d]: sent %v\n", self, right)
			b := make([]message.Bound, 1)
			s0 = s1.W_2_Gather_Bound(b)
			scributil.Debugf("W[%d]: received %v\n", self, b[0])
			prevRight = b[0]
			prevLeft = left
		}
		conv := []message.Converged{message.Converged{}}

		scributil.Delay(1500)
		sEnd := s0.W_2_Scatter_Converged(conv)
		scributil.Debugf("W[%d]: sent %v\n", self, conv)
		return *sEnd
	})
	wg.Done()
}

// W2 implements W[2].
func W2(p *Jacobi.Jacobi, K, self int, w1 scributil.ServerConn, w1Port int, wnext scributil.ClientConn, wnextHost string, wnextPort int, wall scributil.ClientConn, wallHost string, wallBasePort int, wg *sync.WaitGroup) {
	// W[2,2..K-1]\[1,3..K]
	W2 := p.New_family_1_W_2to2and2toKsub1_not_1to1and3toK(K, self)

	// Neighbour connections
	scributil.Debugf("[connection] W[%d]: dialling to W[%d] at %s:%d.\n", self, self+1, wnextHost, wnextPort)
	// W[3..K-1]
	if err := W2.W_2toKsub1and3toK_not_1to1and2to2_Dial(self+1, wnextHost, wnextPort, wnext.Dial, wnext.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	ln, err := w1.Listen(w1Port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	scributil.Debugf("[connection] W[%d]: listening for W[%d] at :%d.\n", self, 1, w1Port)
	// W[2..K-1]
	if err := W2.W_1to1_not_2to2and2toKsub1and3toK_Accept(self-1, ln, w1.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}

	// Scatter connections W[2] -> W[4..K]
	for k := 4; k <= K; k++ {
		scributil.Debugf("[connection/scatter] W[%d]: dialling to W[%d] at %s:%d.\n", self, k, wallHost, wallBasePort+k)
		if k < K {
			// W[3..K-1]
			if err := W2.W_2toKsub1and3toK_not_1to1and2to2_Dial(k, wallHost, wallBasePort+k, wall.Dial, wnext.Formatter()); err != nil {
				log.Fatalf("cannot dial: %v", err)
			}
		} else {
			// W[K]
			if err := W2.W_3toK_not_1to1and2to2and2toKsub1_Dial(k, wallHost, wallBasePort+k, wall.Dial, wall.Formatter()); err != nil {
				log.Fatalf("cannot dial: %v", err)
			}
		}
	}
	scributil.Debugf("W[%d]: Ready.\n", self)

	W2.Run(func(s *W_2to2and2toKsub1_not_1to1and3toK.Init) W_2to2and2toKsub1_not_1to1and3toK.End {
		dimen := []message.Dimen{message.Dimen{}}
		s0 := s.W_1_Gather_Dimen(dimen)
		nextDimen := make([]message.Dimen, K-3+1)
		populate(nextDimen, 3, K) // populate subgrid[3..K]

		scributil.Delay(1500)
		s1 := s0.W_3toK_Scatter_Dimen(nextDimen)
		var (
			prevLeft    message.Bound
			prevRight   message.Bound
			left, right message.Bound
		)
		// calculate initial prevLeft and prevRight
		prevLeft, prevRight = initialBounds(dimen[0])
		for {
			switch s2 := s1.W_1_Branch().(type) {
			case *W_2to2and2toKsub1_not_1to1and3toK.Bound:
				left, right = calculate(dimen[0], prevLeft, prevRight)
				s3 := s2.Recv_Bound(&prevLeft)
				scributil.Debugf("W[%d]: received %v\n", self, prevLeft)

				scributil.Delay(1500)
				s4 := s3.W_1_Scatter_Bound([]message.Bound{left})
				scributil.Debugf("W[%d]: sent %v\n", self, left)

				scributil.Delay(1500)
				s5 := s4.W_3_Scatter_Bound([]message.Bound{right})
				scributil.Debugf("W[%d]: sent %v\n", self, right)
				b := make([]message.Bound, 1)
				s1 = s5.W_3_Gather_Bound(b)
				prevRight = b[0]
				scributil.Debugf("W[%d]: received %v\n", self, prevRight)
			case *W_2to2and2toKsub1_not_1to1and3toK.Converged:
				var conv message.Converged
				scributil.Delay(1500)
				sEnd := s2.Recv_Converged(&conv).W_3_Scatter_Converged([]message.Converged{conv})
				scributil.Debugf("W[%d]: received %v\n", self, conv)
				scributil.Debugf("W[%d]: sent %v\n", self, conv)
				return *sEnd
			}
		}
	})
	wg.Done()
}

// Wi implements W[3..K-1].
func Wi(p *Jacobi.Jacobi, K, self int, wprev scributil.ServerConn, wprevPort int, wnext scributil.ClientConn, wnextHost string, wnextPort int, w2 scributil.ServerConn, w2Port int, wg *sync.WaitGroup) {
	// W[2..K-1,3..K]\[1,2]
	Wi := p.New_family_1_W_2toKsub1and3toK_not_1to1and2to2(K, self)

	// Neighbour connections
	scributil.Debugf("[connection] W[%d]: dialling to W[%d] at %s:%d.\n", self, self+1, wnextHost, wnextPort)
	if self+1 < K {
		// 3..K-1
		if err := Wi.W_2toKsub1and3toK_not_1to1and2to2_Dial(self+1, wnextHost, wnextPort, wnext.Dial, wnext.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	} else {
		// K
		if err := Wi.W_3toK_not_1to1and2to2and2toKsub1_Dial(self+1, wnextHost, wnextPort, wnext.Dial, wnext.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	}
	ln, err := wprev.Listen(wprevPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	scributil.Debugf("[connection] W[%d]: listening for W[%d] at :%d.\n", self, self-1, wprevPort)
	if self-1 == 2 {
		// 2
		if err := Wi.W_2to2and2toKsub1_not_1to1and3toK_Accept(self-1, ln, wprev.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	} else {
		// 3..K-1
		if err := Wi.W_2toKsub1and3toK_not_1to1and2to2_Accept(self-1, ln, wprev.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	}

	if self > 3 {
		// Scatter connection
		ln2, err := w2.Listen(w2Port)
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer ln2.Close()
		scributil.Debugf("[connection] W[%d]: listening for W[%d] at :%d.\n", self, 2, w2Port)
		// W[2]
		if err := Wi.W_2to2and2toKsub1_not_1to1and3toK_Accept(2, ln2, w2.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	}
	scributil.Debugf("W[%d]: Ready.\n", self)

	Wi.Run(func(s *W_2toKsub1and3toK_not_1to1and2to2.Init) W_2toKsub1and3toK_not_1to1and2to2.End {
		dimen := []message.Dimen{message.Dimen{}}
		s0 := s.W_2_Gather_Dimen(dimen)
		var (
			prevLeft    message.Bound
			prevRight   message.Bound
			left, right message.Bound
		)
		// calculate initial prevLeft and prevRight
		prevLeft, prevRight = initialBounds(dimen[0])
		for {
			switch s1 := s0.W_selfsub1_Branch().(type) {
			case *W_2toKsub1and3toK_not_1to1and2to2.Bound_W_State6:
				left, right = calculate(dimen[0], prevLeft, prevRight)
				s2 := s1.Recv_Bound(&prevLeft)
				scributil.Debugf("W[%d]: received %v\n", self, prevLeft)
				scributil.Delay(1500)
				s3 := s2.W_selfsub1_Scatter_Bound([]message.Bound{left})
				scributil.Debugf("W[%d]: sent %v\n", self, left)

				scributil.Delay(1500)
				s4 := s3.W_selfplus1_Scatter_Bound([]message.Bound{right})
				scributil.Debugf("W[%d]: sent %v\n", self, right)
				b := make([]message.Bound, 1)
				s0 = s4.W_selfplus1_Gather_Bound(b)
				prevRight = b[0]
				scributil.Debugf("W[%d]: received %v\n", self, prevRight)
			case *W_2toKsub1and3toK_not_1to1and2to2.Converged_W_State6:
				var conv message.Converged

				scributil.Delay(1500)
				sEnd := s1.Recv_Converged(&conv).W_selfplus1_Scatter_Converged([]message.Converged{conv})
				scributil.Debugf("W[%d]: received %v\n", self, conv)
				scributil.Debugf("W[%d]: sent %v\n", self, conv)
				return *sEnd
			}
		}
	})
	wg.Done()
}

// Wk implements W[K].
func Wk(p *Jacobi.Jacobi, K, self int, wlast scributil.ServerConn, wlastPort int, w2 scributil.ServerConn, w2Port int, wg *sync.WaitGroup) {
	// W[3..K]\[1,2,2..K-1]
	Wk := p.New_family_1_W_3toK_not_1to1and2to2and2toKsub1(K, self)

	// Neighbour connection
	ln, err := wlast.Listen(wlastPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	scributil.Debugf("[connection] W[%d]: listening for W[%d] at :%d.\n", self, self-1, wlastPort)
	// 3..K-1
	if err := Wk.W_2toKsub1and3toK_not_1to1and2to2_Accept(self-1, ln, wlast.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}

	// Scatter connection
	if self > 3 {
		ln2, err := w2.Listen(w2Port)
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer ln2.Close()
		scributil.Debugf("[connection] W[%d]: listening for W[%d] at :%d.\n", self, 2, w2Port)
		// 2
		if err := Wk.W_2to2and2toKsub1_not_1to1and3toK_Accept(2, ln2, w2.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	}
	scributil.Debugf("W[%d]: Ready.\n", self)

	Wk.Run(func(s *W_3toK_not_1to1and2to2and2toKsub1.Init) W_3toK_not_1to1and2to2and2toKsub1.End {
		dimen := []message.Dimen{message.Dimen{}}
		s0 := s.W_2_Gather_Dimen(dimen)
		var (
			prevLeft    message.Bound
			prevRight   message.Bound
			left, right message.Bound
		)
		// calculate initial prevLeft and prevRight
		prevLeft, prevRight = initialBounds(dimen[0])
		for {
			switch s1 := s0.W_selfsub1_Branch().(type) {
			case *W_3toK_not_1to1and2to2and2toKsub1.Bound_W_State2:
				left, right = calculate(dimen[0], prevLeft, prevRight)
				s2 := s1.Recv_Bound(&prevLeft)
				scributil.Debugf("W[%d]: received %v\n", self, prevLeft)

				scributil.Delay(1500)
				s0 = s2.W_selfsub1_Scatter_Bound([]message.Bound{left})
				scributil.Debugf("W[%d]: sent %v\n", self, left)
				prevRight = right
			case *W_3toK_not_1to1and2to2and2toKsub1.Converged_W_State2:
				var conv message.Converged
				sEnd := s1.Recv_Converged(&conv)
				scributil.Debugf("W[%d]: received %v\n", self, conv)
				return *sEnd
			}
		}
	})
	wg.Done()
}

var count = 0

// converged returns true if data is converged.
func converged(l, r message.Bound) bool {
	count++
	return count > 10 // assume converge after 10 runs
}

const (
	height = 4
	width  = 4
)

func populate(dim []message.Dimen, first, last int) []message.Dimen {
	// load partition[first] into dim[0]
	// load partition[last] into dim[last-first+1]
	for i := first; i <= last; i++ {
		dim[i-first].Height = height
		dim[i-first].Width = width
		dim[i-first].Grid = make([][]float32, dim[i-first].Height)
		for h := 0; h < dim[i-first].Height; h++ {
			dim[i-first].Grid[h] = make([]float32, dim[i-first].Width)
			for w := 0; w < dim[i-first].Width; w++ {
				dim[i-first].Grid[h][w] = float32(i)
			}
		}
	}
	return dim
}

// initialBounds extract the bounds of the sub-grid
// to be used as initial bounds.
func initialBounds(dim message.Dimen) (left, right message.Bound) {
	for h := 0; h < dim.Height; h++ {
		left.Bounds = append(left.Bounds, dim.Grid[h][0])
		right.Bounds = append(right.Bounds, dim.Grid[h][dim.Width-1])
	}
	return left, right
}

// calculate computes the next iteration of dim
// using prevLeft and prevRight.
func calculate(dim message.Dimen, prevLeft, prevRight message.Bound) (left, right message.Bound) {
	// Loop through the sub-grid, perform computation using prevLeft and prevRight
	for h := 0; h < dim.Height; h++ {
		left.Bounds = append(left.Bounds, (dim.Grid[h][0]+prevLeft.Bounds[h])/2)
		right.Bounds = append(right.Bounds, (dim.Grid[h][dim.Width-1]+prevRight.Bounds[h])/2)
	}
	return left, right
}
