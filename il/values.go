package il

type Value interface {
	GetKind() string
}

type IntValue struct {
	Int int32
}

type BooleanValue struct {
	Boolean bool
}

type StringValue struct {
	String string
}

type EnumValue struct {
	String string
}

type FloatValue struct {
	Float float64
}

type ListValue struct {
	Values []Value
}

type ObjectValue struct {
	Fields []ObjectValueField
}

type ObjectValueField struct {
	Name  string
	Value Value
}

type VariableValue struct {
	Name string
}

func (IntValue) GetKind() string {
	return "IntValue"
}

func (BooleanValue) GetKind() string {
	return "BooleanValue"
}

func (StringValue) GetKind() string {
	return "StringValue"
}

func (FloatValue) GetKind() string {
	return "FloatValue"
}

func (ListValue) GetKind() string {
	return "ListValue"
}

func (VariableValue) GetKind() string {
	return "VariableValue"
}

func (EnumValue) GetKind() string {
	return "EnumValue"
}

func (ObjectValue) GetKind() string {
	return "ObjectValue"
}
