package shared

type Type int

const (
	TypeInt Type = iota
	TypeFloat
	TypeString
	TypeError
)

func (t Type) String() string {
	switch t {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeString:
		return "string"
	default:
		return "error"
	}
}

type Variable struct {
	Name    string
	Type    Type
	Line    int
	Column  int
	Address int
}

type Function struct {
	Name             string
	Parameters       []Variable
	Line             int
	Column           int
	StartQuad        int
	IntVarsCounter   int
	FloatVarsCounter int
}

type FunctionInfo struct {
	Name           string
	StartQuad      int
	IntVarsCount   int
	FloatVarsCount int
	Parameters     []Variable
}

type Stack struct {
	items []interface{}
}

func NewStack() *Stack {
	return &Stack{
		items: make([]interface{}, 0),
	}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *Stack) Top() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	return s.items[len(s.items)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Size() int {
	return len(s.items)
}

type Quadruple struct {
	Operator string      // The operation to be performed
	LeftOp   interface{} // Left operand
	RightOp  interface{} // Right operand
	Result   interface{} // Where the result will be stored
}
