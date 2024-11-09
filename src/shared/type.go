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
	Name       string
	Parameters []Variable
	Line       int
	Column     int
}
