package encoding

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectionUnmarshaler(t *testing.T) {
	body := `{
		"data": {
			"field": true
		},
		"errors": []
	}`
	type selection struct {
		Field bool `graphql:"field"`
	}
	var q selection
	var res struct {
		Data   interface{}
		Errors []interface{}
	}
	res.Data = &SelectionUnmarshaler{Selection: &q}
	assert.NoError(t, json.Unmarshal([]byte(body), &res))
	var expected selection
	expected.Field = true
	assert.Equal(t, expected, q)
}
