package graphql

import (
	"fmt"

	"github.com/deref/graphql-go/gql"
	"github.com/deref/graphql-go/internal/shurcool"
)

func MustMarshalOperation(op gql.Operation) string {
	var typ string
	switch op.Type {
	case gql.Query:
		typ = "query"
	case gql.Mutation:
		typ = "mutation"
	case gql.Subscription:
		typ = "subscription"
	default:
		panic(fmt.Errorf("unknown operation type: %v", op.Type))
	}
	return shurcool.ConstructOperation(typ, op.Name, op.Selection, op.Variables)
}

// Unmarshal JSON in to the selection structure that uses "graphql:" tags as
// per <https://github.com/shurcooL/graphql>.
type SelectionUnmarshaler struct {
	Selection interface{}
}

func (sel *SelectionUnmarshaler) UnmarshalJSON(bs []byte) error {
	return shurcool.UnmarshalData(bs, sel.Selection)
}
