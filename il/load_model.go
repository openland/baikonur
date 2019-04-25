package il

import (
	"encoding/json"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"io/ioutil"
)

type ClientModel struct {
	Schema        Schema
	Fragments     map[string]*ast.FragmentDefinition
	Subscriptions map[string]*ast.OperationDefinition
	Queries       map[string]*ast.OperationDefinition
	Mutations     map[string]*ast.OperationDefinition
}

type SchemaType struct {
	Kind   string         `json:"kind"`
	Name   string         `json:"name"`
	Fields []*SchemaField `json:"fields"`
	OfType *SchemaType    `json:"ofType"`
}

type SchemaField struct {
	Name string     `json:"name"`
	Type SchemaType `json:"type"`
}

type Schema struct {
	Types []SchemaType `json:"types"`
}

type SchemaRoot struct {
	Schema Schema `json:"__schema"`
}

// Process All Fragment Dependencies

func prepareSelectionSet(selectionSet *ast.SelectionSet, model *Model, clModel *ClientModel) {
	if selectionSet == nil {
		return
	}
	for s := range selectionSet.Selections {
		ss := selectionSet.Selections[s].(ast.Node)
		if ss.GetKind() == "FragmentSpread" {
			fs := ss.(*ast.FragmentSpread)
			prepareFragment(clModel.Fragments[fs.Name.Value], model, clModel)
		} else if ss.GetKind() == "Field" {
			f := ss.(*ast.Field)
			prepareSelectionSet(f.SelectionSet, model, clModel)
		} else if ss.GetKind() == "InlineFragment" {
			fs := ss.(*ast.InlineFragment)
			prepareSelectionSet(fs.SelectionSet, model, clModel)
		} else {
			panic("Unknown selection: " + ss.GetKind())
		}
	}
}

func prepareFragment(fragment *ast.FragmentDefinition, model *Model, clModel *ClientModel) {
	if _, ok := model.FragmentsMap[fragment.Name.Value]; ok {
		return
	}
	fr := NewFragment(fragment.Name.Value, fragment.TypeCondition.Name.Value)
	model.FragmentsMap[fragment.Name.Value] = fr
	prepareSelectionSet(fragment.SelectionSet, model, clModel)
	fr.SelectionSet = convertSelection(fr.TypeName, fragment.SelectionSet, model, clModel)
	deps := collectDependencies(fragment.SelectionSet, model, clModel)
	for i := range deps {
		fr2 := model.FragmentsMap[deps[i]]
		fr.Uses = append(fr.Uses, fr2)
		fr2.UsedBy = append(fr2.UsedBy, fr)
	}

	model.Fragments = append(model.Fragments, fr)
}

// Collect Dependencies

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func collectDependencies(selectionSet *ast.SelectionSet, model *Model, clModel *ClientModel) []string {
	res := make([]string, 0)
	if selectionSet == nil {
		return res
	}
	for s := range selectionSet.Selections {
		ss := selectionSet.Selections[s].(ast.Node)
		if ss.GetKind() == "FragmentSpread" {
			fs := ss.(*ast.FragmentSpread)
			r := collectDependencies(clModel.Fragments[fs.Name.Value].SelectionSet, model, clModel)
			if !contains(res, fs.Name.Value) {
				res = append(res, fs.Name.Value)
			}
			for i := range r {
				if !contains(res, r[i]) {
					res = append(res, r[i])
				}
			}
		} else if ss.GetKind() == "Field" {
			f := ss.(*ast.Field)
			r := collectDependencies(f.SelectionSet, model, clModel)
			for i := range r {
				if !contains(res, r[i]) {
					res = append(res, r[i])
				}
			}
		} else if ss.GetKind() == "InlineFragment" {
			fs := ss.(*ast.InlineFragment)
			r := collectDependencies(fs.SelectionSet, model, clModel)
			for i := range r {
				if !contains(res, r[i]) {
					res = append(res, r[i])
				}
			}
		} else {
			panic("Unknown selection: " + ss.GetKind())
		}
	}
	return res
}

// Convert Selections

func Find(a []SchemaType, name string) *SchemaType {
	for i, n := range a {
		if name == n.Name {
			return &a[i]
		}
	}
	return nil
}

func FindField(a []*SchemaField, name string) *SchemaField {
	for i, n := range a {
		if name == n.Name {
			return a[i]
		}
	}
	return nil
}

func convertType(field *ast.Field, ff SchemaType, model *Model, clModel *ClientModel) Type {
	if ff.Kind == "NON_NULL" {
		// Not null
		return NotNull{convertType(field, *ff.OfType, model, clModel)}
	} else if ff.Kind == "SCALAR" {
		// Scalar
		return Scalar{ff.Name}
	} else if ff.Kind == "OBJECT" {
		// Object
		return Object{ff.Name, convertSelection(ff.Name, field.SelectionSet, model, clModel)}
	} else if ff.Kind == "UNION" {
		// Union
		return Union{ff.Name, convertSelection(ff.Name, field.SelectionSet, model, clModel)}
	} else if ff.Kind == "LIST" {
		// List
		return List{convertType(field, *ff.OfType, model, clModel)}
	} else if ff.Kind == "INTERFACE" {
		// Interface
		return Interface{ff.Name, convertSelection(ff.Name, field.SelectionSet, model, clModel)}
	} else if ff.Kind == "ENUM" {
		// Enum
		return Enum{ff.Name}
	} else {
		panic("Unexpected type kind: " + ff.Kind)
	}
}

func convertSelection(typeName string, selection *ast.SelectionSet, model *Model, clModel *ClientModel) *SelectionSet {
	if selection == nil {
		return nil
	}
	tp := Find(clModel.Schema.Types, typeName)
	fields := make([]*SelectionField, 0)
	fragments := make([]*Fragment, 0)
	inlineFragments := make([]*InlineFragment, 0)
	for s := range selection.Selections {
		ss := selection.Selections[s].(ast.Node)
		if ss.GetKind() == "FragmentSpread" {
			fs := ss.(*ast.FragmentSpread)
			fragments = append(fragments, model.FragmentsMap[fs.Name.Value])
		} else if ss.GetKind() == "Field" {
			f := ss.(*ast.Field)
			alias := f.Name.Value
			if f.Alias != nil {
				alias = f.Alias.Value
			}
			var typ Type
			if f.Name.Value == "__typename" {
				// Special Case
				typ = Scalar{"String"}
			} else {
				ff := FindField(tp.Fields, f.Name.Value)
				typ = convertType(f, ff.Type, model, clModel)
			}
			fld := &SelectionField{Name: f.Name.Value, Alias: alias, Type: typ}
			fields = append(fields, fld)
		} else if ss.GetKind() == "InlineFragment" {
			fs := ss.(*ast.InlineFragment)
			fr := convertSelection(fs.TypeCondition.Name.Value, fs.SelectionSet, model, clModel)
			ifr := &InlineFragment{
				TypeName:  fs.TypeCondition.Name.Value,
				Selection: fr,
			}
			inlineFragments = append(inlineFragments, ifr)
		} else {
			panic("Unknown selection: " + ss.GetKind())
		}
	}
	return &SelectionSet{Fields: fields, Fragments: fragments, InlineFragments: inlineFragments}
}

// Load Model from files

func LoadModel(schemaPath string, files []string) *Model {
	schemaBody, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		panic(err)
	}
	var schemaRoot SchemaRoot
	err = json.Unmarshal(schemaBody, &schemaRoot)
	if err != nil {
		panic(err)
	}
	schema := schemaRoot.Schema
	model := &ClientModel{
		Schema:        schema,
		Fragments:     make(map[string]*ast.FragmentDefinition),
		Queries:       make(map[string]*ast.OperationDefinition),
		Mutations:     make(map[string]*ast.OperationDefinition),
		Subscriptions: make(map[string]*ast.OperationDefinition),
	}
	for i := 0; i < len(files); i++ {
		path := files[i]
		body, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		src := source.NewSource(&source.Source{Body: body})
		parsed, err := parser.Parse(parser.ParseParams{Source: src})
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(parsed.Definitions); i++ {
			node := parsed.Definitions[i]
			if node.GetKind() == "OperationDefinition" {
				op := node.(*ast.OperationDefinition)
				if op.Operation == "query" {
					model.Queries[op.Name.Value] = op
				} else if op.Operation == "mutation" {
					model.Mutations[op.Name.Value] = op
				} else if op.Operation == "subscription" {
					model.Subscriptions[op.Name.Value] = op
				} else {
					panic("Unknown operation: " + op.Operation)
				}
			} else if node.GetKind() == "FragmentDefinition" {
				fr := node.(*ast.FragmentDefinition)
				model.Fragments[fr.Name.Value] = fr
			} else {
				panic("Unknown node: " + node.GetKind())
			}
		}
	}
	ilModel := NewModel()
	for _, v := range model.Fragments {
		prepareFragment(v, ilModel, model)
	}
	return ilModel
}