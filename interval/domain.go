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

func InterTop() Interval {
    panic("Not implemented")
}
func InterBot() Interval {
    panic("Not implemented")
}

func InterOrder(i1, i2 Interval) bool {
    panic ("Not implemented")
}

func InterJoin (i1, i2 Interval) Interval {
    panic("Not implemented")
}

func InterWiden(i1, i2 Interval) Interval {
    panic("Not implemented")
}

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
