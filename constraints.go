package fme

import (
	"errors"
	"fmt"
	"log"
	"sort"

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
	ValidateSchema(s *Schema) error


	// RollBack operation in case of schema integrity failure
	Rollback(a, b string, s *Schema)


	// Take a given set of flags and apply each rules
	VerifyCombination(c *Combination, s *Schema) *CombinationError
}



// Require constraint represent a depedency between two flags
type Require struct {}

// Add two new vertex in the graph and connect them with a directed edge
func (r *Require) Add(a, b string, s *Schema) error {
	ensureFlag(s.Graph, a)
	ensureFlag(s.Graph, b)

	_, _ = s.Graph.AddEdge(string(a), string(b), 0)

	// Rollback if schema is invalid
	if err := r.ValidateSchema(s); err != nil {
		r.Rollback(a, b, s)
		return err
	}
	return nil
}

// Verify the integrity of the graph for Require constraint
func (r *Require) ValidateSchema(s *Schema) error {
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

// If the Schema is invalid, Rollback() deletes the two vertex in the graph
func (r *Require) Rollback(a, b string, s *Schema) {
	s.Graph.RemoveVertex(a)
	s.Graph.RemoveVertex(b)
}

func (r *Require) VerifyCombination(c *Combination, s *Schema) *CombinationError {
	need := make(map[string]struct{})

	// Make sure all selected Flags exist as vertices.
	for _, id := range c.Set {
		ensureFlag(s.Graph, id)
	}

	// For each selected Flag, run BFS to collect all dependencies.
	for _, id := range c.Set {
		res, err := bfs.BFS(s.Graph, string(id))
		log.Println(res)
		
		if err != nil {
			path := NewPathInstance(res.Order[0], res.Order[1], res.Order)
			return NewCombinationError(path, err)
		}
		for _, ID := range res.Order {
			need[string(ID)] = struct{}{}
		}
	}

	c.Need = need

	// Convert to a sorted slice for deterministic output and testinGraph.
	final := make([]string, 0, len(need))
	for id := range need {
		final = append(final, id)
	}
	sort.Slice(final, func(i, j int) bool { return final[i] < final[j] })

	log.Println(final, need)

	return nil
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
	if err := i.ValidateSchema(s); err != nil {
		i.Rollback(a, b, s)
		return err
	}

	return nil
}

func (i *Interfer) ValidateSchema(s *Schema) error  {
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

func (i *Interfer) Rollback(a, b string, s *Schema) {
	s.Graph.RemoveVertex(a)
	s.Graph.RemoveVertex(b)
}

func (i *Interfer) VerifyCombination(c *Combination, s *Schema) *CombinationError {
	for a, row := range s.Interferences {
		for b := range row {
			if a >= b {
				continue
			}
			_, hasA := c.Need[a]
			_, hasB := c.Need[b]
			if hasA && hasB {
				return NewCombinationError(ShortestPath(a, b, s.Graph), ErrCombinationInterfer)
			}
		}
	}
	return nil
}
