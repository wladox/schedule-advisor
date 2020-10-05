package app

import (
	"context"
	pb "github.com/wladox/schedule-advisor/internal/grpc"
	"github.com/wladox/schedule-advisor/pkg/forecast"
	"log"
)

type GrpcServer struct {
	pb.UnimplementedSchedulerAdvisorServer
	Calculator *forecast.Calculator
}

// GetCosts implements
func (s *GrpcServer) GetCosts(ctx context.Context, in *pb.ReschedulingRequest) (*pb.ReconfigurationCost, error) {
	log.Printf("Received: %v", in.GetTime())

	currAssignments := in.CurrentPlacement.Assignments
	newAssignments := in.NewPlacement.Assignments

	request := forecast.Request{
		CurrentPlacement: forecast.Placement{
			Assignments: currAssignments,
		},
		NewPlacement: forecast.Placement{
			Assignments: newAssignments,
		},
	}

	resp := s.Calculator.CalculateCosts(ctx, request)
	return &pb.ReconfigurationCost{
		CurrentCost: int64(resp.CurrentCosts),
		OptimalCost: int64(resp.OptimalCosts),
	}, nil
}
