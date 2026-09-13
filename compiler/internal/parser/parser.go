package parser

/*
#cgo CFLAGS: -I${SRCDIR}/../tree-sitter/src

#include "tree_sitter/parser.h"

const TSLanguage *tree_sitter_naoslang(void);
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/NaosLang/naoslang/internal/ast"
	tsitter "github.com/tree-sitter/go-tree-sitter"
)

type ParseError struct {
	Message string
	Line    uint
	Column  uint
	Text    string
}

func (e ParseError) Error() string {
	return fmt.Sprintf(
		"%d:%d -> %s (%q)",
		e.Line+1,
		e.Column+1,
		e.Message,
		e.Text,
	)
}

type ParseErrors []ParseError

func (e ParseErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	return e[0].Error()
}

func NewParserError(node *tsitter.Node, source []byte, message string) ParseError {
	if node == nil {
		panic("node can't be nil in 'NewParserError' function")
	}

	pos := node.StartPosition()
	return ParseError{
		Message: message,
		Line:    pos.Row,
		Column:  pos.Column,
		Text:    node.Utf8Text(source),
	}
}

// Tree-sitter node kinds
const (
	// +---------+
	// | Imports |
	// +---------+
	nodeKindGlobalImport = "global_import"
	nodeKindAliasImport  = "alias_import"

	// +----------+
	// | Literals |
	// +----------+
	nodeKindStringLiteral    = "str_literal"
	nodeKindRawStringLiteral = "raw_str_literal"
)

func Parse(source []byte) (*ast.Program, error) {
	parser := tsitter.NewParser()
	parser.SetLanguage(naosLang())
	defer parser.Close()

	tree := parser.Parse(source, nil)
	defer tree.Close()

	root := tree.RootNode()
	if root.HasError() {
		return nil, collectParseErrors(root, source)
	}

	program := &ast.Program{
		Imports:      []ast.ImportNode{},
		Declarations: []ast.DeclarationNode{},
	}

	for i := uint(0); i < root.NamedChildCount(); i++ {
		node := root.NamedChild(i)

		switch node.Kind() {
		case nodeKindGlobalImport, nodeKindAliasImport:
			impNode, err := parseImportNode(node, source, node.Kind() == nodeKindAliasImport)
			if err != nil {
				return nil, err
			}
			program.Imports = append(program.Imports, impNode)

		default:
			return nil, NewParserError(node, source, fmt.Sprintf("invalid %s expression, found at", node.Kind()))
		}
	}

	return program, nil
}

func expectFieldName(node *tsitter.Node, source []byte, fieldName string, errFmt string, args ...any) (*tsitter.Node, error) {
	nodeField := node.ChildByFieldName(fieldName)
	if nodeField == nil {
		return nil, NewParserError(node, source, fmt.Sprintf(errFmt, args...))
	}

	return nodeField, nil
}

func collectParseErrors(node *tsitter.Node, source []byte) ParseErrors {
	var errors []ParseError

	if node.IsError() {
		errors = append(errors, NewParserError(node, source, "syntax error"))
	}
	if node.IsMissing() {
		errors = append(errors, NewParserError(node, source, fmt.Sprintf("missing %q", node.Kind())))
	}

	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)

		errors = append(
			errors,
			collectParseErrors(child, source)...,
		)
	}

	return errors
}

func naosLang() *tsitter.Language {
	return tsitter.NewLanguage(unsafe.Pointer(C.tree_sitter_naoslang()))
}
