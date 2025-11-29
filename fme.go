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
	"fmt"
	"sort"
	"errors"

	"github.com/katalvlaran/lvlath/bfs"
	"github.com/katalvlaran/lvlath/core"
	"github.com/katalvlaran/lvlath/dfs"
)

// Schema holds the static constraint model:
//
//   - a directed dependency graph (requires edges) over string values;
//   - a symmetric interfer relation stored as a map of sets.
//
// The graph is used for:
//   - computing the transitive closure of dependencies via BFS;
//   - checking for cycles and building execution order via DFS.
//
// The interfer map is kept separate for simplicity and cheaper lookups.
type Schema struct {
	Graph *core.Graph
	Constraints []Constraint
}

// CombinationResult is the outcome of validating a concrete user combination.
//
//   - Final   is the closure(combination): selected Flags + all dependencies;
//   - Interfer is non-nil if a interfering pair was detected.
type CombinationResult struct {
	Final    []string
	Interfer *PathInstance
}


// NewSchema constructs a directed, unweighted dependency graph backed by
// lvlath/core and an empty interfer relation.
//
// Complexity of operations on this structure is dominated by BFS/DFS:
//   - Schema validation: O(V + E + C * (V + E)) in the worst case;
//   - Combination validation: O(|S| * (V + E) + C) for practical sizes.
func NewSchema() *Schema {
	return &Schema{
		Graph:	core.NewGraph(core.WithDirected(true)),
		Constraints: make([]Constraint, 0),
	}
}


func (s *Schema) AddConstraint(constraints ...Constraint) {
	for _, c := range constraints {
		s.Constraints = append(s.Constraints, c)
	}
}


// ValidateSchema performs static validation of the constraint schema.
//
//  1. Ensures that the dependency graph is a DAG via dfs.TopologicalSort.
//  2. Ensures that no interfer pair {A, B} is such that one is reachable
//     from the other through “requires” edges (which would be a contradiction).
//
// This function is intended to be called once at service startup.
func (s *Schema) ValidateSchema() error {
    for _, c := range s.Constraints {
        if err := c.ValidateSchema(s.Graph); err != nil {
            return err
        }
    }
    return nil
}


// ValidateCombination:
//
//   - expands the initial combination by adding all transitive dependencies;
//   - checks for Interferences in the resulting closure;
//   - returns CombinationResult and either nil or ErrCombinationInterfer.
//
// This is the function you would typically call per user request.
func (s *Schema) ValidateCombination(flags []string) (*CombinationResult, error) {
    need := make(map[string]struct{})

    for _, id := range flags {
        // Ensure vertex exists (matching your previous ensureFlag logic)
        if !s.Graph.HasVertex(id) {
            return nil, fmt.Errorf("unknown flag %q", id)
        }

        // BFS closure
        res, err := bfs.BFS(s.Graph, id)
        if err != nil {
            return nil, fmt.Errorf("expand: BFS from %q failed: %w", id, err)
        }

        for _, v := range res.Order {
            need[string(v)] = struct{}{}
        }
    }

    for _, c := range s.Constraints {
        if err := c.VerifyCombination(need, s.Graph); err != nil {
            return nil, err
        }
    }
	final := make([]string, 0, len(need))
	for id := range need {
		final = append(final, id)
	}
	sort.Slice(final, func(i, j int) bool { return final[i] < final[j] })

    return &CombinationResult{Final: final}, nil
}


// ExecutionOrder computes a deterministic execution order for the subset
// of Flags given in Flags, respecting all dependency edges.
//
// Internally it builds an induced subgraph on the subset and runs a
// topological sort via dfs.TopologicalSort.
func (s *Schema) ExecutionOrder(Flags []string) ([]string, error) {
	if len(Flags) == 0 {
		return nil, nil
	}

	needed := make(map[string]struct{}, len(Flags))
	for _, id := range Flags {
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
				Detail: "cycle detected in induced subgraph for combination",
			}
		}
		return nil, err
	}

	out := make([]string, len(order))
	for i, v := range order {
		out[i] = string(v)
	}
	return out, nil
}
