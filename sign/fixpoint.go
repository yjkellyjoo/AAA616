package main

import (
)

func InputOf(here Node, cfg Cfg, tbl Table) State {
    // 시작 block 이면 
    if cfg.IsEntry(here) {
        res := EmptyState()     // initialize an empty state of map[string] type
        for _, p := range cfg.params {      // CFG parameters 가져와서
            res.Bind(p.Ident(), SignTop())  // p의 정보랑 Top sign 지정
        }
        return res
    } else {    // 시작 block이 아니면
        res := EmptyState()     // initialize an empty state of map[string] type
        for _, p := range cfg.Pred(here) {      //현재 block의 직전 block 정보 가져와서
            res = StateJoin(res, tbl.Find(p))   // 직전 block 정보를 새로 만든 빈거에 저장
        }
        return res
    }
    panic("Unreachable")
}

func NeedWiden(b Node) bool {
    // if b is loophead, then true.
    return true
}

// worklist-based fixpoint algorithm
func Analyze(cfg Cfg) Table {
    var tbl Table
    var worklist Worklist
    tbl = NewTable()    // domain.go에 정의
    worklist = NewWorklist()      // worklist.go에 정의
    worklist.AddSet(cfg.blocks)   // 모든 basic block 추가

    // worklist 순서대로 작업 진행하기
    for !worklist.IsEmpty() {
        here := worklist.Choose()
        state := InputOf(here, cfg, tbl)    // 받은 block의 state 정보를 정리해서 return 해주는 함수
        state.TransferBlock(here.Insts)     // Insts: instruction
        old_state := tbl.Find(here)
        if !StateOrder(state, old_state) {  // fixpoint가 아직 아니면
            if NeedWiden(here) {   // widening이 필요하면
                tbl.Bind(here, StateWiden(old_state, state))
            } else {               // 필요없으면
                tbl.Bind(here, StateJoin(old_state, state))
            }
            worklist.AddSet(cfg.Succ(here)) // 변화에 맞춰서 worklist update
        }
    }
    return tbl
}
