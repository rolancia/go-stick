package stick

import (
	"context"
)

func Spin(spinCtx context.Context, ch chan context.Context, bunch BunchType) error {
	for {
		select {
		case <-spinCtx.Done():
			return spinCtx.Err()
		case ctx := <-ch:
			go func() {
				defer func() {
					for _, f := range bunch.defers {
						ctx = f(ctx)
					}
				}()

				for _, stick := range bunch.sticks {
					if stick.Ignore(ctx) {
						continue
					}
					ctx = stick.Handle(ctx)
				}
			}()
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
	Ignore(ctx context.Context) bool
	Handle(ctx context.Context) context.Context
}

func Straw(fn func(ctx context.Context) context.Context) Stick {
	return StrawType{fn: fn}
}

type StrawType struct {
	fn func(ctx context.Context) context.Context
}

func (StrawType) Ignore(_ context.Context) bool {
	return false
}

func (s StrawType) Handle(ctx context.Context) context.Context {
	return s.fn(ctx)
}
