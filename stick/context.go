package stick

import (
	"context"
)

func With(ctx context.Context, key, val any) context.Context {
	return context.WithValue(ctx, key, val)
}

func GetFrom[T any](ctx context.Context, key any) T {
	return ctx.Value(key).(T)
}

func Has(ctx context.Context, key any) bool {
	return ctx.Value(key) != nil
}
