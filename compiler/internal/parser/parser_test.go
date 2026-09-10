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
		{
			name:   "function type",
			source: `type test = fn() -> i32;`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimpleFunction{
							Parameters: []ast.SimpleTypeNode{},
							ReturnType: &ast.TypeSimplePrimitive{Name: "i32"},
						},
					},
				},
			},
		},
		{
			name:   "no return function type",
			source: `type test = fn();`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimpleFunction{
							Parameters: []ast.SimpleTypeNode{},
							ReturnType: &ast.TypeSimplePrimitive{Name: "void"},
						},
					},
				},
			},
		},
		{
			name:   "one param function type",
			source: `type test = fn(i32);`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimpleFunction{
							Parameters: []ast.SimpleTypeNode{
								&ast.TypeSimplePrimitive{Name: "i32"},
							},
							ReturnType: &ast.TypeSimplePrimitive{Name: "void"},
						},
					},
				},
			},
		},
		{
			name:   "multiple params function type",
			source: `type test = fn(i32, i4, u1);`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimpleFunction{
							Parameters: []ast.SimpleTypeNode{
								&ast.TypeSimplePrimitive{Name: "i32"},
								&ast.TypeSimplePrimitive{Name: "i4"},
								&ast.TypeSimplePrimitive{Name: "u1"},
							},
							ReturnType: &ast.TypeSimplePrimitive{Name: "void"},
						},
					},
				},
			},
		},
		{
			name:   "generic type",
			source: `type test = Option<i32>;`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimpleGeneric{
							Name: "Option",
							Arguments: []ast.SimpleTypeNode{
								&ast.TypeSimplePrimitive{Name: "i32"},
							},
						},
					},
				},
			},
		},
		{
			name:   "multiple args generic type",
			source: `type test = Result<i32, none>;`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.DeclAliasType{
						Name: "test",
						Base: &ast.TypeSimpleGeneric{
							Name: "Result",
							Arguments: []ast.SimpleTypeNode{
								&ast.TypeSimplePrimitive{Name: "i32"},
								&ast.TypeSimplePrimitive{Name: "none"},
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

func TestComments(t *testing.T) {
	cases := []testCase{
		{
			name:   "line comment",
			source: `// this is a comment`,
			exptProg: &ast.Program{
				Imports:      []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name:   "empty line comment",
			source: `//`,
			exptProg: &ast.Program{
				Imports:      []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name:   "nested line comment",
			source: `// //`,
			exptProg: &ast.Program{
				Imports:      []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name: "multi-line comment",
			source: `/*
this is a comment
*/`,
			exptProg: &ast.Program{
				Imports:      []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name: "nested multi-line comment",
			source: `/*
/*
this is a comment
*/
*/`,
			exptProg: &ast.Program{
				Imports:      []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name:     "not closed multi-line comment",
			source:   `/*`,
			exptProg: nil,
			isError:  true,
		},
		{
			name:     "not closed nested multi-line comment",
			source:   `/* /*`,
			exptProg: nil,
			isError:  true,
		},
		{
			name:     "semi closed nested nested multi-line comment",
			source:   `/* /* */`,
			exptProg: nil,
			isError:  true,
		},
	}

	for _, ca := range cases {
		test(ca, t)
	}
}
