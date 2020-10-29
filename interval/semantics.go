package main

import (
 	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/llir/llvm/ir/constant"
//	"github.com/llir/llvm/ir/enum"
)

func evalArgument(v value.Value, s *State) Interval {
    switch v := v.(type) {
    case *constant.Int:
        i := v.X.Int64()
        return InterFromInt(i)
    default: // variable
        loc := v.Ident()
        itv := s.Find(loc)
        return itv
    }
    panic("Unreachable")
}

func (s *State) transferInstAdd(inst *ir.InstAdd) {
    loc := inst.LocalIdent.Ident()  // assign되는 변수를 찾는다
	// fmt.Println("loc: ", loc)
	// fmt.Println("inst.X: ", inst.X)
	// fmt.Println("inst.Y: ", inst.Y)
	// fmt.Println("inst.Typ: ", inst.Typ)
	// fmt.Println("")
    interX := evalArgument(inst.X, s)
    interY := evalArgument(inst.Y, s)
    s.Bind(loc, InterPlus(interX, interY))
}


// 양쪽에서 오는 Control Flow 두개 Join 시키는 Phi 함수
func (s *State) transferInstPhi(inst *ir.InstPhi) {
    loc := inst.LocalIdent.Ident()
    itv := InterBot()
    for _, inc := range inst.Incs {
        itv = InterJoin(itv, evalArgument(inc.X, s))
    }
    s.Bind(loc, itv)
}

func (s *State) transferInstCall(inst *ir.InstCall) {
}


func (s *State) transferInst(inst ir.Instruction) {
    switch inst := inst.(type) {
    case *ir.InstAdd: s.transferInstAdd(inst)
    // case *ir.InstMul: s.transferInstMul(inst)
    // case *ir.InstICmp: s.transferInstICmp(inst)
    case *ir.InstPhi: s.transferInstPhi(inst)
    case *ir.InstCall: s.transferInstCall(inst)
    default: fmt.Printf("Unsupported instruction: %T\n", inst)
    }
}

func (s *State) TransferBlock(insts []ir.Instruction) {
    for _, inst := range insts {
        s.transferInst(inst)
    }
}
