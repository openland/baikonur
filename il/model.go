package il

// Model

type Model struct {
	Fragments     []*Fragment
	FragmentsMap  map[string]*Fragment
	Queries       []*Operation
	Mutations     []*Operation
	Subscriptions []*Operation
}

func NewModel() *Model {
	return &Model{Fragments: make([]*Fragment, 0), FragmentsMap: make(map[string]*Fragment), Queries: make([]*Operation, 0), Mutations: make([]*Operation, 0), Subscriptions: make([]*Operation, 0)}
}

// Fragments

type Fragment struct {
	Name         string
	TypeName     string
	SelectionSet *SelectionSet
	Uses         []*Fragment
	UsedBy       []*Fragment
}

func NewFragment(name string, typeName string) *Fragment {
	return &Fragment{Name: name, TypeName: typeName, SelectionSet: nil, Uses: make([]*Fragment, 0), UsedBy: make([]*Fragment, 0)}
}

type InlineFragment struct {
	TypeName  string
	Selection *SelectionSet
}

// Selection

type SelectionSet struct {
	Fields          []*SelectionField
	Fragments       []*Fragment
	InlineFragments []*InlineFragment
}

type SelectionField struct {
	Name      string
	Alias     string
	Type      Type
	Selection *SelectionSet
}

// Operation
type Operation struct {
	Type         string
	Name         string
	SelectionSet *SelectionSet
}
