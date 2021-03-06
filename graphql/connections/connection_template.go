package connections

import (
	"fmt"

	"github.com/MichaelMure/git-bug/graphql/models"
	"github.com/cheekybits/genny/generic"
)

// NodeType define the node type handled by this relay connection
type NodeType generic.Type

// EdgeType define the edge type handled by this relay connection
type EdgeType generic.Type

// ConnectionType define the connection type handled by this relay connection
type ConnectionType generic.Type

// NodeTypeEdger define a function that take a NodeType and an offset and
// create an Edge.
type NodeTypeEdger func(value NodeType, offset int) Edge

// NodeTypeConMaker define a function that create a ConnectionType
type NodeTypeConMaker func(
	edges []EdgeType,
	nodes []NodeType,
	info models.PageInfo,
	totalCount int) (ConnectionType, error)

// NodeTypeCon will paginate a source according to the input of a relay connection
func NodeTypeCon(source []NodeType, edger NodeTypeEdger, conMaker NodeTypeConMaker, input models.ConnectionInput) (ConnectionType, error) {
	var nodes []NodeType
	var edges []EdgeType
	var pageInfo models.PageInfo

	emptyCon, _ := conMaker(edges, nodes, pageInfo, 0)

	offset := 0

	if input.After != nil {
		for i, value := range source {
			edge := edger(value, i)
			if edge.GetCursor() == *input.After {
				// remove all previous element including the "after" one
				source = source[i+1:]
				offset = i + 1
				break
			}
		}
	}

	if input.Before != nil {
		for i, value := range source {
			edge := edger(value, i+offset)

			if edge.GetCursor() == *input.Before {
				// remove all after element including the "before" one
				break
			}

			edges = append(edges, edge.(EdgeType))
			nodes = append(nodes, value)
		}
	} else {
		edges = make([]EdgeType, len(source))
		nodes = source

		for i, value := range source {
			edges[i] = edger(value, i+offset).(EdgeType)
		}
	}

	if input.First != nil {
		if *input.First < 0 {
			return emptyCon, fmt.Errorf("first less than zero")
		}

		if len(edges) > *input.First {
			// Slice result to be of length first by removing edges from the end
			edges = edges[:*input.First]
			nodes = nodes[:*input.First]
			pageInfo.HasNextPage = true
		}
	}

	if input.Last != nil {
		if *input.Last < 0 {
			return emptyCon, fmt.Errorf("last less than zero")
		}

		if len(edges) > *input.Last {
			// Slice result to be of length last by removing edges from the start
			edges = edges[len(edges)-*input.Last:]
			nodes = nodes[len(nodes)-*input.Last:]
			pageInfo.HasPreviousPage = true
		}
	}

	return conMaker(edges, nodes, pageInfo, len(source))
}
