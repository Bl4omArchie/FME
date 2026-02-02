package fme

import (
	"testing"
)


func TestSchema(t *testing.T) {
	schema := NewSchema()

	err := schema.Constraints[0].Add("a", "b", schema)
	err =  schema.Constraints[0].Add("b", "c", schema)
	if err != nil {
		t.Fatalf("incorrect schema : %v", err)
	}

	_, _ = schema.ValidateCombination([]string{"a"})

	// order, err := schema.ExecutionOrder(c)
	// if err != nil {
	// 	t.Fatalf("%v", err)
	// } else {
	// 	fmt.Println(order)
	// }
}
