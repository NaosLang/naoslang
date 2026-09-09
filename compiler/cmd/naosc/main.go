package main

import (
	"fmt"
	"regexp"

	"github.com/NaosLang/naoslang/internal/parser"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	source := `using @import("works");`
	prog, err := parser.Parse([]byte(source))
	if err != nil {
		panic(err)
	}

	config := spew.ConfigState{
		Indent:                  "    ",
		DisablePointerAddresses: true,
		DisableCapacities:       true,
	}

	dumpStr := config.Sdump(prog)

	reLen := regexp.MustCompile(`\s*\(len=\d+(?:\s+cap=\d+)?\)`)
	dumpStr = reLen.ReplaceAllString(dumpStr, "")

	rePtr := regexp.MustCompile(`\(0xc[0-9a-fA-F]+\)`)
	dumpStr = rePtr.ReplaceAllString(dumpStr, "")

	fmt.Print(dumpStr)
}
