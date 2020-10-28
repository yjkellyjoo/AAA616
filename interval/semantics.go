package main

import (
 	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/llir/llvm/ir/constant"
//	"github.com/llir/llvm/ir/enum"
)

func evalArgument(x value.Value, y value.Value, s *State) Interval {
    switch x := x.(type) {
    case *constant.Int:
    	switch y := y.(type) {
    	case *constant.Int:
	        i := x.X.Int64()
	        j := y.X.Int64()
	        return InterFromInt(i, j)
    	}
    default: // variable
        loc := x.Ident()
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
    int := evalArgument(inst.X, inst.Y, s)   // Sign과 다르게 두 매개변수가 동시에 필요. 
    s.Bind(loc, InterPlus(int))

}

func (s *State) transferInstCall(inst *ir.InstCall) {
}


func (s *State) transferInst(inst ir.Instruction) {
    switch inst := inst.(type) {
    case *ir.InstAdd: s.transferInstAdd(inst)
    // case *ir.InstMul: s.transferInstMul(inst)
    // case *ir.InstICmp: s.transferInstICmp(inst)
    // case *ir.InstPhi: s.transferInstPhi(inst)
    // case *ir.InstCall: s.transferInstCall(inst)
    default: fmt.Printf("Unsupported instruction: %T\n", inst)
    }
}

func (s *State) TransferBlock(insts []ir.Instruction) {
    for _, inst := range insts {
        s.transferInst(inst)
    }
}
