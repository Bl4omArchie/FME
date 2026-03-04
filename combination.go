package fme



// Combination is a given set of flag bind to a schema
// 
type Combination struct {
	Set			map[string]struct{}
	
	Schema		*Schema
	
	Error		error
}


func NewCombination(flags map[string]struct{}, schema *Schema, err error) *Combination {
	return &Combination{
		Set: flags,
		Schema: schema,
		Error: err,
	}
}
