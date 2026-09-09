package parser

import (
	"testing"

	"github.com/NaosLang/naoslang/internal/ast"
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
	if err != nil {
		if tcase.isError {
			return
		}
		t.Fatalf("[%s]: expected 0 parsing error, received 1 -> %v", tcase.name, err)
	} else if tcase.isError {
		t.Fatalf("[%s]: expected 1 parsing error, received 0", tcase.name)
	}

	if diff := cmp.Diff(tcase.exptProg, prog); diff != "" {
		t.Fatalf("[%s]: program mismatch (-expected, +got)\n%s", tcase.name, diff)
	}
}

func TestImports(t *testing.T) {
	cases := []testCase{
		{
			name:   "global import",
			source: `using @import("std.bool");`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportGlobal{Path: "std.bool"},
				},
			},
		},
		{
			name:   "global import empty path",
			source: `using @import("");`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportGlobal{Path: ""},
				},
			},
		},
		{
			name:   "alias import",
			source: `math = @import("std.math");`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportAlias{Alias: "math", Path: "std.math"},
				},
			},
		},
		{
			name:   "alias import empty path",
			source: `math = @import("");`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportAlias{Alias: "math", Path: ""},
				},
			},
		},
	}

	for _, ca := range cases {
		test(ca, t)
	}
}
