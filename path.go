package fme

import (
	"github.com/katalvlaran/lvlath/bfs"
	"github.com/katalvlaran/lvlath/core"
)

type PathInstance struct {
	A, B string
	Path []string
}

func NewPathInstance(a, b string, path []string) *PathInstance {
	return &PathInstance{
		A:    a,
		B:    b,
		Path: path,
	}
}

// shortestPath reconstructs a shortest path (in edges) between from and to
// in the unweighted dependency graph, using the PathTo helper from BFSResult.
// If there is no path, it returns nil.
func ShortestPath(a, b string, g *core.Graph) *PathInstance {
	res, err := bfs.BFS(g, string(a))
	if err != nil {
		return nil
	}

	// Use the BFSResult helper, which already relies on Depth/Parent
	// and returns an error when dest has not been reached.
	strPath, err := res.PathTo(string(b))
	if err != nil {
		return nil
	}

	path := make([]string, len(strPath))
	for i, v := range strPath {
		path[i] = string(v)
	}
	return NewPathInstance(a, b, path)
}
