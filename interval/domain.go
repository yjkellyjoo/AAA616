package main

import (

    "github.com/llir/llvm/ir"
)

/////////////////////
// Interval Domain //
/////////////////////
type Interval interface {
    String() string
}

type Top struct {}
type Inter struct {
    l int
    u int
}
type Bot struct {}

type Pos struct {}  // struct for +infinite upp
type Neg struct {}  // struct for -infinite low

func (b Bot) String() string { return "Bot" }
func (b Inter) String() string { return "Inter" }
func (b Top) String() string { return "Top" }

func InterTop() Interval { return Top{} }
func InterInter() Interval { return Inter{} }
func InterBot() Interval { return Bot{} }

func InterOrder(i1, i2 Interval) bool {
    if (i1 == i2 || i2 == Top{} || i1 == Bot{}) {
        return true
    } else if (i1.l >= i2.l && i1.u <= i2.u) {
        return true
    } else {
        return false
    }
}

func InterJoin (i1, i2 Interval) Interval {
    if (i1 == Bot{}) {
        return i2
    } else if (i2 == Bot{}) {
        return i1
    } else {
        var inter Inter
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
    if (i1 == Bot{}) {
        return i2
    } else if (i2 == Bot{}) {
        return i1
    } else {
        var inter Inter
        if (i2.l < i1.l) {
            inter.l = Neg{}
        } else {
            inter.l = i1.l
        }
        if (i2.u > i1.u) {
            inter.u = Pos{}
        } else {
            inter.u = i1.u
        }
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
