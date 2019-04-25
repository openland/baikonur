package il

type Type interface {
	GetKind() string
}

type NotNull struct {
	Inner Type
}

type List struct {
	Inner Type
}

type Input struct {
	Name string
}

type Scalar struct {
	Name string
}

type Object struct {
	Name         string
	SelectionSet *SelectionSet
}

type Union struct {
	Name         string
	SelectionSet *SelectionSet
}

type Interface struct {
	Name         string
	SelectionSet *SelectionSet
}

type Enum struct {
	Name string
}

func (NotNull) GetKind() string {
	return "NotNull"
}

func (List) GetKind() string {
	return "List"
}

func (Input) GetKind() string {
	return "Input"
}

func (Scalar) GetKind() string {
	return "Scalar"
}

func (Object) GetKind() string {
	return "Object"
}

func (Union) GetKind() string {
	return "Union"
}

func (Interface) GetKind() string {
	return "Interface"
}

func (Enum) GetKind() string {
	return "Enum"
}
