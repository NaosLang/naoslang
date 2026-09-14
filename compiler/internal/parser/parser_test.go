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
					&ast.ImportGlobalNode{ImportPath: "std.bool"},
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
					&ast.ImportGlobalNode{ImportPath: "std.bool"},
					&ast.ImportGlobalNode{ImportPath: "std.string"},
					&ast.ImportGlobalNode{ImportPath: "std.math"},
				},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name:   "raw string global import",
			source: "using @import(`std.bool`);",
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportGlobalNode{ImportPath: "std.bool"},
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
					&ast.ImportAliasNode{ImportAlias: "bool", ImportPath: "std.bool"},
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
					&ast.ImportAliasNode{ImportAlias: "bool", ImportPath: "std.bool"},
					&ast.ImportAliasNode{ImportAlias: "string", ImportPath: "std.string"},
					&ast.ImportAliasNode{ImportAlias: "math", ImportPath: "std.math"},
				},
				Declarations: []ast.DeclarationNode{},
			},
		},
		{
			name:   "raw string alias import",
			source: "bool = @import(`std.bool`);",
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{
					&ast.ImportAliasNode{ImportAlias: "bool", ImportPath: "std.bool"},
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

func TestPrimitiveTypes(t *testing.T) {
	cases := []testCase{
		{
			name:   "id type",
			source: `type test = i32;`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeIdNode{TypeName: "i32"},
					},
				},
			},
		},
		{
			name:   "void id type",
			source: `type test = void;`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeVoidNode{},
					},
				},
			},
		},
		{
			name: "pointer types",
			source: `
type test = *i32;
type test = *const i32;
`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypePointerNode{
							IsPointingConst: false,
							PointedType:     &ast.TypeIdNode{TypeName: "i32"},
						},
					},
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypePointerNode{
							IsPointingConst: true,
							PointedType:     &ast.TypeIdNode{TypeName: "i32"},
						},
					},
				},
			},
		},
		{
			name: "array pointer types",
			source: `
type test = [*]i32;
type test = [*]const i32;
`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeArrayPointerNode{
							IsDataConst:     false,
							PointedDataType: &ast.TypeIdNode{TypeName: "i32"},
						},
					},
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeArrayPointerNode{
							IsDataConst:     true,
							PointedDataType: &ast.TypeIdNode{TypeName: "i32"},
						},
					},
				},
			},
		},
		{
			name: "array types",
			source: `
type test = [1]i32;
type test = [0xf]const i32;
`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeArrayNode{
							Size:        1,
							IsDataConst: false,
							DataType:    &ast.TypeIdNode{TypeName: "i32"},
						},
					},
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeArrayNode{
							Size:        15,
							IsDataConst: true,
							DataType:    &ast.TypeIdNode{TypeName: "i32"},
						},
					},
				},
			},
		},
		{
			name: "slice types",
			source: `
type test = []i32;
type test = []const i32;
`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeSliceNode{
							IsDataConst:     false,
							PointedDataType: &ast.TypeIdNode{TypeName: "i32"},
						},
					},
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeSliceNode{
							IsDataConst:     true,
							PointedDataType: &ast.TypeIdNode{TypeName: "i32"},
						},
					},
				},
			},
		},
		{
			name: "module type",
			source: `
type test = math.vec3;
`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeModuleNode{
							ModuleName: "math",
							ModuleType: &ast.TypeIdNode{TypeName: "vec3"},
						},
					},
				},
			},
		},
		{
			name: "function types",
			source: `
type test = fn() -> i32;
type test = fn();
type test = fn(i32);
type test = fn(i32, u32, i8) -> void;
`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeFunctionNode{
							ParameterTypes: []ast.TypeNode{},
							ReturnType:     &ast.TypeIdNode{TypeName: "i32"},
						},
					},
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeFunctionNode{
							ParameterTypes: []ast.TypeNode{},
							ReturnType:     &ast.TypeVoidNode{},
						},
					},
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeFunctionNode{
							ParameterTypes: []ast.TypeNode{
								&ast.TypeIdNode{TypeName: "i32"},
							},
							ReturnType: &ast.TypeVoidNode{},
						},
					},
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeFunctionNode{
							ParameterTypes: []ast.TypeNode{
								&ast.TypeIdNode{TypeName: "i32"},
								&ast.TypeIdNode{TypeName: "u32"},
								&ast.TypeIdNode{TypeName: "i8"},
							},
							ReturnType: &ast.TypeVoidNode{},
						},
					},
				},
			},
		},
		{
			name: "concrete generic types",
			source: `
type test = option<i32>;
type test = result<i32, none>;
`,
			exptProg: &ast.Program{
				Imports: []ast.ImportNode{},
				Declarations: []ast.DeclarationNode{
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeConcreteGenericNode{
							TypeName: "option",
							GenericArgumentTypes: []ast.TypeNode{
								&ast.TypeIdNode{TypeName: "i32"},
							},
						},
					},
					&ast.TypeDefAliasNode{
						TypeName: "test",
						BaseType: &ast.TypeConcreteGenericNode{
							TypeName: "result",
							GenericArgumentTypes: []ast.TypeNode{
								&ast.TypeIdNode{TypeName: "i32"},
								&ast.TypeIdNode{TypeName: "none"},
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
