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

const (
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
	if tree.RootNode().HasError() {
		tree.Close()
		return nil, fmt.Errorf("error while parsing naos source")
	}
	defer tree.Close()

	program := &ast.Program{
		Imports:      []ast.ImportNode{},
		Declarations: []ast.DeclarationNode{},
	}
	root := tree.RootNode()

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

		default:
			return nil, fmt.Errorf("invalid %s node, found at %v", node.Kind(), node.StartPosition())
		}
	}

	return program, nil
}

func naosLang() *tsitter.Language {
	return tsitter.NewLanguage(unsafe.Pointer(C.tree_sitter_naoslang()))
}
