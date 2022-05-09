package stick

import (
	"context"
)

func Spin(spinCtx context.Context, ch chan context.Context, bunch BunchType) error {
	cfg := defaultConfig
	if Has(spinCtx, ConfigCtxKey()) {
		cfg = GetFrom[Config](spinCtx, ConfigCtxKey())
	}
	for {
		select {
		case <-spinCtx.Done():
			return spinCtx.Err()
		case ctx := <-ch:
			cfg.Worker(func() {
				defer func() {
					for _, f := range bunch.defers {
						ctx = f(ctx)
					}
				}()

				var err error
				for _, stick := range bunch.sticks {
					if stick.Ignore(ctx, err) {
						continue
					}
					ctx, err = stick.Handle(ctx, err)
				}
			})
		}
	}
}

func Bunch() BunchType {
	return BunchType{}
}

type BunchType struct {
	sticks []Stick
	defers []func(ctx context.Context) context.Context
}

func (b BunchType) I(stick Stick) BunchType {
	b.sticks = append(b.sticks, stick)
	return b
}

func (b BunchType) L(branch BunchType) BunchType {
	b.sticks = append(b.sticks, branch.sticks...)
	return b
}

func (b BunchType) Defer(fn func(ctx context.Context) context.Context) BunchType {
	b.defers = append(b.defers, fn)
	return b
}

type Stick interface {
	Ignore(ctx context.Context, err error) bool
	Handle(ctx context.Context, err error) (context.Context, error)
}

func Straw(fn func(ctx context.Context, err error) (context.Context, error)) Stick {
	return StrawType{fn: fn}
}

type StrawType struct {
	fn func(ctx context.Context, err error) (context.Context, error)
}

func (StrawType) Ignore(_ context.Context, _ error) bool {
	return false
}

func (s StrawType) Handle(ctx context.Context, err error) (context.Context, error) {
	return s.fn(ctx, err)
}
