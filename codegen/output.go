package codegen

import "strings"

type Output struct {
	builder     strings.Builder
	indentValue int
	maxScope    int64
	counter     int64
	counters    []int64
}

func NewOutput() *Output {
	return &Output{builder: strings.Builder{}, indentValue: 0, counter: 0, counters: []int64{0}}
}

func (o *Output) IndentAdd() {
	o.indentValue++
}

func (o *Output) NextScope() {
	o.maxScope++
	o.counters = append(o.counters, o.counter)
	o.counter = o.maxScope
}

func (o *Output) ScopePop() {
	nc := o.counters[len(o.counters)-1]
	o.counters = o.counters[:len(o.counters)-1]
	o.counter = nc
}

func (o *Output) GetScope() int64 {
	return o.counter
}

func (o *Output) ParentScope() int64 {
	return o.counters[len(o.counters)-1]
}

func (o *Output) IndentRemove() {
	if o.indentValue == 0 {
		panic("inconsistent ident")
	}
	o.indentValue--
}

func (o *Output) WriteLine(src string) {
	o.builder.WriteString(strings.Repeat(" ", o.indentValue*4) + src + "\n")
}

func (o *Output) String() string {
	return o.builder.String()
}
