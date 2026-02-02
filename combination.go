package fme



// CombinationResult is the outcome of validating a concrete user combination.
//
//  - Set is the closure(combination): selected Flags + all dependencies;
//	- Need are depe
type Combination struct {
	Set		[]string
	
	Need	map[string]struct{}
}


type CombinationError struct {
	Path *PathInstance
	Error error
}


func NewCombination(flags []string) *Combination {
	return &Combination{
		Set: flags,
		Need: make(map[string]struct{}),
	}
}

func NewCombinationError(path *PathInstance, err error) *CombinationError {
	return &CombinationError{
		Path: path,
		Error: err,
	}
}