package il

type Value interface {
	getKind() string
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

func (IntValue) getKind() string {
	return "IntValue"
}

func (BooleanValue) getKind() string {
	return "BooleanValue"
}

func (StringValue) getKind() string {
	return "StringValue"
}

func (FloatValue) getKind() string {
	return "FloatValue"
}

func (ListValue) getKind() string {
	return "ListValue"
}

func (VariableValue) getKind() string {
	return "VariableValue"
}

func (EnumValue) getKind() string {
	return "EnumValue"
}

func (ObjectValue) getKind() string {
	return "ObjectValue"
}
