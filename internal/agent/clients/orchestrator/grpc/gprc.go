package grpc

import (
	"context"
	"fmt"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	distributedcalculatorv1 "github.com/k6mil6/distributed-calculator-protobuf/gen/go/distributed-calculator"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type Client struct {
	api distributedcalculatorv1.OrchestratorClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	address string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	op := "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadSent, grpclog.PayloadReceived),
	}

	cc, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: distributedcalculatorv1.NewOrchestratorClient(cc),
	}, nil
}

func (c *Client) GetFreeExpressions(ctx context.Context) (model.Subexpression, error) {
	op := "grpc.GetFreeExpressions"

	log := c.log.With(slog.String("op", op))

	log.Info("requesting free expressions")

	resp, err := c.api.GetFreeExpressions(ctx, &distributedcalculatorv1.GetFreeExpressionsRequest{})
	if err != nil {
		return model.Subexpression{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("free expressions received")

	return model.Subexpression{
		ID:            int(resp.SubexpressionID),
		Subexpression: resp.Subexpression,
		Timeout:       resp.Timeout,
		WorkerId:      int(resp.WorkerID),
	}, nil
}

func (c *Client) SaveResult(ctx context.Context, subexpressionID int, result float64) (int, error) {
	op := "grpc.SaveResult"

	log := c.log.With(slog.String("op", op))

	log.Info("sending result")

	resp, err := c.api.SendResult(ctx, &distributedcalculatorv1.SendResultRequest{
		SubexpressionID: int32(subexpressionID),
		Result:          result,
	})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("result sent")

	return int(resp.SubexpressionID), nil
}

func (c *Client) SendHeartbeat(ctx context.Context) error {
	op := "grpc.SendHeartbeat"

	log := c.log.With(slog.String("op", op))

	log.Info("sending heartbeat")

	_, err := c.api.SendHeartbeat(ctx, &distributedcalculatorv1.SendHeartbeatRequest{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("heartbeat sent")

	return nil
}
func InterceptorLogger(log *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), msg, fields...)
	})
}
