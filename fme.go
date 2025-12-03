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
	"sort"

	"github.com/katalvlaran/lvlath/bfs"
	"github.com/katalvlaran/lvlath/core"
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
	Graph *core.Graph
	Constraints []Constraint
}

// CombinationResult is the outcome of validating a concrete user combination.
//
//   - Set is the closure(combination): selected Flags + all dependencies;
type Combination struct {
	Set []string
}

// NewSchema constructs a directed, unweighted dependency graph backed by
// lvlath/core and an empty interfer relation.
//
// Complexity of operations on this structure is dominated by BFS/DFS:
//   - Schema validation: O(V + E + C * (V + E)) in the worst case;
//   - Combination validation: O(|S| * (V + E) + C) for practical sizes.
func NewSchema() *Schema {
	var constraints []Constraint = []Constraint{&Require{}, &Interfer{}} 
	return &Schema{
		Graph:			core.NewGraph(core.WithMixedEdges()),
		Constraints:	constraints,
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
func (s *Schema) ValidateSchema() (bool, error) {
    for _, c := range s.Constraints {
        if err := c.ValidateSchema(s.Graph); err != nil {
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
func (s *Schema) ValidateCombination(flags []string) (*Combination, error) {
    need := make(map[string]struct{})
	var combination *Combination
	combination.Set = flags

    for _, c := range s.Constraints {
        if err := c.VerifyCombination(combination, s.Graph); err != nil {
            return nil, err
        }
    }

	final := make([]string, 0, len(need))
	for id := range need {
		final = append(final, id)
	}
	sort.Slice(final, func(i, j int) bool { return final[i] < final[j] })

    return &Combination{Set: final}, nil
}


// reachable runs BFS once and checks whether "to" is reachable from "from"
// using the Depth map from BFSResult. This is enough for schema-level checks
// where only reachability matters, not the exact path.
func reachable(g *core.Graph, from, to string) bool {
	res, err := bfs.BFS(g, string(from))
	if err != nil {
		// In a well-formed schema this should not fail; treat as not reachable.
		return false
	}

	_, ok := res.Depth[string(to)]
	return ok
}
