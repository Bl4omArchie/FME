package fme


import (
	"github.com/katalvlaran/lvlath/bfs"
	"github.com/katalvlaran/lvlath/core"
)


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


// ensureFlag makes sure there is a vertex for id in the underlying graph.
// It is intentionally idempotent: calling it multiple times is safe.
func ensureFlag(g *core.Graph, id string) {
	_ = g.AddVertex(string(id))
}

