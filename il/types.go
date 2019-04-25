package il

type Type interface {
	getKind() string
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
	Name string
	SelectionSet *SelectionSet
}

type Interface struct {
	Name string
	SelectionSet *SelectionSet
}

type Enum struct {
	Name string
}

func (NotNull) getKind() string {
	return "NotNull"
}

func (List) getKind() string {
	return "List"
}

func (Input) getKind() string {
	return "Input"
}

func (Scalar) getKind() string {
	return "Scalar"
}

func (Object) getKind() string {
	return "Object"
}

func (Union) getKind() string {
	return "Union"
}

func (Interface) getKind() string {
	return "Interface"
}

func (Enum) getKind() string {
	return "Enum"
}
