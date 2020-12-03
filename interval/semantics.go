package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
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
}

func (s *State) transferInstAdd(inst *ir.InstAdd) {
	loc := inst.LocalIdent.Ident() // assign되는 변수를 찾는다
	// fmt.Println("loc: ", loc)
	// fmt.Println("inst.X: ", inst.X)
	// fmt.Println("inst.Y: ", inst.Y)
	// fmt.Println("inst.Typ: ", inst.Typ)
	// fmt.Println("")
	interX := evalArgument(inst.X, s)
	interY := evalArgument(inst.Y, s)
	s.Bind(loc, InterPlus(interX, interY))
}

func (s *State) transferInstSub(inst *ir.InstSub) {
	loc := inst.LocalIdent.Ident() // assign되는 변수를 찾는다
	// fmt.Println("loc: ", loc)
	// fmt.Println("inst.X: ", inst.X)
	// fmt.Println("inst.Y: ", inst.Y)
	// fmt.Println("inst.Typ: ", inst.Typ)
	// fmt.Println("")
	interX := evalArgument(inst.X, s)
	interY := evalArgument(inst.Y, s)
	s.Bind(loc, InterMinus(interX, interY))
}

func (s *State) transferInstMul(inst *ir.InstMul) {
	loc := inst.LocalIdent.Ident()
	vx := evalArgument(inst.X, s)
	vy := evalArgument(inst.Y, s)
	s.Bind(loc, InterMult(vx, vy))
}

func (s *State) transferInstICmp(inst *ir.InstICmp) {
	loc := inst.LocalIdent.Ident()
	vx := evalArgument(inst.X, s)
	vy := evalArgument(inst.Y, s)
	itv := InterTop()
	switch inst.Pred {
	case enum.IPredSLE:
		itv = InterSLE(vx, vy)
		// 더 다양한 종류의 comp ir 있음
		// case enum.IPredEQ:
		//     ...
	}
	s.Bind(loc, itv)
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

func (s *State) transferInstLoad(inst *ir.InstLoad) {
}

func (s *State) transferInstPtrToInt(inst *ir.InstPtrToInt) {
}

func (s *State) transferInstGetElementPtr(inst *ir.InstGetElementPtr) {
}

func (s *State) transferInst(inst ir.Instruction) {
	switch inst := inst.(type) {
	case *ir.InstAdd:
		s.transferInstAdd(inst)
	case *ir.InstSub:
		s.transferInstSub(inst)
	case *ir.InstMul:
		s.transferInstMul(inst)
	case *ir.InstICmp:
		s.transferInstICmp(inst)
	case *ir.InstPhi:
		s.transferInstPhi(inst)
	case *ir.InstCall:
		s.transferInstCall(inst)

	//**********************************
	case *ir.InstLoad:
		s.transferInstLoad(inst)
	case *ir.InstPtrToInt:
		s.transferInstPtrToInt(inst)
	case *ir.InstGetElementPtr:
		s.transferInstGetElementPtr(inst)
	//**********************************

	default:
		fmt.Printf("Unsupported instruction: %T\n", inst)
	}
}

func (s *State) TransferBlock(insts []ir.Instruction) {
	for _, inst := range insts {
		s.transferInst(inst)
	}
}
