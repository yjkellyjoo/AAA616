package main

import (
    "flag"
    "fmt"

	"github.com/llir/llvm/asm"
)

func main() {
    flag.Parse()
    args := flag.Args()
    filename := args[0]
	m, err := asm.ParseFile(filename)  //LLVM ASM 결과 받아서
	if err != nil {
		panic(err)
	}
    md := NewModule(m)
    for _, cfg := range md.cfgs {   // 그 안에 CFG 활용,
        fmt.Printf("Analysis of %s begins...\n", cfg.GetFid())
        tbl := Analyze(cfg)     // fixpoint.go에 Analyze func으로 분석 결과
        fmt.Println()
        fmt.Println("Analysis Results for:", cfg.GetFid())
        fmt.Println(tbl.String())
        fmt.Println()
    }
}
