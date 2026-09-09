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
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name:   "global import empty path",
			source: `using @import("");`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportGlobal{Path: ""},
				},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name:   "alias import",
			source: `math = @import("std.math");`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportAlias{Alias: "math", Path: "std.math"},
				},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name:   "alias import empty path",
			source: `math = @import("");`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportAlias{Alias: "math", Path: ""},
				},
				Declarations: []ast.DeclarationNode{},
			},
		},
	}

	for _, ca := range cases {
		test(ca, t)
	}
}

func TestSimpleType(t *testing.T) {
	cases := []testCase{
		{
			name:   "primitive type",
			source: `type test = i32;`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimplePrimitive{Name: "i32"},
					},
				},
			},
		},
		{
			name:   "strange primitive type",
			source: `type test = struct;`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimplePrimitive{Name: "struct"},
					},
				},
			},
		},
		{
			name:   "pointer type",
			source: `type test = *i32;`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimplePointer{Base: &ast.TypeSimplePrimitive{Name: "i32"}},
					},
				},
			},
		},
		{
			name:   "array pointer type",
			source: `type test = [*]i32;`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimpleArrayPointer{Base: &ast.TypeSimplePrimitive{Name: "i32"}},
					},
				},
			},
		},
		{
			name:   "array type",
			source: `type test = [5]i32;`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimpleArray{
							Size: 5,
							Base: &ast.TypeSimplePrimitive{Name: "i32"},
						},
					},
				},
			},
		},
		{
			name:   "slice type",
			source: `type test = []i32;`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimpleSlice{Base: &ast.TypeSimplePrimitive{Name: "i32"}},
					},
				},
			},
		},
		{
			name:   "nested types",
			source: `type test = *[*][5][]i32;`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimplePointer{
							Base: &ast.TypeSimpleArrayPointer{
								Base: &ast.TypeSimpleArray{
									Size: 5,
									Base: &ast.TypeSimpleSlice{
										Base: &ast.TypeSimplePrimitive{Name: "i32"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, ca := range cases {
		test(ca, t)
	}
}
