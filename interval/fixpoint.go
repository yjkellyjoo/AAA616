package main

import (
	"math"
)

func InputOf(here Node, cfg Cfg, tbl Table) State {
	if cfg.IsEntry(here) {
		res := EmptyState()
		for _, p := range cfg.params {
			res.Bind(p.Ident(), InterTop())
		}
		return res
	} else {
		res := EmptyState()
		for _, p := range cfg.Pred(here) {
			res = StateJoin(res, tbl.Find(p))
		}
		return res
	}
}

//***************************
func NeedWiden(count int) bool {
	// index번 까지는 StateJoin으로 해보고 안 끝나겠다 싶으면 widening
	if count < index {
		return false
	} else {
		return true
	}
}

func Analyze(cfg Cfg) Table {
	var tbl Table
	var worklist Worklist
	tbl = NewTable()
	worklist = NewWorklist()
	worklist.AddSet(cfg.blocks)

	//***************************
	count := 0
	var wide_threshold [index]int64
	for i := 0; i < index/2; i++ {
		c := math.Pow(2, float64(i))
		wide_threshold[index/2+i] = int64(c)
		wide_threshold[index/2-i] = int64(-c)
	}
	//***************************

	for !worklist.IsEmpty() {
		here := worklist.Choose()
		state := InputOf(here, cfg, tbl)
		state.TransferBlock(here.Insts)
		old_state := tbl.Find(here)

		if !StateOrder(state, old_state) {
			if NeedWiden(count) {
				tbl.Bind(here, StateWiden(old_state, state, wide_threshold))
				// TODO: Narrowing after widening
				// must fix - doing widening and narrowing all the time may cause an infinite loop
				// tbl.Bind(here, StateNarrow(old_state, state))
			} else {
				tbl.Bind(here, StateJoin(old_state, state))
				count++
			}
			worklist.AddSet(cfg.Succ(here))
		}
	}
	return tbl
}
