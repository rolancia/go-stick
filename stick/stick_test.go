package stick_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rolancia/go-stick/stick"
)

func TestStick(t *testing.T) {
	bunch := stick.Bunch().
		Defer(func(ctx context.Context) context.Context {
			if v := recover(); v != nil {
				t.Error(v)
			}
			return ctx
		}).
		I(Add{1}).
		I(Add{2}).
		I(Add{3}).
		L(stick.Bunch().
			I(Add{4}).
			I(Add{5}).
			I(Add{6})).
		I(Add{7}).
		I(stick.Straw(func(ctx context.Context) context.Context {
			v := stick.GetFrom[int](ctx, valCtxKey(""))
			t.Log(v)
			assert.Equal(t, (7*8)/2, v)
			return ctx
		}))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ch := make(chan context.Context)
	defer close(ch)

	go func() {
		_ = stick.Spin(ctx, ch, bunch)
	}()
	go func() {
		ctx = stick.With(ctx, valCtxKey(""), 0)
		for {
			ch <- ctx
			time.Sleep(100 * time.Millisecond)
		}
	}()

	<-ctx.Done()
}

type valCtxKey string

var _ stick.Stick = &Add{}

type Add struct {
	Val int
}

func (s Add) Ignore(ctx context.Context) bool {
	return !stick.Has(ctx, valCtxKey(""))
}

func (s Add) Handle(ctx context.Context) context.Context {
	v := stick.GetFrom[int](ctx, valCtxKey(""))
	v += s.Val
	return stick.With(ctx, valCtxKey(""), v)
}
