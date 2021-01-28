package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: add tests

func GrpcErrFromContext(ctx context.Context) error {
	switch err := ctx.Err(); err {
	case nil:
		return nil
	case context.Canceled:
		return status.Error(codes.Aborted, err.Error())
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	default:
		panic("unreachable")
	}
}

func ArgErrf(format string, a ...interface{}) error {
	return status.Errorf(codes.InvalidArgument, format, a...)
}

// Returns true is `err` is a gRPC status with code AlreadyExists.
func IsAlreadyExists(err error) bool {
	s, ok := status.FromError(err)
	return ok && s.Code() == codes.AlreadyExists
}
