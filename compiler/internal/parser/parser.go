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

func NewParserError(node *tsitter.Node, message string, source []byte) ParseError {
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

const (
	// +----------+
	// | Comments |
	// +----------+
	nodeKindLineComment      = "comment_line"
	nodeKindMultiLineComment = "comment_multiline"

	// +---------+
	// | Imports |
	// +---------+
	nodeKindGlobalImport = "global_import"
	nodeKindAliasImport  = "alias_import"

	// +--------------+
	// | Simple Types |
	// +--------------+
	nodeKindPrimitiveType    = "id_type"
	nodeKindPointerType      = "ptr_type"
	nodeKindArrayPointerType = "array_ptr_type"
	nodeKindArrayType        = "array_type"
	nodeKindSliceType        = "slice_type"
	nodeKindGenericType      = "generic_id_type"
	nodeKindFunctionType     = "function_type"

	// +-------------------+
	// | Type Declarations |
	// +-------------------+
	nodeKindAliasTypeDecl = "alias_type"
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
		case nodeKindGlobalImport:
			gimp, err := parseImportGlobal(node, source)
			if err != nil {
				return nil, err
			}
			program.Imports = append(program.Imports, gimp)

		case nodeKindAliasImport:
			aimp, err := parseImportAlias(node, source)
			if err != nil {
				return nil, err
			}
			program.Imports = append(program.Imports, aimp)

		case nodeKindAliasTypeDecl:
			atypedecl, err := parseTypedeclAlias(node, source)
			if err != nil {
				return nil, err
			}
			program.Declarations = append(program.Declarations, atypedecl)

		case nodeKindLineComment, nodeKindMultiLineComment:
			continue

		default:
			return nil, NewParserError(node, fmt.Sprintf("invalid %s node, found at", node.Kind()), source)
		}
	}

	return program, nil
}

func collectParseErrors(node *tsitter.Node, source []byte) ParseErrors {
	var errors []ParseError

	if node.IsError() {
		errors = append(errors, NewParserError(node, "syntax error", source))
	}
	if node.IsMissing() {
		errors = append(errors, NewParserError(node, fmt.Sprintf("missing %q", node.Kind()), source))
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
