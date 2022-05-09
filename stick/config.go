package stick

import (
	"context"
)

type ConfigCtxKey string

type Config struct {
	Worker func(job func())
}

func GetConfigFrom(ctx context.Context) Config {
	return GetFrom[Config](ctx, ConfigCtxKey(""))
}

func WithConfig(ctx context.Context, cfg Config) context.Context {
	return With(ctx, ConfigCtxKey(""), cfg)
}

var defaultConfig = Config{
	Worker: func(job func()) {
		go func() {
			job()
		}()
	},
}
