package ast

type Node interface {
	Node()
}

type Program struct {
	Imports []ImportNode
}

// +---------+
// | Imports |
// +---------+

type ImportNode interface {
	Node
	ImpNode()
}

// ImportAlias -> id = @import("...");
type ImportAlias struct {
	Alias string
	Path  string
}

func (n *ImportAlias) Node()    {}
func (n *ImportAlias) ImpNode() {}

// ImportGlobal -> using @import("...");
type ImportGlobal struct {
	Path string
}

func (n *ImportGlobal) Node()    {}
func (n *ImportGlobal) ImpNode() {}
