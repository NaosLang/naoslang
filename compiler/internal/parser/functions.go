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
		return nil, fmt.Errorf("null path in global import, found at %v", node.StartPosition())
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
		return nil, fmt.Errorf("null identifier in alias import, found at %v", node.StartPosition())
	}

	path := node.ChildByFieldName("path")
	if path == nil {
		return nil, fmt.Errorf("null path in alias import, found at %v", node.StartPosition())
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
	}
	return nil, fmt.Errorf("invalid %s node as simple type, found at %v", node.Kind(), node.StartPosition())
}

// parseTypeSimplePrimitive -> i32
func parseTypeSimplePrimitive(node *tsitter.Node, source []byte) (*ast.TypeSimplePrimitive, error) {
	name := node.NamedChild(0)
	if name == nil {
		return nil, fmt.Errorf("null type name in primitive type, found at %v", node.StartPosition())
	}

	return &ast.TypeSimplePrimitive{
		Name: name.Utf8Text(source),
	}, nil
}

// parseTypeSimplePointer -> *i32
func parseTypeSimplePointer(node *tsitter.Node, source []byte) (*ast.TypeSimplePointer, error) {
	baseNode := node.ChildByFieldName("base")
	if baseNode == nil {
		return nil, fmt.Errorf("null base type in pointer type, found at %v", node.StartPosition())
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
		return nil, fmt.Errorf("null base type in array pointer type, found at %v", node.StartPosition())
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
		return nil, fmt.Errorf("null size in array type, found at %v", node.StartPosition())
	}
	sizeNumb, err := strconv.Atoi(size.Utf8Text(source))
	if err != nil {
		return nil, fmt.Errorf("invalid integer as size in array type, found at %v", node.StartPosition())
	}

	baseNode := node.ChildByFieldName("base")
	if baseNode == nil {
		return nil, fmt.Errorf("null base type in slice type, found at %v", node.StartPosition())
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
		return nil, fmt.Errorf("null base type in slice type, found at %v", node.StartPosition())
	}
	baseType, err := parseTypeSimple(baseNode, source)
	if err != nil {
		return nil, err
	}

	return &ast.TypeSimpleSlice{
		Base: baseType,
	}, nil
}

// +-------------------+
// | Type Declarations |
// +-------------------+

func parseTypedeclAlias(node *tsitter.Node, source []byte) (*ast.DeclAliasType, error) {
	pub := node.ChildByFieldName("visibility")

	name := node.ChildByFieldName("name")
	if name == nil {
		return nil, fmt.Errorf("null alias name in alias typedecl, found at %v", node.StartPosition())
	}

	baseNode := node.ChildByFieldName("base")
	if baseNode == nil {
		return nil, fmt.Errorf("null base type in alias typedecl, found at %v", node.StartPosition())
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
