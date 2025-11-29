package fme

import (
	"fmt"
	"errors"

	"github.com/katalvlaran/lvlath/dfs"
	"github.com/katalvlaran/lvlath/core"
)

// A constraint is a type of rule that will be applied to every combination
// Each constraint can implement this interface in order to be used witth the Schema struct
type Constraint interface {
	// Add a new rule
	Add(a, b string) error


	// Take every rules and verify there is no conflict between them and rules from other constraints
	// If there is a conflict, the schema is rolled back to keep its integrity
	ValidateSchema(graph *core.Graph) error


	// RollBack operation in case of schema integrity failure
	Rollback(a, b string, graph *core.Graph)


	// Take a given set of flags and apply each rules
	VerifyCombination(combination map[string]struct{}, graph *core.Graph) error
}


// Require constraint represent a depedency between two flags
type Require struct {}
func (r *Require) Add(a, b string) error {
	return nil
}
func (r *Require) ValidateSchema(graph *core.Graph) error {
	if _, err := dfs.TopologicalSort(graph); err != nil {
		if errors.Is(err, dfs.ErrCycleDetected) {
			return &SchemaValidationError{
				Kind:   ErrSchemaCycle,
				Detail: "invalid schema: dependency cycle detected in Flag graph",
			}
		}
		// Any other error is unexpected and should be surfaced as-is.
		return err
	}
	return nil
}
func (r *Require) Rollback(a, b string, graph *core.Graph) {
	graph.RemoveVertex(a)
	graph.RemoveVertex(b)
}
func (r *Require) VerifyCombination(combination map[string]struct{}, graph *core.Graph) error {
	return nil
}


// Interfer constraint represent the impossibility for two flags to be mixed together into one combination
type Interfer struct {
	Interferences map[string]map[string]struct{}
}
// Interfer registers a symmetric interfer between A and B.
// If A and B end up together in the closure of a combination, that combination
// is considered invalid.
func (i *Interfer) Add(a, b string) error {
	if a == b {
		// A self-interfer does not make sense; ignore defensively.
		return fmt.Errorf("Self interference")
	}

	if i.Interferences[a] == nil {
		i.Interferences[a] = make(map[string]struct{})
	}
	if i.Interferences[b] == nil {
		i.Interferences[b] = make(map[string]struct{})
	}
	i.Interferences[a][b] = struct{}{}
	i.Interferences[b][a] = struct{}{}

	return nil
}

func (i *Interfer) ValidateSchema(graph *core.Graph) error  {
	for a, row := range i.Interferences {
		for b := range row {
			// Work with each unordered pair only once (a < b).
			if a >= b {
				continue
			}

			if reachable(graph, a, b) || reachable(graph, b, a) {
				msg := fmt.Sprintf(
					"invalid schema: %q and %q are declared as interfering, "+
						"but one is reachable from the other via requires edges",
					a, b,
				)
				return &SchemaValidationError{
					Kind:   ErrSchemaContradiction,
					Detail: msg,
				}
			}
		}
	}
	return nil
}

func (i *Interfer) Rollback(a, b string, graph *core.Graph) {
	delete(i.Interferences[a], b)
	delete(i.Interferences[b], a)

	if len(i.Interferences[a]) == 0 {
		delete(i.Interferences, a)
	}
	if len(i.Interferences[b]) == 0 {
		delete(i.Interferences, b)
	}
}

func (i *Interfer) VerifyCombination(combination map[string]struct{}, graph *core.Graph) error {
	for a, row := range i.Interferences {
		for b := range row {
			if a >= b {
				continue
			}
			_, hasA := combination[a]
			_, hasB := combination[b]
			if hasA && hasB {
				return &CombinationVerificationError{
					Kind: ErrCombinationInterfer,
					Detail: "Your interference is creating a conflict with the graph",
					Path: ShortestPath(a, b, graph),
				}
			}
		}
	}
	return nil
}


// type Position struct {
// 	SortedSlice []string
// }

// type Scale struct {
// 	Name string
// 	Description string
// 	Values map[string]int
// }

// func NewScale(name, description string, values map[string]int) *Scale {
// 	return &Scale {
// 		Name: name,
// 		Description: description,
// 		Values: values,
// 	}
// }

// // Compute the weigth of a given combination
// func (s *Scale) ComputeWeigth(c []string) int {
// 	var total int = 0
// 	for _, flag := range c {
// 		f, ok := i.Values[flag]
// 		if ok {
// 			total += f
// 		}
// 	}
// 	return total
// }

//  //=== TODO ===

// // Among a population of flag, return the highest possible combination
// func (s *Scale) GetHighestCombination(p []string) *Combination {
// 	return NewCombination(p)
// }

// // Among a population of flag, return the lowest possible combination
// func (s *Scale) GetLowestCombination(p []string) *Combination {
// 	return NewCombination(p)
// }
