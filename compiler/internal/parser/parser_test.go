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

// +---------+
// | Imports |
// +---------+

func TestGlobalImport(t *testing.T) {
	cases := []testCase{
		{
			name:   "one global import",
			source: `using @import("std.bool");`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportGlobalNode{Path: "std.bool"},
				},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name: "multiple global import",
			source: `
using @import("std.bool");
using @import("std.string");
using @import("std.math");
`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportGlobalNode{Path: "std.bool"},
					&ast.ImportGlobalNode{Path: "std.string"},
					&ast.ImportGlobalNode{Path: "std.math"},
				},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name:   "raw string global import",
			source: "using @import(`std.bool`);",
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportGlobalNode{Path: "std.bool"},
				},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name:     "no semicolon at end of global import",
			source:   `using @import("std.bool")`,
			exptProg: nil,
			isError:  true,
		},
	}

	for _, ca := range cases {
		test(ca, t)
	}
}

func TestAliasImport(t *testing.T) {
	cases := []testCase{
		{
			name:   "one alias import",
			source: `bool = @import("std.bool");`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportAliasNode{Alias: "bool", Path: "std.bool"},
				},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name: "multiple alias import",
			source: `
bool = @import("std.bool");
string = @import("std.string");
math = @import("std.math");
`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportAliasNode{Alias: "bool", Path: "std.bool"},
					&ast.ImportAliasNode{Alias: "string", Path: "std.string"},
					&ast.ImportAliasNode{Alias: "math", Path: "std.math"},
				},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name:   "raw string alias import",
			source: "bool = @import(`std.bool`);",
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportAliasNode{Alias: "bool", Path: "std.bool"},
				},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name:     "no semicolon at end of alias import",
			source:   `bool = @import("std.bool")`,
			exptProg: nil,
			isError:  true,
		},
	}

	for _, ca := range cases {
		test(ca, t)
	}
}
