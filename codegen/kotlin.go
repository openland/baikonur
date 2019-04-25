package codegen

import (
	"github.com/openland/baikonur/il"
	"io/ioutil"
	"strconv"
)

func generateReadScalar(alias string, tp string, name string, output *Output) {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	output.WriteLine(scope + ".set(\"" + name + "\", " + scope + ".read" + tp + "(\"" + alias + "\"))")
}

func generateReadOptionalScalar(alias string, tp string, name string, output *Output) {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	output.WriteLine(scope + ".set(\"" + name + "\", " + scope + ".read" + tp + "Optional(\"" + alias + "\"))")
}

func generateReadOptionalListScalar(tp string, output *Output) {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	output.WriteLine(scope + ".next(" + scope + ".read" + tp + "Optional(i))")
}

func generateReadListScalar(tp string, output *Output) {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	output.WriteLine(scope + ".next(" + scope + ".read" + tp + "(i))")
}

func newScope(output *Output, field *il.SelectionField) {
	output.NextScope()
	output.WriteLine("val scope" + strconv.FormatInt(output.GetScope(), 10) + " = scope" + strconv.FormatInt(output.ParentScope(), 10) +
		".child(\"" + field.Alias + "\")")
}

func newScopeInList(output *Output) {
	output.NextScope()
	output.WriteLine("val scope" + strconv.FormatInt(output.GetScope(), 10) + " = scope" + strconv.FormatInt(output.ParentScope(), 10) +
		".child(i)")
}

func newListScope(output *Output, field *il.SelectionField) {
	output.NextScope()
	output.WriteLine("val scope" + strconv.FormatInt(output.GetScope(), 10) + " = scope" + strconv.FormatInt(output.ParentScope(), 10) +
		".childList(\"" + field.Alias + "\")")
}

func newListScopeInList(output *Output) {
	output.NextScope()
	output.WriteLine("val scope" + strconv.FormatInt(output.GetScope(), 10) + " = scope" + strconv.FormatInt(output.ParentScope(), 10) +
		".childList(i)")
}

var nextLevel int64

func generateListNormalizer(level int64, fld *il.SelectionField, list il.List, output *Output) {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	output.WriteLine("for (i in 0 until " + scope + ".size) {")
	output.IndentAdd()
	if list.Inner.GetKind() == "NotNull" {
		inner := list.Inner.(il.NotNull).Inner
		if inner.GetKind() == "Scalar" {
			scalar := inner.(il.Scalar)
			if scalar.Name == "String" {
				generateReadListScalar("String", output)
			} else if scalar.Name == "Int" {
				generateReadListScalar("Int", output)
			} else if scalar.Name == "Float" {
				generateReadListScalar("Float", output)
			} else if scalar.Name == "ID" {
				generateReadListScalar("String", output)
			} else if scalar.Name == "Boolean" {
				generateReadListScalar("Boolean", output)
			} else if scalar.Name == "Date" {
				generateReadListScalar("String", output)
			} else {
				panic("Unsupported scalar: " + scalar.Name)
			}
		} else if inner.GetKind() == "Object" || inner.GetKind() == "Union" || inner.GetKind() == "Interface" {
			output.WriteLine(scope + ".assertObject(i, \"" + fld.Alias + "\")")
			newScopeInList(output)
			if inner.GetKind() == "Object" {
				obj := inner.(il.Object)
				generateNormalizer(obj.SelectionSet, output)
			} else if inner.GetKind() == "Interface" {
				obj := inner.(il.Interface)
				generateNormalizer(obj.SelectionSet, output)
			} else {
				obj := inner.(il.Union)
				generateNormalizer(obj.SelectionSet, output)
			}
			output.ScopePop()
		} else if inner.GetKind() == "Enum" {
			generateReadListScalar("String", output)
		} else {
			panic("Unsupported list inner type " + inner.GetKind())
		}
	} else {
		if list.Inner.GetKind() == "Scalar" {
			scalar := list.Inner.(il.Scalar)
			if scalar.Name == "String" {
				generateReadOptionalListScalar("String", output)
			} else if scalar.Name == "Int" {
				generateReadOptionalListScalar("Int", output)
			} else if scalar.Name == "Float" {
				generateReadOptionalListScalar("Float", output)
			} else if scalar.Name == "ID" {
				generateReadOptionalListScalar("String", output)
			} else if scalar.Name == "Boolean" {
				generateReadOptionalListScalar("Boolean", output)
			} else if scalar.Name == "Date" {
				generateReadOptionalListScalar("String", output)
			} else {
				panic("Unsupported scalar: " + scalar.Name)
			}
		} else if list.Inner.GetKind() == "Object" || list.Inner.GetKind() == "Union" || list.Inner.GetKind() == "Interface" {

		} else if list.Inner.GetKind() == "Enum" {
			generateReadOptionalListScalar("String", output)
		} else if list.Inner.GetKind() == "List" {
			output.WriteLine("if (" + scope + ".isNotNull(i)) {")
			output.IndentAdd()
			newListScopeInList(output)
			generateListNormalizer(nextLevel, fld, list.Inner.(il.List), output)
			output.ScopePop()
			output.IndentRemove()
			output.WriteLine("}")
		} else {
			panic("Unsupported list inner type " + list.Inner.GetKind())
		}
	}
	output.IndentRemove()
	output.WriteLine("}")
	output.WriteLine(scope + ".completed()")
}

func generateNormalizer(set *il.SelectionSet, output *Output) {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	for _, fld := range set.Fields {
		if fld.Type.GetKind() == "NotNull" {
			inner := fld.Type.(il.NotNull).Inner
			if inner.GetKind() == "Scalar" {
				scalar := inner.(il.Scalar)
				if scalar.Name == "String" {
					generateReadScalar(fld.Alias, "String", fld.Name, output)
				} else if scalar.Name == "Int" {
					generateReadScalar(fld.Alias, "Int", fld.Name, output)
				} else if scalar.Name == "Float" {
					generateReadScalar(fld.Alias, "Float", fld.Name, output)
				} else if scalar.Name == "ID" {
					generateReadScalar(fld.Alias, "String", fld.Name, output)
				} else if scalar.Name == "Boolean" {
					generateReadScalar(fld.Alias, "Boolean", fld.Name, output)
				} else if scalar.Name == "Date" {
					generateReadScalar(fld.Alias, "String", fld.Name, output)
				} else {
					panic("Unsupported scalar: " + scalar.Name)
				}
			} else if inner.GetKind() == "Object" || inner.GetKind() == "Union" || inner.GetKind() == "Interface" {
				output.WriteLine(scope + ".assertObject(\"" + fld.Alias + "\")")
				newScope(output, fld)
				if inner.GetKind() == "Object" {
					obj := inner.(il.Object)
					generateNormalizer(obj.SelectionSet, output)
				} else if inner.GetKind() == "Interface" {
					obj := inner.(il.Interface)
					generateNormalizer(obj.SelectionSet, output)
				} else {
					obj := inner.(il.Union)
					generateNormalizer(obj.SelectionSet, output)
				}
				output.ScopePop()
			} else if inner.GetKind() == "List" {
				output.WriteLine("if (" + scope + ".assertList(\"" + fld.Alias + "\")) {")
				output.IndentAdd()
				newListScope(output, fld)
				nextLevel++
				generateListNormalizer(nextLevel, fld, inner.(il.List), output)
				output.ScopePop()
				output.IndentRemove()
				output.WriteLine("}")
			} else if inner.GetKind() == "Enum" {
				generateReadScalar(fld.Alias, "String", fld.Name, output)
			} else {
				panic("Unsupported type: " + inner.GetKind())
			}
			//
		} else {
			if fld.Type.GetKind() == "Scalar" {
				scalar := fld.Type.(il.Scalar)
				if scalar.Name == "String" {
					generateReadOptionalScalar(fld.Alias, "String", fld.Name, output)
				} else if scalar.Name == "Int" {
					generateReadOptionalScalar(fld.Alias, "Int", fld.Name, output)
				} else if scalar.Name == "Float" {
					generateReadOptionalScalar(fld.Alias, "Float", fld.Name, output)
				} else if scalar.Name == "ID" {
					generateReadOptionalScalar(fld.Alias, "String", fld.Name, output)
				} else if scalar.Name == "Boolean" {
					generateReadOptionalScalar(fld.Alias, "Boolean", fld.Name, output)
				} else if scalar.Name == "Date" {
					generateReadOptionalScalar(fld.Alias, "String", fld.Name, output)
				} else {
					panic("Unsupported scalar: " + scalar.Name)
				}
			} else if fld.Type.GetKind() == "Object" || fld.Type.GetKind() == "Union" || fld.Type.GetKind() == "Interface" {
				output.WriteLine("if (" + scope + ".hasKey(\"" + fld.Alias + "\")) {")
				output.IndentAdd()
				newScope(output, fld)
				if fld.Type.GetKind() == "Object" {
					obj := fld.Type.(il.Object)
					generateNormalizer(obj.SelectionSet, output)
				} else if fld.Type.GetKind() == "Interface" {
					obj := fld.Type.(il.Interface)
					generateNormalizer(obj.SelectionSet, output)
				} else {
					obj := fld.Type.(il.Union)
					generateNormalizer(obj.SelectionSet, output)
				}
				output.IndentRemove()
				output.ScopePop()
				output.WriteLine("}")
			} else if fld.Type.GetKind() == "List" {
				output.WriteLine("if (" + scope + ".hasKey(\"" + fld.Alias + "\")) {")
				output.IndentAdd()
				newListScope(output, fld)
				nextLevel++
				generateListNormalizer(nextLevel, fld, fld.Type.(il.List), output)
				output.ScopePop()
				output.IndentRemove()
				output.WriteLine("}")
			} else if fld.Type.GetKind() == "Enum" {
				generateReadOptionalScalar(fld.Alias, "String", fld.Name, output)
			} else {
				panic("Unsupported type: " + fld.Type.GetKind())
			}
		}
	}
	for _, inf := range set.InlineFragments {
		output.WriteLine("if (" + scope + ".isType(\"" + inf.TypeName + "\")) {")
		output.IndentAdd()
		generateNormalizer(inf.Selection, output)
		output.IndentRemove()
		output.WriteLine("}")
	}
	for _, fr := range set.Fragments {
		output.WriteLine("normalize" + fr.Name + "(scope" + strconv.FormatInt(output.GetScope(), 10) + ")")
	}
}

func GenerateKotlin(model *il.Model) {
	output := NewOutput()
	output.WriteLine("package com.openland.soyuz.gen")
	output.WriteLine("")
	for _, f := range model.Fragments {
		output.NextScope()
		output.WriteLine("fun normalize" + f.Name + "(scope: Scope) {")
		output.IndentAdd()
		output.WriteLine("val scope" + strconv.FormatInt(output.GetScope(), 10) + " = scope")
		output.WriteLine("scope" + strconv.FormatInt(output.GetScope(), 10) + ".assertType(\"" + f.TypeName + "\")")
		generateNormalizer(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("}")
	}

	for _, f := range model.Queries {
		output.NextScope()
		output.WriteLine("fun normalize" + f.Name + "(scope: Scope) {")
		output.IndentAdd()
		output.WriteLine("val scope" + strconv.FormatInt(output.GetScope(), 10) + " = scope")
		generateNormalizer(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("}")
	}

	for _, f := range model.Mutations {
		output.NextScope()
		output.WriteLine("fun normalize" + f.Name + "(scope: Scope) {")
		output.IndentAdd()
		output.WriteLine("val scope" + strconv.FormatInt(output.GetScope(), 10) + " = scope")
		generateNormalizer(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("}")
	}

	for _, f := range model.Subscriptions {
		output.NextScope()
		output.WriteLine("fun normalize" + f.Name + "(scope: Scope) {")
		output.IndentAdd()
		output.WriteLine("val scope" + strconv.FormatInt(output.GetScope(), 10) + " = scope")
		generateNormalizer(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("}")
	}

	ioutil.WriteFile("/Users/steve/Develop/soyuz/src/commonMain/kotlin/com.openland.soyuz/gen/Generated.kt",
		[]byte(output.String()), 0644)
}
