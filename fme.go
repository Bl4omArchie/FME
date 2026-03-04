package fme

// High-level pseudocode (how lvlath is used to solve OTO-style constraints):
//
//  1. Build a directed dependency graph G over Flag IDs using core.Graph.
//     For each rule "A requires B", add a directed edge A -> B with weight 0.
//  2. Maintain a separate symmetric interfer relation C ⊆ V×V in a Go map.
//  3. At service startup:
//     a) Run dfs.TopologicalSort(G) to ensure that dependencies form a DAG.
//        If a cycle is found, reject the schema (misconfigured rules).
//     b) For each interfer {a, b} in C, run bfs.BFS(G, a) and bfs.BFS(G, b).
//        If a and b are mutually reachable via "requires" edges, reject
//        the schema as contradictory.
//  4. For each user combination S:
//     a) Compute closure(S) by running BFS from every a ∈ S over G,
//        adding all reachable Flags to the set "need".
//     b) If "need" contains any interfering pair {a, b} ∈ C, reject the
//        combination, optionally explaining the shortest implication chain.
//     c) Otherwise, build an induced subgraph on "need" and run
//        dfs.TopologicalSort to get a safe execution order for tasks.
//
// This example is a skeleton for an "Flag-constraints engine" that can
// sit underneath a scheduler like OTO while staying small, explicit, and
// easy to explain to users.

import (
	"errors"
	"fmt"

	"github.com/katalvlaran/lvlath/core"
	"github.com/katalvlaran/lvlath/dfs"
)

// Schema holds the static constraint model:
//
//   - a directed dependency graph (requires edges) over string values;
//   - a slice of constraints
//
// The graph is used for:
//   - computing the transitive closure of dependencies via BFS;
//   - checking for cycles and building execution order via DFS.
//
// The interfer map is kept separate for simplicity and cheaper lookups.
type Schema struct {
	Graph			*core.Graph

	Interferences 	map[string]map[string]struct{}
	
	Constraints		map[string]Constraint
}


// NewSchema constructs a directed, unweighted dependency graph backed by
// lvlath/core and an empty interfer relation.
//
// Complexity of operations on this structure is dominated by BFS/DFS:
//   - Schema validation: O(V + E + C * (V + E)) in the worst case;
//   - Combination validation: O(|S| * (V + E) + C) for practical sizes.
func NewSchema(constraints map[string]Constraint) *Schema {
	return &Schema{
		Graph:			core.NewGraph(core.WithDirected(true)),
		Interferences:  map[string]map[string]struct{}{},
		Constraints:	constraints,
	}
}

func InitSchema() *Schema {
	return NewSchema(map[string]Constraint{"req": &Require{}, "interfer": &Interfer{}})
}


func (s *Schema) Add(key, a, b string) error {
	if ok := s.Constraints[key]; ok == nil {
		return fmt.Errorf("incorrect constraint key, couldn't find : %s", key)
	} else {
		return ok.Add(a, b, s)
	}
}


// ValidateSchema performs static validation of the constraint schema.
//
//  1. Ensures that the dependency graph is a DAG via dfs.TopologicalSort.
//  2. Ensures that no interfer pair {A, B} is such that one is reachable
//     from the other through “requires” edges (which would be a contradiction).
//
// This function is intended to be called once at service startup.
func (s *Schema) ValidateSchema() (bool, error) {
    for _, c := range s.Constraints {
        if err := c.SchemaValidation(s); err != nil {
            return false, err
        }
    }
    return true, nil
}


// ValidateCombination:
//
//   - expands the initial combination by adding all transitive dependencies;
//   - checks for Interferences in the resulting closure;
//   - returns CombinationResult and either nil or ErrCombinationInterfer.
//
// This is the function you would typically call per user request.
func (s *Schema) ValidateCombination(flags map[string]struct{}) (bool, error) {
    for _, c := range s.Constraints {
        if err := c.CombinationValidation(flags, s); err != nil {
            return false, err
        }
    }

	return true, nil
}


// ExecutionOrder computes a deterministic execution order for the subset
// of arguments given in args, respecting all dependency edges.
//
// Internally it builds an induced subgraph on the subset and runs a
// topological sort via dfs.TopologicalSort.
func (s *Schema) ExecutionOrder(c *Combination) ([]string, error) {
	if len(c.Set) == 0 {
		return nil, nil
	}

	needed := make(map[string]struct{}, len(c.Set))
	for id := range c.Set {
		needed[string(id)] = struct{}{}
	}

	// Build an induced subgraph on the needed vertices.
	sub := s.Graph.CloneEmpty()
	for id := range needed {
		_ = sub.AddVertex(id)
	}

	for _, e := range s.Graph.Edges() {
		if _, ok := needed[e.From]; !ok {
			continue
		}
		if _, ok := needed[e.To]; !ok {
			continue
		}
		_, _ = sub.AddEdge(e.From, e.To, e.Weight, core.WithEdgeDirected(e.Directed))
	}

	order, err := dfs.TopologicalSort(sub)
	if err != nil {
		if errors.Is(err, dfs.ErrCycleDetected) {
			// Should not happen if ValidateSchema has already passed.
			return nil, &SchemaValidationError{
				Kind:   ErrSchemaCycle,
				Detail: "cycle detected in induced subgraph for selection",
			}
		}
		return nil, err
	}

	return order, nil
}
