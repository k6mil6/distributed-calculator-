package grpc

import (
	"context"
	"errors"
	distributedcalculatorv1 "github.com/k6mil6/distributed-calculator-protobuf/gen/go/distributed-calculator"
	errs "github.com/k6mil6/distributed-calculator/internal/errors"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Orchestrator interface {
	GetFreeExpressions(ctx context.Context) (model.Subexpression, error)
	SaveResult(ctx context.Context, subexpressionID int, result float64) (int, error)
}

type Heartbeat interface {
	SaveHeartbeat(ctx context.Context, workerID int) error
}

type serverApi struct {
	distributedcalculatorv1.UnimplementedOrchestratorServer
	orchestrator Orchestrator
	heartbeat    Heartbeat
}

func Register(gRPC *grpc.Server, orchestrator Orchestrator, heartbeat Heartbeat) {
	distributedcalculatorv1.RegisterOrchestratorServer(gRPC, &serverApi{
		orchestrator: orchestrator,
		heartbeat:    heartbeat,
	})
}

func (s *serverApi) GetFreeExpressions(
	ctx context.Context,
	request *distributedcalculatorv1.GetFreeExpressionsRequest,
) (*distributedcalculatorv1.GetFreeExpressionsResponse, error) {
	freeExpression, err := s.orchestrator.GetFreeExpressions(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrSubexpressionNotFound) {
			return nil, status.Error(codes.NotFound, "error getting free expressions")
		}
		return nil, status.Error(codes.Internal, "error getting free expressions")
	}

	if request.GetWorkerID() != 0 {
		freeExpression.WorkerId = int(request.GetWorkerID())
	}

	return &distributedcalculatorv1.GetFreeExpressionsResponse{
		SubexpressionID: int32(freeExpression.ID),
		Subexpression:   freeExpression.Subexpression,
		Timeout:         freeExpression.Timeout,
		WorkerID:        int32(freeExpression.WorkerId),
	}, nil
}

func (s *serverApi) SendResult(
	ctx context.Context,
	request *distributedcalculatorv1.SendResultRequest,
) (*distributedcalculatorv1.SendResultResponse, error) {
	id, err := s.orchestrator.SaveResult(ctx, int(request.GetSubexpressionID()), request.GetResult())

	if err != nil {
		return nil, status.Error(codes.Internal, "error saving result")
	}

	return &distributedcalculatorv1.SendResultResponse{SubexpressionID: int32(id)}, nil
}

func (s *serverApi) SendHeartbeat(
	ctx context.Context,
	request *distributedcalculatorv1.SendHeartbeatRequest,
) (*distributedcalculatorv1.Empty, error) {
	if err := s.heartbeat.SaveHeartbeat(ctx, int(request.GetWorkerID())); err != nil {
		return nil, status.Error(codes.Internal, "error saving heartbeat")
	}

	return &distributedcalculatorv1.Empty{}, nil
}
