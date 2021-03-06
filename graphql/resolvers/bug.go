package resolvers

import (
	"context"
	"github.com/MichaelMure/git-bug/bug"
	"github.com/MichaelMure/git-bug/graphql/connections"
	"github.com/MichaelMure/git-bug/graphql/models"
)

type bugResolver struct{}

func (bugResolver) Status(ctx context.Context, obj *bug.Snapshot) (models.Status, error) {
	return convertStatus(obj.Status)
}

func (bugResolver) Comments(ctx context.Context, obj *bug.Snapshot, after *string, before *string, first *int, last *int) (models.CommentConnection, error) {
	input := models.ConnectionInput{
		Before: before,
		After:  after,
		First:  first,
		Last:   last,
	}

	edger := func(comment bug.Comment, offset int) connections.Edge {
		return models.CommentEdge{
			Node:   comment,
			Cursor: connections.OffsetToCursor(offset),
		}
	}

	conMaker := func(edges []models.CommentEdge, nodes []bug.Comment, info models.PageInfo, totalCount int) (models.CommentConnection, error) {
		return models.CommentConnection{
			Edges:      edges,
			Nodes:      nodes,
			PageInfo:   info,
			TotalCount: totalCount,
		}, nil
	}

	return connections.BugCommentCon(obj.Comments, edger, conMaker, input)
}

func (bugResolver) Operations(ctx context.Context, obj *bug.Snapshot, after *string, before *string, first *int, last *int) (models.OperationConnection, error) {
	input := models.ConnectionInput{
		Before: before,
		After:  after,
		First:  first,
		Last:   last,
	}

	edger := func(op bug.Operation, offset int) connections.Edge {
		return models.OperationEdge{
			Node:   op,
			Cursor: connections.OffsetToCursor(offset),
		}
	}

	conMaker := func(edges []models.OperationEdge, nodes []bug.Operation, info models.PageInfo, totalCount int) (models.OperationConnection, error) {
		return models.OperationConnection{
			Edges:      edges,
			Nodes:      nodes,
			PageInfo:   info,
			TotalCount: totalCount,
		}, nil
	}

	return connections.BugOperationCon(obj.Operations, edger, conMaker, input)
}
