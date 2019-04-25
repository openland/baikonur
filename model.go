package main

import "github.com/graphql-go/graphql/language/ast"

type ClientModel struct {
	Fragments map[string]*ast.FragmentDefinition
	Subscriptions map[string]*ast.OperationDefinition
	Queries map[string]*ast.OperationDefinition
	Mutations map[string]*ast.OperationDefinition
}