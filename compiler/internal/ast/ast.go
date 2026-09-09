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

	SimpleTypeNode interface {
		Node
		STypeNode()
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

// ImportAlias -> ID = @import("...");
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

// +-----------+
// | Typedecls |
// +-----------+

// DeclAliasType -> [pub] type ID = SIMPLE_TYPE_NODE
type DeclAliasType struct {
	IsPublic bool
	Name     string
	Base     SimpleTypeNode
}

func (n *DeclAliasType) Node()     {}
func (n *DeclAliasType) DeclNode() {}

// +-------------------+
// | Simple Type Nodes |
// +-------------------+

// TypeSimplePrimitive -> i32
type TypeSimplePrimitive struct {
	Name string
}

func (n *TypeSimplePrimitive) Node()      {}
func (n *TypeSimplePrimitive) STypeNode() {}

// TypeSimplePointer -> *SIMPLE_TYPE_NODE
type TypeSimplePointer struct {
	Base SimpleTypeNode
}

func (n *TypeSimplePointer) Node()      {}
func (n *TypeSimplePointer) STypeNode() {}

// TypeSimpleArrayPointer -> [*]SIMPLE_TYPE_NODE
type TypeSimpleArrayPointer struct {
	Base SimpleTypeNode
}

func (n *TypeSimpleArrayPointer) Node()      {}
func (n *TypeSimpleArrayPointer) STypeNode() {}

// TypeSimpleArray -> [N]SIMPLE_TYPE_NODE
type TypeSimpleArray struct {
	Size int
	Base SimpleTypeNode
}

func (n *TypeSimpleArray) Node()      {}
func (n *TypeSimpleArray) STypeNode() {}

// TypeSimpleSlice -> []SIMPLE_TYPE_NODE
type TypeSimpleSlice struct {
	Base SimpleTypeNode
}

func (n *TypeSimpleSlice) Node()      {}
func (n *TypeSimpleSlice) STypeNode() {}

// TypeSimpleFunction -> fn([ SIMPLE_TYPE_NODE [, ... ] ]) [ -> SIMPLE_TYPE_NODE ]
type TypeSimpleFunction struct {
	Parameters []SimpleTypeNode
	ReturnType SimpleTypeNode
}

func (n *TypeSimpleFunction) Node()      {}
func (n *TypeSimpleFunction) STypeNode() {}

// TypeSimpleGeneric -> ID<SIMPLE_TYPE_NODE [ , ... ]>
type TypeSimpleGeneric struct {
	Name      string
	Arguments []SimpleTypeNode
}

func (n *TypeSimpleGeneric) Node()      {}
func (n *TypeSimpleGeneric) STypeNode() {}
