package message

// Dimen is the sub-grid and the dimension for each worker.
type Dimen struct {
	Grid          [][]float32 // Width x Height grid
	Width, Height int         // Width and Height of the sub-grid
}

func (Dimen) GetOp() string {
	return "Dimen"
}

// Bound is the boundary (shadow) values.
type Bound struct {
	Bounds []float32
}

func (Bound) GetOp() string {
	return "Bound"
}

// Converged represents a signal where the data has converged
// hence no more calculation needed
type Converged struct {
}

func (Converged) GetOp() string {
	return "Converged"
}

func (Converged) String() string {
	return "Converged"
}
