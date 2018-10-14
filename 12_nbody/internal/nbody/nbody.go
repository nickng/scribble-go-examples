//go:generate scribblec-param.sh ../../NBody.scr -d ../../ -param NBody github.com/nickng/scribble-go-examples/12_nbody/NBody -param-api W

package nbody

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/12_nbody/NBody/NBody"
	"github.com/nickng/scribble-go-examples/12_nbody/NBody/NBody/family_2/W_1to1_not_2to2and2toKsub1and3toKandKtoK"
	"github.com/nickng/scribble-go-examples/12_nbody/NBody/NBody/family_2/W_2toKsub1and3toK_not_1to1and2to2andKtoK"
	"github.com/nickng/scribble-go-examples/12_nbody/NBody/NBody/family_2/W_3toKandKtoK_not_1to1and2to2and2toKsub1"
	"github.com/nickng/scribble-go-examples/12_nbody/message"
	"github.com/nickng/scribble-go-examples/scributil"
)

// NIter is the number of iterations.
var NIter = 2

func init() {
	gob.Register(new(message.Particles))
	gob.Register(new(message.Stop))
}

// W1 is the main worker W[1] which also decides if the iteration should continue.
func W1(p *NBody.NBody, K, self int, w2 scributil.ClientConn, w2Host string, w2Port int, wn scributil.ClientConn, wNHost string, wNPort int, wg *sync.WaitGroup) {
	// W[1]\[2,2..K-1,3..K,K]
	W1 := p.New_family_2_W_1to1_not_2to2and2toKsub1and3toKandKtoK(K, self)
	particles := loadParticles("particles.txt", self)
	rcvd := []message.Particles{particles}
	var velocities message.Vector3 // Stores the velocities

	scributil.Debugf("[connection] W[%d]: dialling to W[%d] at %s:%d.\n", self, self+1, w2Host, w2Port)
	if err := W1.W_2toKsub1and3toK_not_1to1and2to2andKtoK_Dial(self+1, w2Host, w2Port, w2.Dial, w2.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("[connection] W[%d]: dialling to W[%d] at %s:%d.\n", self, K, wNHost, wNPort)
	if err := W1.W_3toKandKtoK_not_1to1and2to2and2toKsub1_Dial(K, wNHost, wNPort, wn.Dial, wn.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("W[%d]: Ready.\n", self)

	W1.Run(func(s *W_1to1_not_2to2and2toKsub1and3toKandKtoK.Init) W_1to1_not_2to2and2toKsub1and3toKandKtoK.End {
		for i := 0; i < NIter*K; i++ {
			rcvd[0].Update = completedLoop(i, K)
			scributil.Delay(1500)
			s0 := s.W_2_Scatter_Particles(rcvd)
			scributil.Debugf("W[%d]: sent %v.\n", self, rcvd)
			s = s0.W_K_Gather_Particles(rcvd) // overwrite previous received particles
			velocities = compute(particles, rcvd[0], velocities)
			scributil.Debugf("[nbody] W[%d]: compute forces.\n", self)
			if completedLoop(i, K) {
				update(particles, velocities)
				scributil.Debugf("[nbody] W[%d]: update positions.\n", self)
			}
		}
		stop := []message.Stop{message.Stop{}}
		scributil.Delay(1500)
		sEnd := s.W_2_Scatter_Stop(stop)
		scributil.Debugf("W[%d]: sent %v.\n", self, stop)
		return *sEnd
	})
	wg.Done()
}

// Wi is the common middle worker W[i].
func Wi(p *NBody.NBody, K, self int, wPrev scributil.ServerConn, wPrevPort int, wNext scributil.ClientConn, wNextHost string, wNextPort int, wg *sync.WaitGroup) {
	// W[i]
	Wi := p.New_family_2_W_2toKsub1and3toK_not_1to1and2to2andKtoK(K, self)
	particles := loadParticles("particles.txt", self)
	rcvd := particles
	var velocities message.Vector3 // Stores the velocities

	scributil.Debugf("[connection] W[%d]: dialling to W[%d] at %s:%d.\n", self, self+1, wNextHost, wNextPort)
	if err := Wi.W_3toKandKtoK_not_1to1and2to2and2toKsub1_Dial(self+1, wNextHost, wNextPort, wNext.Dial, wNext.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	ln, err := wPrev.Listen(wPrevPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	scributil.Debugf("[connection] W[%d]: listening for W[%d] at :%d.\n", self, self-1, wPrevPort)
	if err := Wi.W_2to2and2toKsub1_not_1to1and3toKandKtoK_Accept(self-1, ln, wPrev.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("W[%d]: Ready.\n", self)

	Wi.Run(func(s *W_2toKsub1and3toK_not_1to1and2to2andKtoK.Init) W_2toKsub1and3toK_not_1to1and2to2andKtoK.End {
		for {
			switch s0 := s.W_selfsub1_Branch().(type) {
			case *W_2toKsub1and3toK_not_1to1and2to2andKtoK.Particles_W_Init:
				var p message.Particles
				scributil.Delay(1500)
				s = s0.Recv_Particles(&p).W_selfplus1_Scatter_Particles([]message.Particles{rcvd})
				scributil.Debugf("W[%d]: received %v.\n", self, p)
				rcvd = p
				velocities = compute(particles, p, velocities)
				scributil.Debugf("[nbody] W[%d]: compute forces.\n", self)
				if p.Update {
					update(particles, velocities)
					scributil.Debugf("[nbody] W[%d]: update positions.\n", self)
				}
			case *W_2toKsub1and3toK_not_1to1and2to2andKtoK.Stop_W_Init:
				var stop message.Stop
				scributil.Delay(1500)
				sEnd := s0.Recv_Stop(&stop).W_selfplus1_Scatter_Stop([]message.Stop{stop})
				//		scributil.Debugf("W[%d]: received %v.\n", self, stop)
				return *sEnd
			}
		}
	})
	wg.Done()
}

// WK is the last worker W[K], which loops back to W[1].
func WK(p *NBody.NBody, K, self int, wPrev scributil.ServerConn, wPrevPort int, w1 scributil.ServerConn, w1Port int, wg *sync.WaitGroup) {
	// W[K]
	WK := p.New_family_2_W_3toKandKtoK_not_1to1and2to2and2toKsub1(K, self)
	particles := loadParticles("particles.txt", self)
	rcvd := particles
	var velocities message.Vector3 // Stores the velocities

	lnPrev, err := wPrev.Listen(wPrevPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer lnPrev.Close()
	scributil.Debugf("[connection] W[%d]: listening for W[%d] at :%d.\n", self, self-1, wPrevPort)
	if err := WK.W_2toKsub1and3toK_not_1to1and2to2andKtoK_Accept(self-1, lnPrev, wPrev.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	ln, err := w1.Listen(w1Port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	scributil.Debugf("[connection] W[%d]: listening for W[%d] at :%d.\n", self, 1, w1Port)
	if err := WK.W_1to1_not_2to2and2toKsub1and3toKandKtoK_Accept(1, ln, w1.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("W[%d]: Ready.\n", self)

	WK.Run(func(s *W_3toKandKtoK_not_1to1and2to2and2toKsub1.Init) W_3toKandKtoK_not_1to1and2to2and2toKsub1.End {
		for {
			switch s0 := s.W_selfsub1_Branch().(type) {
			case *W_3toKandKtoK_not_1to1and2to2and2toKsub1.Particles:
				var p message.Particles
				scributil.Delay(1500)
				s = s0.Recv_Particles(&p).W_1_Scatter_Particles([]message.Particles{rcvd})
				scributil.Debugf("W[%d]: received %v.\n", self, p)
				rcvd = p
				velocities = compute(particles, p, velocities)
				scributil.Debugf("[nbody] W[%d]: compute forces.\n", self)
				if p.Update {
					update(particles, velocities)
					scributil.Debugf("[nbody] W[%d]: update positions.\n", self)
				}
			case *W_3toKandKtoK_not_1to1and2to2and2toKsub1.Stop:
				var stop message.Stop
				sEnd := s0.Recv_Stop(&stop)
				scributil.Debugf("W[%d]: received %v.\n", self, stop)
				return *sEnd
			}
		}
	})
	wg.Done()
}

func loadParticles(filename string, offset int) message.Particles {
	return message.Particles{
		Coords: []message.Vector3{
			message.Vector3{
				X: float32(offset),
				Y: float32(offset),
				Z: float32(offset),
			}},
	}
}

// completedLoop returns true if the processes have seen particles
// of all other N-1 processes (i.e. every N iterations).
func completedLoop(i, K int) bool {
	return (i+1)%K == 0
}

var v message.Vector3

// compute calculates and accumulates forces.
func compute(local, received message.Particles, velocities message.Vector3) message.Vector3 {
	// mock calculated vector.
	return v
}

// update updates positions based on forces.
func update(local message.Particles, velocities message.Vector3) {
	// local should be updated using velocities.
}
