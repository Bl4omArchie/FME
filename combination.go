package fme



// Combination is a given set of flag bind to a schema
// 
type Combination struct {
	Set			map[string]struct{}
	
	Schema		*Schema
	
	Valid		bool
}


func NewCombination(flags []string, s *Schema) (*Combination, error) {
	set := make(map[string]struct{}, len(flags))
	for _, v := range flags {
		set[v] = struct{}{}
	}

	ok, err := s.ValidateCombination(set)

	return &Combination{
		Set: set,
		Schema: s,
		Valid: ok,
	}, err
}


func (c *Combination) Update() (bool, error){
	ok, err := c.Schema.ValidateCombination(c.Set)
	c.Valid = ok
	return ok, err
}
