package testutil

import (
	"context"
)

type Result[T any] struct {
	Value T
	Err   error
}

func Ok[T any](v T) (T, error) {
	return v, nil
}

func Err[T any](err error) (T, error) {
	var zero T
	return zero, err
}

func Ptr[T any](v T) *T {
	return &v
}

func GetCancelledCtx(ctx context.Context) context.Context {
	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()
	return canceledCtx
}
