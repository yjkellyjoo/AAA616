package main

import (
    "github.com/llir/llvm/ir"
)

/////////////////////
// Sign Domain //
/////////////////////
type Sign interface {
    String() string
}

type Top struct { }
type Bot struct { }
type Pos struct { }
type Neg struct { }
type Zero struct { }

func (b Bot) String() string { return "Bot" }
func (b Top) String() string { return "Top" }
func (b Pos) String() string { return "Pos" }
func (b Neg) String() string { return "Neg" }
func (b Zero) String() string { return "Zero" }

func SignTop() Sign { return Top{} }
func SignBot() Sign { return Bot{} }

func SignFromInt (l int) Sign {
    if l > 0 {
        return Pos{}
    } else if l < 0 {
        return Neg{}
    } else {
        return Zero{}
    }
}

func SignPlus (i1, i2 Sign) Sign {
    switch i1.(type) {
    case Bot: return Bot{}
    case Top: return Top{}
    case Pos:
        switch i2.(type) {
        case Bot: return Bot{}
        case Pos: return Pos{}
        case Zero: return Pos{}
        case Top: return Top{}
        case Neg: return Top{}
        }
    case Neg:
        switch i2.(type) {
        case Bot: return Bot{}
        case Pos: return Top{}
        case Zero: return Neg{}
        case Top: return Top{}
        case Neg: return Neg{}
        }
    case Zero:
        return i2
    }
    panic("Unreachable")
}

func SignMinus (i1, i2 Sign) Sign {
    switch i1.(type) {
    case Bot: return Bot{}
    case Top: return Top{}
    case Pos:
        switch i2.(type) {
        case Bot: return Bot{}
        case Pos: return Top{}
        case Zero: return Pos{}
        case Top: return Top{}
        case Neg: return Pos{}
        }
    case Neg:
        switch i2.(type) {
        case Bot: return Bot{}
        case Pos: return Top{}
        case Zero: return Neg{}
        case Top: return Top{}
        case Neg: return Top{}
        }
    case Zero:
        switch i2.(type) {
        case Bot: return Bot{}
        case Top: return Top{}
        case Pos: return Neg{}
        case Neg: return Pos{}
        case Zero: return Zero{}
        }
    }
    panic("Unreachable")
}

func SignMult (i1, i2 Sign) Sign {
    switch i1.(type) {
    case Bot: return Bot{}
    case Top: return Top{}
    case Pos:
        switch i2.(type) {
        case Bot: return Bot{}
        case Pos: return Pos{}
        case Zero: return Zero{}
        case Top: return Top{}
        case Neg: return Neg{}
        }
    case Neg:
        switch i2.(type) {
        case Bot: return Bot{}
        case Pos: return Neg{}
        case Zero: return Zero{}
        case Top: return Top{}
        case Neg: return Pos{}
        }
    case Zero:
        switch i2.(type) {
        case Bot: return Bot{}
        default: return Zero{}
        }
    }
    panic("Unreachable")
}

func SignSLE(i1, i2 Sign) Sign {
    switch i1.(type) {
    case Bot: return Bot{}
    case Top: return Top{}
    case Pos:
        switch i2.(type) {
        case Bot: return Bot{}
        case Zero: return Zero{}
        case Neg: return Zero{}
        default: return Top{}
        }
    case Neg:
        switch i2.(type) {
        case Bot: return Bot{}
        case Pos: return Pos{}
        case Zero: return Pos{}
        default: return Top{}
        }
    case Zero:
        switch i2.(type) {
        case Bot: return Bot{}
        case Pos: return Pos{}
        case Zero: return Pos{}
        case Neg: return Zero{}
        default: return Top{}
        }
    }
    panic("Unreachable")
}

func SignOrder(i1, i2 Sign) bool {
    if (i1 == i2 || i2 == Top{} || i1 == Bot{}) {
        return true
    } else {
        return false
    }
}

func SignJoin (i1, i2 Sign) Sign {
    if SignOrder(i1, i2) {
        return i2
    } else if SignOrder(i2, i1) {
        return i1
    } else {
        return Top{}
    }
}

func SignWiden(i1, i2 Sign) Sign {
    return SignJoin(i1, i2)
}

type State map[string]Sign

func EmptyState() State {
    return make(map[string]Sign)
}

func (s *State) Bind(x string, v Sign) {
    (*s)[x] = v
}

func (s *State) Find(x string) Sign {
    v, ok := (*s)[x]
    if !ok {
        return Bot{}
    }
    return v
}

func StateOrder(s1, s2 State) bool {
    for k, v1 := range s1 {
        v2, ok := s2[k]
        if !ok {
            v2 = Bot{}
        }
        if !SignOrder(v1, v2) {
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
            v2 = Bot{}
        }
        s3[k] = SignJoin(v1, v2)
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
            v2 = Bot{}
        }
        s3[k] = SignWiden(v1, v2)
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
