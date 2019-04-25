package main

import (
	"github.com/graphql-go/graphql/language/ast"
	"github.com/openland/baikonur/codegen"
	"sort"
	"strconv"
	"strings"
)

func generateSelectionSetNormalizer(model *ClientModel, s *ast.SelectionSet, output *codegen.Output) {
	for i := 0; i < len(s.Selections); i++ {
		s := s.Selections[i].(ast.Node)
		if s.GetKind() == "Field" {
			f := s.(*ast.Field)
			responseName := f.Name.Value
			if f.Alias != nil {
				responseName = f.Alias.Value
			}
			storeName := f.Name.Value
			if f.Arguments != nil && len(f.Arguments) > 0 {
				argkeys := make([]string, 0)
				for j := 0; j < len(f.Arguments); j++ {
					arg := f.Arguments[j]
					var argumentKey string
					if arg.Value.GetKind() == "IntValue" {
						argumentKey = arg.Value.(*ast.IntValue).Value
					} else if arg.Value.GetKind() == "FloatValue" {
						argumentKey = arg.Value.(*ast.FloatValue).Value
					} else if arg.Value.GetKind() == "StringValue" {
						argumentKey = "\"" + arg.Value.(*ast.StringValue).Value + "\""
					} else if arg.Value.GetKind() == "BooleanValue" {
						argumentKey = strconv.FormatBool(arg.Value.(*ast.BooleanValue).Value)
					} else if arg.Value.GetKind() == "EnumValue" {
						argumentKey = arg.Value.(*ast.EnumValue).Value
					} else if arg.Value.GetKind() == "ListValue" {
						panic("List Value is not supported yet")
					} else if arg.Value.GetKind() == "ObjectValue" {
						panic("Object Value is not supported yet")
					} else if arg.Value.GetKind() == "Variable" {
						v := arg.Value.(*ast.Variable)
						argumentKey = "scope.argumentKey(\"" + v.Name.Value + "\")"
					} else {
						panic("Unknown value kind: " + arg.Value.GetKind())
					}
					argkeys = append(argkeys, "\""+arg.Name.Value+":\"+"+argumentKey)
				}
				sort.Strings(argkeys)
				storeName = "\"" + f.Name.Value + "(\"+" + strings.Join(argkeys, "+ \",\" + ") + "+\")\""
			}
			if f.SelectionSet != nil {
				output.WriteLine("if (scope.push(\"" + responseName + "\", " + storeName + ")) {")
				output.IndentAdd()
				generateSelectionSetNormalizer(model, f.SelectionSet, output)
				output.WriteLine("scope.pop()")
				output.IndentRemove()
				output.WriteLine("}")
			} else {
				output.WriteLine("scope.output(" + storeName + ", scope.input(\"" + responseName + "\"))")
			}
		} else if s.GetKind() == "FragmentSpread" {
			f := s.(*ast.FragmentSpread)
			fr := model.Fragments[f.Name.Value]
			output.WriteLine("if (scope.typename == \"" + fr.TypeCondition.Name.Value + "\") {")
			output.IndentAdd()
			if fr.SelectionSet != nil {
				generateSelectionSetNormalizer(model, fr.SelectionSet, output)
			}
			output.IndentRemove()
			output.WriteLine("}")
		} else if s.GetKind() == "InlineFragment" {
			f := s.(*ast.InlineFragment)
			output.WriteLine("if (scope.typename == \"" + f.TypeCondition.Name.Value + "\") {")
			output.IndentAdd()
			if f.SelectionSet != nil {
				generateSelectionSetNormalizer(model, f.SelectionSet, output)
			}
			output.IndentRemove()
			output.WriteLine("}")
		} else {
			panic("Unknown selection type: " + s.GetKind())
		}
	}
}

func generateNormalizer(model *ClientModel, op *ast.OperationDefinition, output *codegen.Output) {
	output.WriteLine("fun normalize" + op.Name.Value + "(src: JsonObject) {")
	output.IndentAdd()
	output.WriteLine("var scope = Scope(src)")
	generateSelectionSetNormalizer(model, op.SelectionSet, output)
	output.IndentRemove()
	output.WriteLine("}")
}

func generate(model *ClientModel) {
	output := codegen.NewOutput()
	for _, v := range model.Queries {
		generateNormalizer(model, v, output)
	}
	print(output.String())
}
