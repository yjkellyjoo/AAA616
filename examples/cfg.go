package main

import (
    "flag"
	"fmt"
	"strings"

	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
)

type Cfg struct {
    fid string
    blocks []*ir.Block
}

// example usage: go run cfg.go ../test/gzip-1.2.4a.ll @shorten_name
func main() {
    flag.Parse()
    args := flag.Args()
    filename := args[0]
    funcname := args[1]
	m, err := asm.ParseFile(filename)
	if err != nil {
		panic(err)
	}
	cfg := genCFG(m, funcname)
	fmt.Println(cfg)
}

func genCFG(m *ir.Module, funcname string) string {
	buf := &strings.Builder{}
	buf.WriteString("digraph {\n")
	for _, f := range m.Funcs {
        if f.Ident() == funcname {
            for _, block := range f.Blocks {
                term := block.Term
                for _, succ := range term.Succs() {
                    fmt.Fprintf(buf, "\t%q -> %q\n", block, succ)
                }
            }
        }
	}
	buf.WriteString("}")
	return buf.String()
}
