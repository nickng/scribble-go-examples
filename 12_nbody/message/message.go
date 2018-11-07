package message

import (
	"fmt"
)

// Vector3 is a 3D vector.
type Vector3 struct {
	X, Y, Z float32
}

func (v Vector3) String() string {
	return fmt.Sprintf("(%.1f,%.1f,%.1f)", v.X, v.Y, v.Z)
}

// Particles represents a set of particles
// to be sent to the next neighbour.
type Particles struct {
	Coords []Vector3 // Coords stores the coordinates of vectors
	Update bool      // Update indicates if particles should be updated
}

func (Particles) GetOp() string {
	return "Particles"
}

// Stop represents a message that indicates no more iterations are needed.
type Stop struct {
}

func (Stop) GetOp() string {
	return "Stop"
}
