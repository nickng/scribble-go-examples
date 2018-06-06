package A_1to1

import "github.com/nickng/scribble-go-examples/microbenchmarks/message"

type Init struct {
Err error
id uint64
Ept *A_1to1
}




func (s *Init) B_1toN_Scatter_Int(arg0 []message.Int) *End {
if s.Err != nil {
panic(s.Err)
}
if s.id != s.Ept.lin {
panic("Linear resource already used")
}
for i, j := 1, 0; i <= s.Ept.N; i, j = i+1, j+1 {
s.Ept._End.Err = s.Ept.MPChan.MSend("B", i, &arg0[j])
}
s.Ept.lin = s.Ept.lin + 1
s.Ept._End.id = s.Ept.lin
return s.Ept._End

}