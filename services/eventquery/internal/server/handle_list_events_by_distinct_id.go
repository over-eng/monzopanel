package server

import (
	"context"

	"github.com/over-eng/monzopanel/protos/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ListEventsByDistinctID(
	ctx context.Context,
	req *event.ListEventsByDistinctIDRequest,
) (*event.ListEventsByDistinctIDResponse, error) {
	if req.TeamId == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}

	if req.DistinctId == "" {
		return nil, status.Error(codes.InvalidArgument, "distinct_id is required")
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	result, err := s.eventstore.ListEventsByDistinctID(
		req.TeamId,
		req.DistinctId,
		pageSize,
		req.PaginationToken,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events: %v", err)
	}

	response := &event.ListEventsByDistinctIDResponse{
		Events:              result.Events,
		NextPaginationToken: result.NextPaginationToken,
	}

	return response, nil
}
