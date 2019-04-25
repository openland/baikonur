package codegen

import "strings"

type Output struct {
	builder     strings.Builder
	indentValue int
}

func NewOutput() *Output {
	return &Output{builder: strings.Builder{}, indentValue: 0}
}

func (o *Output) IndentAdd() {
	o.indentValue++
}

func (o *Output) IndentRemove() {
	if o.indentValue == 0 {
		panic("inconsistent ident")
	}
	o.indentValue--
}

func (o *Output) WriteLine(src string) {
	o.builder.WriteString(strings.Repeat(" ", o.indentValue * 4) + src + "\n")
}

func (o *Output) String() string {
	return o.builder.String()
}
