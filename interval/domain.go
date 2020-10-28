package main

import (
    "math"
    "github.com/llir/llvm/ir"
)

/////////////////////
// Interval Domain //
/////////////////////
type Interval interface {
    String() string
}

type Bot struct {}
type Inter struct {
    l int64
    u int64
}
type Top struct {
    l int64
    u int64
}

func (b Bot) String() string { return "Bot" }
func (b Inter) String() string { return "Inter" }
func (b Top) String() string { return "Top" }

func InterBot() Interval { return Bot{} }
// func InterInter() Interval { return Inter{} }
func InterTop() Interval { return Top{ math.MinInt64, math.MaxInt64 } }


/* Integer abstraction function 
*  alpha : int -> sign
*/
func InterFromInt (x int64, y int64) Interval {
    var inter Inter

    // upper interval
    if math.MaxInt64 - x < y {    // overflow -> +infinite
        inter.u = math.MaxInt64
    } else {
        inter.u = x + y
    }
    // lower interval
    if x < y {
        inter.l = x
    } else {
        inter.l = y
    }

    return inter
}

/// sign domain에서의 operation들 정의 ///
func InterPlus (i1 Interval) Interval {
    switch i1.(type) {
    case Bot: return Bot{}  // i1이 Bot일 때
    case Top: return Top{}  // i1이 Top일 때
    // case Pos:   //i1이 +이면서
    //     switch i2.(type) {  // i2가 뭐일 땐?
    //     case Bot: return Bot{}
    //     case Pos: return Pos{}
    //     case Zero: return Pos{}
    //     case Top: return Top{}
    //     case Neg: return Top{}
    //     }
    // case Neg:   //i1이 -이면서
    //     switch i2.(type) {  // i2가 뭐일 땐?
    //     case Bot: return Bot{}
    //     case Pos: return Top{}
    //     case Zero: return Neg{}
    //     case Top: return Top{}
    //     case Neg: return Neg{}
    //     }
    // case Zero:
    //     return i2
    }
    panic("Unreachable")    // sanity check
}

func InterOrder(i1, i2 Interval) bool {
    if (i1 == i2 || i2 == InterTop() || i1 == InterBot()) {
        return true
    } else if (i1 == InterTop() || i2 == InterBot()) {
            return false
    } else {
        i1 := i1.(Inter)
        i2 := i2.(Inter)
        if (i1.l >= i2.l && i1.u <= i2.u) {
            return true
        } else {
            return false
        }
    }
}

func InterJoin (i1, i2 Interval) Interval {
    if (i1 == InterBot()) {
        return i2
    } else if (i2 == InterBot()) {
        return i1
    } else if (i1 == InterTop() || i2 == InterTop()) {
        return i1
    } else {
        var inter Inter
        i1 := i1.(Inter)
        i2 := i2.(Inter)

        if (i1.l < i2.l) {
            inter.l = i1.l
        } else {
            inter.l = i2.l
        }
        if (i1.u > i2.u) {
            inter.u = i1.u
        } else {
            inter.u = i2.u
        }
        return inter
    }
}

// InterBottom?

func InterWiden(i1, i2 Interval) Interval {
    if (i1 == InterBot()) {
        return i2
    } else if (i2 == InterBot()) {
        return i1
    } else {
        var inter Inter
        i1 := i1.(Inter)
        i2 := i2.(Inter)

        if (i2.l < i1.l) {
            inter.l = math.MinInt64
        } else {
            inter.l = i1.l
        }
        if (i2.u > i1.u) {
            inter.u = math.MaxInt64
        } else {
            inter.u = i1.u
        }
        // Inter가 Top이 될 때 확인?
        // if (inter.l == math.MinInt32 && inter.u == math.MaxInt32) {
        //     var top Top = Top(inter)
        //     return top
        // }
        return inter
    }
}

// InterNarrow?

/// end of Interval ///

type State map[string]Interval

func EmptyState() State {
    return make(map[string]Interval)
}

func (s *State) Bind(x string, v Interval) {
    (*s)[x] = v
}

func (s *State) Find(x string) Interval {
    v, ok := (*s)[x]
    if !ok {
        return InterBot()
    }
    return v
}

func StateOrder(s1, s2 State) bool {
    for k, v1 := range s1 {
        v2, ok := s2[k]
        if !ok {
            v2 = InterBot()
        }
        if !InterOrder(v1, v2) {
            return false
        }
    }
    return true
}

func StateJoin(s1, s2 State) State {
    s3 := make(State)
    for k, v := range s2 {
        s3[k] = v
    }
    for k, v1 := range s1 {
        v2, ok := s3[k]
        if !ok {
            v2 = InterBot()
        }
        s3[k] = InterJoin(v1, v2)
    }
    return s3
}

func StateWiden(s1, s2 State) State {
    s3 := make(State)
    for k, v := range s2 {
        s3[k] = v
    }
    for k, v1 := range s1 {
        v2, ok := s3[k]
        if !ok {
            v2 = InterBot()
        }
        s3[k] = InterWiden(v1, v2)
    }
    return s3
}

func (s State) String() string {
    if len(s) == 0 {
        return "{ }"
    }
    var res string
    for k, v := range s {
        res = res + "\t" + k + " |-> " + v.String() + "\n"
    }
    return res
}

type Node *ir.Block
type Table map[Node]State

func (t *Table) Bind (n Node, s State) { (*t)[n] = s }
func (t *Table) Find (n Node) State {
    s, ok := (*t)[n]
    if !ok {
        t.Bind(n, EmptyState())
        return (*t)[n]
    }
    return s
}

func (t *Table) String() string {
    var res string
    for n, s := range *t {
        res += "   " + n.LocalIdent.Ident() + "\n"
        res += s.String() + "\n"
    }
    return res
}

func NewTable() Table {
    return make(Table)
}
