package parser

import (
	"fmt"
	"strconv"

	"github.com/NaosLang/naoslang/internal/ast"
	tsitter "github.com/tree-sitter/go-tree-sitter"
)

// +---------+
// | Imports |
// +---------+

// parseImportGlobal -> using @import("...");
func parseImportGlobal(node *tsitter.Node, source []byte) (*ast.ImportGlobal, error) {
	path := node.ChildByFieldName("path")
	if path == nil {
		return nil, NewParserError(node, "null path in global import", source)
	}
	pathStr := path.Utf8Text(source)

	return &ast.ImportGlobal{
		Path: pathStr[1 : len(pathStr)-1],
	}, nil
}

// parseImportAlias -> ID = @import("...");
func parseImportAlias(node *tsitter.Node, source []byte) (*ast.ImportAlias, error) {
	alias := node.NamedChild(0)
	if alias == nil {
		return nil, NewParserError(node, "null identifier in alias import", source)
	}

	path := node.ChildByFieldName("path")
	if path == nil {
		return nil, NewParserError(node, "null path in alias import", source)
	}
	pathStr := path.Utf8Text(source)

	return &ast.ImportAlias{
		Alias: alias.Utf8Text(source),
		Path:  pathStr[1 : len(pathStr)-1],
	}, nil
}

// +--------------+
// | Simple Types |
// +--------------+

// parseTypeSimple -> i32 | *i32 | [*]i32 | [5]i32 | []i32
func parseTypeSimple(node *tsitter.Node, source []byte) (ast.SimpleTypeNode, error) {
	switch node.Kind() {
	case nodeKindPrimitiveType:
		return parseTypeSimplePrimitive(node, source)
	case nodeKindPointerType:
		return parseTypeSimplePointer(node, source)
	case nodeKindArrayPointerType:
		return parseTypeSimpleArrayPointer(node, source)
	case nodeKindArrayType:
		return parseTypeSimpleArray(node, source)
	case nodeKindSliceType:
		return parseTypeSimpleSlice(node, source)
	case nodeKindFunctionType:
		return parseTypeSimpleFunction(node, source)
	case nodeKindGenericType:
		return parseTypeSimpleGeneric(node, source)
	}
	return nil, NewParserError(node, fmt.Sprintf("invalid %s node as simple type", node.Kind()), source)
}

// parseTypeSimplePrimitive -> i32
func parseTypeSimplePrimitive(node *tsitter.Node, source []byte) (*ast.TypeSimplePrimitive, error) {
	name := node.NamedChild(0)
	if name == nil {
		return nil, NewParserError(node, "null type name in primitive type", source)
	}

	return &ast.TypeSimplePrimitive{
		Name: name.Utf8Text(source),
	}, nil
}

// parseTypeSimplePointer -> *i32
func parseTypeSimplePointer(node *tsitter.Node, source []byte) (*ast.TypeSimplePointer, error) {
	baseNode := node.ChildByFieldName("base")
	if baseNode == nil {
		return nil, NewParserError(node, "null base type in pointer type", source)
	}
	baseType, err := parseTypeSimple(baseNode, source)
	if err != nil {
		return nil, err
	}

	return &ast.TypeSimplePointer{
		Base: baseType,
	}, nil
}

// parseTypeSimpleArrayPointer -> [*]i32
func parseTypeSimpleArrayPointer(node *tsitter.Node, source []byte) (*ast.TypeSimpleArrayPointer, error) {
	baseNode := node.ChildByFieldName("base")
	if baseNode == nil {
		return nil, NewParserError(node, "null base type in array pointer type", source)
	}
	baseType, err := parseTypeSimple(baseNode, source)
	if err != nil {
		return nil, err
	}

	return &ast.TypeSimpleArrayPointer{
		Base: baseType,
	}, nil
}

// parseTypeSimpleArray -> [5]i32
func parseTypeSimpleArray(node *tsitter.Node, source []byte) (*ast.TypeSimpleArray, error) {
	size := node.ChildByFieldName("size")
	if size == nil {
		return nil, NewParserError(node, "null size in array type", source)
	}
	sizeNumb, err := strconv.Atoi(size.Utf8Text(source))
	if err != nil {
		return nil, NewParserError(node, "invalid integer as size in array type", source)
	}

	baseNode := node.ChildByFieldName("base")
	if baseNode == nil {
		return nil, NewParserError(node, "null base type in slice type", source)
	}
	baseType, err := parseTypeSimple(baseNode, source)
	if err != nil {
		return nil, err
	}

	return &ast.TypeSimpleArray{
		Size: sizeNumb,
		Base: baseType,
	}, nil
}

// parseTypeSimpleSlice -> []i32
func parseTypeSimpleSlice(node *tsitter.Node, source []byte) (*ast.TypeSimpleSlice, error) {
	baseNode := node.ChildByFieldName("base")
	if baseNode == nil {
		return nil, NewParserError(node, "null base type in slice type", source)
	}
	baseType, err := parseTypeSimple(baseNode, source)
	if err != nil {
		return nil, err
	}

	return &ast.TypeSimpleSlice{
		Base: baseType,
	}, nil
}

func parseTypeSimpleFunction(node *tsitter.Node, source []byte) (*ast.TypeSimpleFunction, error) {
	params := node.NamedChild(0)
	if params == nil {
		return nil, NewParserError(node, "null parameters node in function type", source)
	}

	fn := &ast.TypeSimpleFunction{
		Parameters: []ast.SimpleTypeNode{},
		ReturnType: &ast.TypeSimplePrimitive{Name: "void"},
	}

	for i := uint(0); i < params.NamedChildCount(); i++ {
		param := params.NamedChild(i)
		if param == nil {
			return nil, NewParserError(node, "null parameter node, in function parameters", source)
		}

		paramType, err := parseTypeSimple(param, source)
		if err != nil {
			return nil, err
		}

		fn.Parameters = append(fn.Parameters, paramType)
	}

	returnNode := node.ChildByFieldName("return_type")
	if returnNode != nil {
		returnType, err := parseTypeSimple(returnNode, source)
		if err != nil {
			return nil, err
		}

		fn.ReturnType = returnType
	}

	return fn, nil
}

func parseTypeSimpleGeneric(node *tsitter.Node, source []byte) (*ast.TypeSimpleGeneric, error) {
	name := node.ChildByFieldName("name")
	if name == nil {
		return nil, NewParserError(node, "null type name in generic type", source)
	}

	args := node.ChildByFieldName("arguments")
	if args == nil {
		return nil, NewParserError(node, "null arguments node, in generic type", source)
	}

	generic := &ast.TypeSimpleGeneric{
		Name:      name.Utf8Text(source),
		Arguments: []ast.SimpleTypeNode{},
	}

	for i := uint(0); i < args.NamedChildCount(); i++ {
		arg := args.NamedChild(i)
		if arg == nil {
			return nil, NewParserError(node, "null argument node, in generic arguments", source)
		}

		argType, err := parseTypeSimple(arg, source)
		if err != nil {
			return nil, err
		}

		generic.Arguments = append(generic.Arguments, argType)
	}
	return generic, nil
}

// +-------------------+
// | Type Declarations |
// +-------------------+

func parseTypedeclAlias(node *tsitter.Node, source []byte) (*ast.DeclAliasType, error) {
	pub := node.ChildByFieldName("visibility")

	name := node.ChildByFieldName("name")
	if name == nil {
		return nil, NewParserError(node, "null alias name in alias typedecl", source)
	}

	baseNode := node.ChildByFieldName("base")
	if baseNode == nil {
		return nil, NewParserError(node, "null base type in alias typedecl", source)
	}
	baseType, err := parseTypeSimple(baseNode, source)
	if err != nil {
		return nil, err
	}

	return &ast.DeclAliasType{
		IsPublic: pub != nil,
		Name:     name.Utf8Text(source),
		Base:     baseType,
	}, nil
}
