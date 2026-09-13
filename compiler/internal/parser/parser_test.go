package parser

import (
	"testing"

	"github.com/NaosLang/naoslang/internal/ast"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
)

type testCase struct {
	name     string
	source   string
	exptProg *ast.Program
	isError  bool
}

func test(tcase testCase, t *testing.T) {
	t.Helper()

	prog, err := Parse([]byte(tcase.source))
	spewCfg := spew.ConfigState{
		Indent:                  "    ",
		DisablePointerAddresses: true,
		DisableCapacities:       true,
	}

	if err != nil {
		if tcase.isError {
			return
		}
		spewCfg.Dump(tcase)
		t.Fatalf("[%s]: expected 0 parsing error, received 1 -> %v", tcase.name, err)
	} else if tcase.isError {
		spewCfg.Dump(tcase)
		t.Fatalf("[%s]: expected 1 parsing error, received 0", tcase.name)
	}

	if diff := cmp.Diff(tcase.exptProg, prog); diff != "" {
		spewCfg.Dump(tcase)
		t.Fatalf("[%s]: program mismatch (-expected, +got)\n%s", tcase.name, diff)
	}
}
