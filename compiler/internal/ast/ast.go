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

	TypeNode interface {
		Node
		TypeNode()
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
	ImportPath string
}

func (n *ImportGlobalNode) Node()    {}
func (n *ImportGlobalNode) ImpNode() {}

type ImportAliasNode struct {
	ImportAlias string
	ImportPath  string
}

func (n *ImportAliasNode) Node()    {}
func (n *ImportAliasNode) ImpNode() {}

// +-----------+
// | Type Defs |
// +-----------+

type TypeDefAliasNode struct {
	TypeName string
	BaseType TypeNode
}

func (n *TypeDefAliasNode) Node()     {}
func (n *TypeDefAliasNode) DeclNode() {}

// +-------+
// | Types |
// +-------+

type TypeIdNode struct {
	TypeName string
}

func (n *TypeIdNode) Node()     {}
func (n *TypeIdNode) TypeNode() {}

type TypePointerNode struct {
	IsPointingConst bool
	PointedType     TypeNode
}

func (n *TypePointerNode) Node()     {}
func (n *TypePointerNode) TypeNode() {}

type TypeArrayPointerNode struct {
	IsDataConst     bool
	PointedDataType TypeNode
}

func (n *TypeArrayPointerNode) Node()     {}
func (n *TypeArrayPointerNode) TypeNode() {}

type TypeArrayNode struct {
	Size        int64
	IsDataConst bool
	DataType    TypeNode
}

func (n *TypeArrayNode) Node()     {}
func (n *TypeArrayNode) TypeNode() {}

type TypeSliceNode struct {
	IsDataConst     bool
	PointedDataType TypeNode
}

func (n *TypeSliceNode) Node()     {}
func (n *TypeSliceNode) TypeNode() {}

type TypeModuleNode struct {
	ModuleName string
	ModuleType *TypeIdNode
}

func (n *TypeModuleNode) Node()     {}
func (n *TypeModuleNode) TypeNode() {}

type TypeFunctionNode struct {
	ParameterTypes []TypeNode
	ReturnType     TypeNode
}

func (n *TypeFunctionNode) Node()     {}
func (n *TypeFunctionNode) TypeNode() {}

type TypeVoidNode struct{}

func (n *TypeVoidNode) Node()     {}
func (n *TypeVoidNode) TypeNode() {}

type TypeConcreteGenericNode struct {
	TypeName             string
	GenericArgumentTypes []TypeNode
}

func (n *TypeConcreteGenericNode) Node()     {}
func (n *TypeConcreteGenericNode) TypeNode() {}
