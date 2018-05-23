package Scatter

import A_1to1 "github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Scatter/A_1to1"
import B_1toN "github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Scatter/B_1toN"

type Scatter struct {
}

func (*Scatter) IsProtocol() {
}

func New() *Scatter {
return &Scatter{ }
}

func (p *Scatter) New_A_1to1(N int, self int) *A_1to1.A_1to1 {
return A_1to1.New(p, N, self)
}

func (p *Scatter) New_B_1toN(N int, self int) *B_1toN.B_1toN {
return B_1toN.New(p, N, self)
}
