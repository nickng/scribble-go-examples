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


func (s *Init) A_1_Gather_Int(arg0 []message.Int) *End {
if s.Err != nil {
panic(s.Err)
}
if s.id != s.Ept.lin {
panic("Linear resource already used")
}
for i := 1; i <= 1; i++ {
var tmp session2.ScribMessage
if s.Ept._End.Err = s.Ept.MPChan.MRecv("A", i, &tmp); s.Ept._End.Err != nil {
s.Ept.lin = s.Ept.lin + 1
s.Ept._End.id = s.Ept.lin
return s.Ept._End

}
arg0[i-1] = *(tmp.(*message.Int))
}
s.Ept.lin = s.Ept.lin + 1
s.Ept._End.id = s.Ept.lin
return s.Ept._End

}