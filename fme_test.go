package fme

import (
	"testing"
)


func TestSchema(t *testing.T) {
	schema := NewSchema()

	schema.Constraints[0].Add("a", "b", schema.Graph)
	schema.Constraints[0].Add("b", "c", schema.Graph)
	schema.Constraints[1].Add("c", "a", schema.Graph)


	if ok,  err := schema.ValidateSchema(); !ok {
		t.Fatalf("%v", err)
	}
}
