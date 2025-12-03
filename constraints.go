package fme


import (
	"fmt"
	"errors"

	"github.com/katalvlaran/lvlath/bfs"
	"github.com/katalvlaran/lvlath/dfs"
	"github.com/katalvlaran/lvlath/core"
)


// A constraint is a type of rule that will be applied to every combination
// Each constraint can implement this interface in order to be used witth the Schema struct
type Constraint interface {
	// Add a new rule
	Add(a, b string, graph *core.Graph) error


	// Take every rules and verify there is no conflict between them and rules from other constraints
	// If there is a conflict, the schema is rolled back to keep its integrity
	ValidateSchema(graph *core.Graph) error


	// RollBack operation in case of schema integrity failure
	Rollback(a, b string, graph *core.Graph)


	// Take a given set of flags and apply each rules
	VerifyCombination(c *Combination, graph *core.Graph) error
}



// Require constraint represent a depedency between two flags
type Require struct {}

// Add two new vertex in the graph and connect them with a directed edge
func (r *Require) Add(a, b string, graph *core.Graph) error {
	ensureFlag(graph, a)
	ensureFlag(graph, b)

	// Unweighted dependency edge: weight = 0.
	_, _ = 	graph.AddEdge(string(a), string(b), 0, core.WithEdgeDirected(true))

	// Rollback if schema is invalid
	if err := r.ValidateSchema(graph); err != nil {
		r.Rollback(a, b, graph)
		return err
	}
	return nil
}

// Verify the integrity of the graph for Require constraint
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

// If the Schema is invalid, Rollback() deletes the two vertex in the graph
func (r *Require) Rollback(a, b string, graph *core.Graph) {
	graph.RemoveVertex(a)
	graph.RemoveVertex(b)
}

func (r *Require) VerifyCombination(c *Combination, g *core.Graph) error {
	for id := range c.Need {
		res, err := bfs.BFS(g, id)
		if err != nil {
			return fmt.Errorf("BFS from %q: %w", id, err)
		}
		for _, dep := range res.Order {
			c.Need[dep] = struct{}{}
		}
	}
	return nil
}



// Interfer constraint represent the impossibility for two flags to be mixed together into one combination
type Interfer struct {}

// Add two new vertex in the graph and connect them with a non-directed edge
// If A and B end up together in the closure of a combination, that combination is considered invalid.
func (i *Interfer) Add(a, b string, graph *core.Graph) error {
	ensureFlag(graph, a)
	ensureFlag(graph, b)

	if a == b {
		// A self-interfer does not make sense; ignore defensively.
		return fmt.Errorf("Self interference")
	}

	_, _ = 	graph.AddEdge(string(a), string(b), 0, core.WithEdgeDirected(false))

	// Rollback if schema is invalid
	if err := i.ValidateSchema(graph); err != nil {
		i.Rollback(a, b, graph)
		return err
	}

	return nil
}

func (i *Interfer) ValidateSchema(graph *core.Graph) error  {
 	sub := core.NewGraph(core.WithDirected(true))
    for _, v := range graph.Vertices() {
        _ = sub.AddVertex(v)
    }

    for _, e := range graph.Edges() {
        if e.Directed {
            _, _ = sub.AddEdge(e.From, e.To, e.Weight, core.WithEdgeDirected(true))
        }
    }

    // check DAG
    if _, err := dfs.TopologicalSort(sub); err != nil {
        if errors.Is(err, dfs.ErrCycleDetected) {
            return fmt.Errorf("dependency cycle detected")
        }
        return err
    }

    // check interference conflicts
    for _, e := range graph.Edges() {
        if !e.Directed {
            from, to := e.From, e.To
            if reachable(sub, from, to) || reachable(sub, to, from) {
                return fmt.Errorf("interference between %s and %s conflicts with requires edges", from, to)
            }
        }
    }

    return nil
}

func (i *Interfer) Rollback(a, b string, graph *core.Graph) {
	graph.RemoveVertex(a)
	graph.RemoveVertex(b)
}

func (i *Interfer) VerifyCombination(c *Combination, g *core.Graph) error {
	for from := range c.Need {
		for to := range c.Need {
			if from >= to {
				continue
			}
			if reachable(g, from, to) || reachable(g, to, from) {
				return fmt.Errorf("interference conflict between %q and %q", from, to)
			}
		}
	}
	return nil
}


// ensureFlag makes sure there is a vertex for id in the underlying graph.
// It is intentionally idempotent: calling it multiple times is safe.
func ensureFlag(g *core.Graph, id string) {
	_ = g.AddVertex(string(id))
}
