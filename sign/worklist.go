package main

import (
)

type Worklist []Node

func (w *Worklist) IsEmpty() bool {
    return len(*w) == 0
}

func (w *Worklist) Add(n Node) {
    *w = append(*w, n)
}

// nodes들을 worklist에 통째로 추가
func (w *Worklist) AddSet(nodes []Node) {
    for _, n := range nodes { w.Add(n) }
}

// node 하나 꺼냄
func (w *Worklist) Choose() Node {
    if len(*w) == 0 {
        panic("Worklist is empty")
    }
    top := (*w)[len(*w) - 1]
    (*w) = (*w)[:len(*w) - 1]
    return top
}

func NewWorklist() Worklist {
    return []Node{}
}
