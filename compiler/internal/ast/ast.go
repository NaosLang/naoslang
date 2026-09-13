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
