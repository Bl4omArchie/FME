package fme

import (
	"errors"
	"fmt"

	"github.com/katalvlaran/lvlath/bfs"
	"github.com/katalvlaran/lvlath/dfs"
)

// A constraint is a type of rule that will be applied to every combination
// Each constraint can implement this interface in order to be used witth the Schema struct
type Constraint interface {
	// Add a new rule
	Add(a, b string, s *Schema) error


	// Take every rules and verify there is no conflict between them and rules from other constraints
	// If there is a conflict, the schema is rolled back to keep its integrity
	SchemaValidation(s *Schema) error


	// Take a given set of flags and apply each rules
	CombinationValidation(c map[string]struct{}, s *Schema) error


	// RollBack operation in case of schema integrity failure
	Rollback(a, b string, s *Schema)
}



// Require constraint represent a depedency between two flags
type Require struct {}

// Add two new vertex in the graph and connect them with a directed edge
func (r *Require) Add(a, b string, s *Schema) error {
	ensureFlag(s.Graph, a)
	ensureFlag(s.Graph, b)

	_, _ = s.Graph.AddEdge(string(a), string(b), 0)

	// Rollback if schema is invalid
	if err := r.SchemaValidation(s); err != nil {
		r.Rollback(a, b, s)
		return err
	}
	return nil
}


func (r *Require) SchemaValidation(s *Schema) error {
	if _, err := dfs.TopologicalSort(s.Graph); err != nil {
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


func (r *Require) CombinationValidation(c map[string]struct{}, s *Schema) error {
	need := make(map[string]struct{})

	// Make sure all selected Flags exist as vertices.
	for id := range c {
		ensureFlag(s.Graph, id)
	}

	// For each selected Flag, run BFS to collect all dependencies.
	for flagID := range c {
		res, err := bfs.BFS(s.Graph, string(flagID))
		if err != nil {
			return &CombinationVerificationError{
				Kind: ErrCombinationRequire,
				Detail: "BFS failed for requirement constraints combination verification",
				Path: nil,
			}
		}
		for _, ID := range res.Order {
			depID := string(ID)
			need[depID] = struct{}{}

			if _, ok := c[depID]; !ok {
				return &CombinationVerificationError{
					Kind:   ErrCombinationRequire,
					Detail: "missing required dependency: " + depID,
					Path:   ShortestPath(flagID, depID, s.Graph),
				}
			}
		}
	}

	return nil
}


func (r *Require) Rollback(a, b string, s *Schema) {
	s.Graph.RemoveVertex(a)
	s.Graph.RemoveVertex(b)
}



// Interfer constraint represent the impossibility for two flags to be mixed together into one combination
type Interfer struct {}

// Add two new vertex in the graph and connect them with a non-directed edge
// If A and B end up together in the closure of a combination, that combination is considered invalid.
func (i *Interfer) Add(a, b string, s *Schema) error {
	ensureFlag(s.Graph, a)
	ensureFlag(s.Graph, b)

	if a == b {
		// A self-interfer does not make sense; ignore defensively.
		return fmt.Errorf("Self interference")
	}

	if s.Interferences[a] == nil {
		s.Interferences[a] = make(map[string]struct{})
	}
	if s.Interferences[b] == nil {
		s.Interferences[b] = make(map[string]struct{})
	}
	s.Interferences[a][b] = struct{}{}
	s.Interferences[b][a] = struct{}{}

	// Rollback if schema is invalid
	if err := i.SchemaValidation(s); err != nil {
		i.Rollback(a, b, s)
		return err
	}

	return nil
}


func (i *Interfer) SchemaValidation(s *Schema) error  {
	for a, row := range s.Interferences {
		for b := range row {
			// Work with each unordered pair only once (a < b).
			if a >= b {
				continue
			}
			
			if reachable(s.Graph, a, b) || reachable(s.Graph, b, a) {
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


func (i *Interfer) CombinationValidation(c map[string]struct{}, s *Schema) error {
	for a, row := range s.Interferences {
		for b := range row {
			if a >= b {
				continue
			}
			_, hasA := c[a]
			_, hasB := c[b]
			if hasA && hasB {
				return &CombinationVerificationError{
					Kind: ErrCombinationInterfer,
					Detail: "Failed combination verification for interference constraints",
					Path: ShortestPath(a, b, s.Graph),
				}
			}
		}
	}
	return nil
}


func (i *Interfer) Rollback(a, b string, s *Schema) {
	s.Graph.RemoveVertex(a)
	s.Graph.RemoveVertex(b)
}
