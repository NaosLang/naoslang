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

func parseImportNode(node *tsitter.Node, source []byte, isAlias bool) (ast.ImportNode, error) {
	if isAlias {
		alias, err := expectFieldName(
			node, source, "alias",
			"invalid alias in alias import, found at %v", node.StartPosition(),
		)
		if err != nil {
			return nil, err
		}

		path, err := expectFieldName(
			node, source, "path",
			"invalid path in alias import, found at %v", node.StartPosition(),
		)
		if err != nil {
			return nil, err
		}

		return &ast.ImportAliasNode{
			ImportAlias: alias.Utf8Text(source),
			ImportPath:  cleanStr(path, source),
		}, nil
	}

	path, err := expectFieldName(
		node, source, "path",
		"invalid path in global import, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	return &ast.ImportGlobalNode{
		ImportPath: cleanStr(path, source),
	}, nil
}

func cleanStr(node *tsitter.Node, source []byte) string {
	clean := node.Utf8Text(source)
	return clean[1 : len(clean)-1]
}

// +-----------+
// | Type Defs |
// +-----------+

func parseAliasTypeDef(node *tsitter.Node, source []byte) (*ast.TypeDefAliasNode, error) {
	name, err := expectFieldName(
		node, source, "name",
		"invalid name in alias type definition, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	baseNode, err := expectFieldName(
		node, source, "base",
		"invalid base type in alias type definition, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	baseType, err := parsePrimitiveType(baseNode, source)
	if err != nil {
		return nil, err
	}

	return &ast.TypeDefAliasNode{
		TypeName: name.Utf8Text(source),
		BaseType: baseType,
	}, nil
}

// +-------+
// | Types |
// +-------+

func parsePrimitiveType(node *tsitter.Node, source []byte) (ast.TypeNode, error) {
	switch node.Kind() {
	case nodeKindIdType:
		return parseIdType(node, source), nil
	case nodeKindPointerType:
		return parsePointerType(node, source)
	case nodeKindArrayPointerType:
		return parseArrayPointerType(node, source)
	case nodeKindArrayType:
		return parseArrayType(node, source)
	case nodeKindSliceType:
		return parseSliceType(node, source)
	case nodeKindModuleType:
		return parseModuleType(node, source)
	case nodeKindFunctionType:
		return parseFunctionType(node, source)
	case nodeKindConcreteGenericType:
		return parseConcreteGenericType(node, source)
	}
	return nil, fmt.Errorf("invalid primitive type, found at %v", node.StartPosition())
}

func parseIdType(node *tsitter.Node, source []byte) ast.TypeNode {
	name := node.Utf8Text(source)
	if name == "void" {
		return &ast.TypeVoidNode{}
	}
	return &ast.TypeIdNode{TypeName: name}
}

func parsePointerType(node *tsitter.Node, source []byte) (*ast.TypePointerNode, error) {
	constNode := node.ChildByFieldName("const")

	baseNode, err := expectFieldName(
		node, source, "base",
		"invalid pointed type in pointer type, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	baseType, err := parsePrimitiveType(baseNode, source)
	if err != nil {
		return nil, err
	}

	return &ast.TypePointerNode{
		IsPointingConst: constNode != nil,
		PointedType:     baseType,
	}, nil
}

func parseArrayPointerType(node *tsitter.Node, source []byte) (*ast.TypeArrayPointerNode, error) {
	constNode := node.ChildByFieldName("const")

	baseNode, err := expectFieldName(
		node, source, "base",
		"invalid pointed type in array pointer type, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	baseType, err := parsePrimitiveType(baseNode, source)
	if err != nil {
		return nil, err
	}

	return &ast.TypeArrayPointerNode{
		IsDataConst:     constNode != nil,
		PointedDataType: baseType,
	}, nil
}

func parseArrayType(node *tsitter.Node, source []byte) (*ast.TypeArrayNode, error) {
	constNode := node.ChildByFieldName("const")

	sizeNode, err := expectFieldName(
		node, source, "size",
		"invalid size in array type, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	size, err := strconv.ParseInt(sizeNode.Utf8Text(source), 0, 64)
	if err != nil {
		return nil, err
	}

	baseNode, err := expectFieldName(
		node, source, "base",
		"invalid pointed type in array type, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	baseType, err := parsePrimitiveType(baseNode, source)
	if err != nil {
		return nil, err
	}

	return &ast.TypeArrayNode{
		Size:        size,
		IsDataConst: constNode != nil,
		DataType:    baseType,
	}, nil
}

func parseSliceType(node *tsitter.Node, source []byte) (*ast.TypeSliceNode, error) {
	constNode := node.ChildByFieldName("const")

	baseNode, err := expectFieldName(
		node, source, "base",
		"invalid pointed type in slice type, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	baseType, err := parsePrimitiveType(baseNode, source)
	if err != nil {
		return nil, err
	}

	return &ast.TypeSliceNode{
		IsDataConst:     constNode != nil,
		PointedDataType: baseType,
	}, nil
}

func parseModuleType(node *tsitter.Node, source []byte) (*ast.TypeModuleNode, error) {
	module, err := expectFieldName(
		node, source, "module",
		"invalid module name in module type, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	typeNode, err := expectFieldName(
		node, source, "type",
		"invalid type name in module type, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	return &ast.TypeModuleNode{
		ModuleName: module.Utf8Text(source),
		ModuleType: &ast.TypeIdNode{TypeName: typeNode.Utf8Text(source)},
	}, nil
}

func parseFunctionType(node *tsitter.Node, source []byte) (*ast.TypeFunctionNode, error) {
	paramsNode, err := expectFieldName(
		node, source, "parameters",
		"invalid parameters in function type, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	paramTypes := []ast.TypeNode{}
	for i := range paramsNode.NamedChildCount() {
		paramNode := paramsNode.NamedChild(i)
		typeNode, err := expectFieldName(
			paramNode, source, "type",
			"invalid type in function type parameter, found at %v", node.StartPosition(),
		)
		if err != nil {
			return nil, err
		}

		paramType, err := parsePrimitiveType(typeNode, source)
		if err != nil {
			return nil, err
		}

		paramTypes = append(paramTypes, paramType)
	}

	returnNode := node.ChildByFieldName("return")
	if returnNode == nil {
		return &ast.TypeFunctionNode{
			ParameterTypes: paramTypes,
			ReturnType:     &ast.TypeVoidNode{},
		}, nil
	}

	returnType, err := parsePrimitiveType(returnNode, source)
	if err != nil {
		return nil, err
	}

	return &ast.TypeFunctionNode{
		ParameterTypes: paramTypes,
		ReturnType:     returnType,
	}, nil
}

func parseConcreteGenericType(node *tsitter.Node, source []byte) (*ast.TypeConcreteGenericNode, error) {
	name, err := expectFieldName(
		node, source, "name",
		"invalid name in concrete generic type, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	argsNode, err := expectFieldName(
		node, source, "parameters",
		"invalid generic arguments in concrete generic type, found at %v", node.StartPosition(),
	)
	if err != nil {
		return nil, err
	}

	argTypes := []ast.TypeNode{}
	for i := range argsNode.NamedChildCount() {
		argNode := argsNode.NamedChild(i)
		argType, err := parsePrimitiveType(argNode, source)
		if err != nil {
			return nil, err
		}

		argTypes = append(argTypes, argType)
	}

	return &ast.TypeConcreteGenericNode{
		TypeName:             name.Utf8Text(source),
		GenericArgumentTypes: argTypes,
	}, nil
}
