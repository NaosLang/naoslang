package parser

import (
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
			Alias: alias.Utf8Text(source),
			Path:  cleanStr(path, source),
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
		Path: cleanStr(path, source),
	}, nil
}

func cleanStr(node *tsitter.Node, source []byte) string {
	if node == nil || (node.Kind() != nodeKindStringLiteral && node.Kind() != nodeKindRawStringLiteral) {
		panic("invalid string node passed in cleanStr function")
	}

	clean := node.Utf8Text(source)
	return clean[1 : len(clean)-1]
}
