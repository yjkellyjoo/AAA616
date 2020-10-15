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
    loc := inst.LocalIdent.Ident()
    vx := evalArgument(inst.X, s)
    vy := evalArgument(inst.Y, s)
    s.Bind(loc, SignPlus(vx, vy))
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
    }
    s.Bind(loc, itv)
}

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

func (s *State) TransferBlock(insts []ir.Instruction) {
    for _, inst := range insts {
        s.transferInst(inst)
    }
}
