// Package gql provides AST-like representations of various graphql constructs.
package gql

type Operation struct {
	OperationDefinition
	Variables map[string]interface{}
}

type OperationDefinition struct {
	Type      OperationType
	Name      string
	Selection interface{}
}

type OperationType uint8

const (
	Query OperationType = iota
	Mutation
	Subscription
)
