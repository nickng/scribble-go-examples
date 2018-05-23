package B_1toN

import "github.com/nickng/scribble-go-examples/microbenchmarks/message"
import "github.com/rhu1/scribble-go-runtime/runtime/session2"
import "sync/atomic"
import "reflect"
import "sort"

var _ = session2.NewMPChan
var _ = atomic.AddUint64
var _ = reflect.TypeOf
var _ = sort.Sort

type Init struct {
Err error
id uint64
Ept *B_1toN
}


func (s *Init) A_1to1_Gather_Int(arg0 []message.Int) *End {
if s.Err != nil {
panic(s.Err)
}
if s.id != s.Ept.lin {
panic("Linear resource already used")
}
var err error
for i := 1; i <= 1; i++ {
var tmp session2.ScribMessage
if err = s.Ept.MPChan.MRecv("A", i, &tmp); err != nil {
succ := s.Ept._End
s.Ept.lin = s.Ept.lin + 1
succ.id = s.Ept.lin
succ.Err = err
return succ
}
arg0[i-1] = *(tmp.(*message.Int))
}
succ := s.Ept._End
s.Ept.lin = s.Ept.lin + 1
succ.id = s.Ept.lin
succ.Err = err
return succ
}