package main

import (
    "fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
)

func evalArgument(v value.Value, s *State) Sign {
    switch v := v.(type) {
    case *constant.Int:
        i := int(v.X.Int64())
        return SignFromInt(i)
    default: // variable
        loc := v.Ident()
        itv := s.Find(loc)
        return itv
    }
    panic("Unreachable")
}

func (s *State) transferInstAdd(inst *ir.InstAdd) {
    loc := inst.LocalIdent.Ident()  // assign되는 변수를 찾는다
    vx := evalArgument(inst.X, s)   // 덧셈 시키는 첫번째 변수
    vy := evalArgument(inst.Y, s)   // 덧셈 시키는 두번째 변수
    s.Bind(loc, SignPlus(vx, vy))   // 덧셈 action
}

func (s *State) transferInstSub(inst *ir.InstSub) {
    loc := inst.LocalIdent.Ident()
    vx := evalArgument(inst.X, s)
    vy := evalArgument(inst.Y, s)
    s.Bind(loc, SignMinus(vx, vy))
}

func (s *State) transferInstMul(inst *ir.InstMul) {
    loc := inst.LocalIdent.Ident()
    vx := evalArgument(inst.X, s)
    vy := evalArgument(inst.Y, s)
    s.Bind(loc, SignMult(vx, vy))
}

func (s *State) transferInstICmp(inst *ir.InstICmp) {
    loc := inst.LocalIdent.Ident()
    vx := evalArgument(inst.X, s)
    vy := evalArgument(inst.Y, s)
    itv := SignTop()
    switch inst.Pred {
    case enum.IPredSLE:
        itv = SignSLE(vx, vy)
    // 더 다양한 종류의 comp ir 있음
    // case enum.IPredEQ:
    //     ...
    }
    s.Bind(loc, itv)
}

// 양쪽에서 오는 Control Flow 두개 Join 시키는 Phi 함수
func (s *State) transferInstPhi(inst *ir.InstPhi) {
    loc := inst.LocalIdent.Ident()
    itv := SignBot()
    for _, inc := range inst.Incs {
        itv = SignJoin(itv, evalArgument(inc.X, s))
    }
    s.Bind(loc, itv)
}

func (s *State) transferInstCall(inst *ir.InstCall) {
}

// instruction 종류에 따라서 action
func (s *State) transferInst(inst ir.Instruction) {
    switch inst := inst.(type) {
    case *ir.InstAdd: s.transferInstAdd(inst)
    case *ir.InstSub: s.transferInstSub(inst)
    case *ir.InstMul: s.transferInstMul(inst)
    case *ir.InstICmp: s.transferInstICmp(inst)
    case *ir.InstPhi: s.transferInstPhi(inst)
    case *ir.InstCall: s.transferInstCall(inst)
    default: fmt.Printf("Unsupported instruction: %T\n", inst)
    }
}

// basic: State -> State
func (s *State) TransferBlock(insts []ir.Instruction) {
    for _, inst := range insts {
        s.transferInst(inst)    // basic block 안에 있는 instruction들을 compose하는 방식으로 transfer할 수 있음. 
    }
}
