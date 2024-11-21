package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

type validator interface {
	Validate() error
}

// ValidateInterceptor проверяет наличие валидации protobuf,
// если она есть, то запускается валидация.
func ValidateInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	if val, ok := req.(validator); ok {
		err := val.Validate()
		if err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}
