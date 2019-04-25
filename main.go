package main

import (
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func collectFiles(path string) ([]string, error) {
	var res []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".graphql") {
			return nil
		}
		res = append(res, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func main() {
	files, err := collectFiles("tests/queries")
	if err != nil {
		panic(err)
	}
	model := ClientModel{
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
			}
		}
	}
	generate(model)
}
