package bytecode

// The set of available bytecode instructions.
// $0 denotes the 0th (top) item in the stack.
// using $0 pops it from the stack
// [arg] denotes the instruction's argument.
const (
	Nop byte = iota

	/* Storage & constants */
	// LoadConst loads a constant by index: [arg]
	LoadConst

	// LoadName loads a name by index: [arg]
	LoadName

	// StoreName stores $0 in the name indexed by [arg]
	StoreName

	// Declarename is the same as StoreName, but only operates in
	// the single enclosing scope, not the parent ones
	DeclareName

	// LoadSubscript pushes $1[$0]
	LoadSubscript

	// StoreSubscript sets $1[$0] to $2
	StoreSubscript

	/* Operators */
	UnaryInvert
	UnaryNegate
	BinaryAdd
	BinarySub
	BinaryMul
	BinaryDiv
	BinaryExp
	BinaryFloorDiv
	BinaryMod
	BinaryLogicOr
	BinaryLogicAnd
	BinaryBitOr
	BinaryBitAnd
	BinaryEqual
	BinaryNotEqual
	BinaryLess
	BinaryMore
	BinaryLessEq
	BinaryMoreEq

	/* Functions & scopes */
	// CallFunctions calls $0 and pops an item for each argument
	CallFunction

	// CallMethod calls $0 (a method) on $1 (a map), popping an item
	// for each argument
	CallMethod

	Return
	OpenScope
	CloseScope

	/* Control flow */
	// The virtual machine stores jump in a list, allowing a jump
	// argument of 8 bits to jump to a 64-bit code offset

	// Jump jumps to target [arg]
	Jump

	// JumpIf jumps to target [arg] if $0 is truthy
	JumpIf

	// JumpUnless jumps to target [arg] unless $0 is truthy
	JumpUnless

	Break
	Next
	StartLoop
	EndLoop

	/* Data */
	// MakeList pushes a list containing $0, $1, ..., $[arg]
	MakeList

	// MakeTuple pushes a tuple containing $0, $1, ..., $[arg]
	MakeTuple

	// MakeMap pushes a map from the top [arg]*2 items, in key, val order
	MakeMap
)