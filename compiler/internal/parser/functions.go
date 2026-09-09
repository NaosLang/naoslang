package parser

import (
	"fmt"

	"github.com/NaosLang/naoslang/internal/ast"
	tsitter "github.com/tree-sitter/go-tree-sitter"
)

// +---------+
// | Imports |
// +---------+

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
