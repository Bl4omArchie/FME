package fme

import (
	"testing"
)


func TestSchema(t *testing.T) {
	schema := NewSchema()

	err := schema.Constraints[0].Add("a", "b", schema)
	err =  schema.Constraints[0].Add("b", "c", schema)
	err =  schema.Constraints[0].Add("f", "d", schema)
	if err != nil {
		t.Fatalf("incorrect schema : %v", err)
	}

	c, err := NewCombination([]string{"a", "b", "c"}, schema)
	if err != nil {
		t.Fatalf("incorrect combination : %v", err)
	} else {
		t.Log(c.Valid)
	}

	order, err := schema.ExecutionOrder(c)
	if err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Log(order)
	}
}
