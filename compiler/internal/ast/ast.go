package ast

// +------------+
// | INTERFACES |
// +------------+
type (
	Node interface{ Node() }

	DeclarationNode interface {
		Node
		DeclNode()
	}

	ImportNode interface {
		Node
		ImpNode()
	}
)

// +-------+
// | NODES |
// +-------+

type Program struct {
	Imports      []ImportNode
	Declarations []DeclarationNode
}

// +---------+
// | Imports |
// +---------+

type ImportGlobalNode struct {
	Path string
}

func (n *ImportGlobalNode) Node()    {}
func (n *ImportGlobalNode) ImpNode() {}

type ImportAliasNode struct {
	Alias string
	Path  string
}

func (n *ImportAliasNode) Node()    {}
func (n *ImportAliasNode) ImpNode() {}
