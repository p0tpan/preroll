package parser

// Node is a node in the parsed expression tree.
type Node interface {
	node()
}

// ModKind is the single modifier a roll may carry. The grammar allows either
// a keep/drop modifier or an explode, never both.
type ModKind int

const (
	ModNone ModKind = iota
	KeepHighest
	KeepLowest
	DropHighest
	DropLowest
	Explode
)

// Modifier is the optional modifier prefixing a roll.
//
// Keep and drop modifiers always act on exactly one die; the grammar has no
// separate count for them. ExplodeMax is the "!{N}" limit, or 0 when the
// roll used a bare "!".
type Modifier struct {
	Kind       ModKind
	ExplodeMax int
}

// Num is a literal integer operand.
type Num struct {
	Value int
}

// Roll is Count dice. A numeric die has Sides faces numbered 1..Sides; a
// custom-faced die has Faces set and Sides equal to len(Faces).
type Roll struct {
	Count int
	Sides int
	Faces []string
	Mod   Modifier
}

// Op is a binary operator.
type Op rune

const (
	Add Op = '+'
	Sub Op = '-'
)

// BinOp combines two operands. Chained terms are left-associative:
// 1d20-2-3 parses as (1d20-2)-3.
type BinOp struct {
	Op    Op
	Left  Node
	Right Node
}

func (Num) node()   {}
func (Roll) node()  {}
func (BinOp) node() {}
