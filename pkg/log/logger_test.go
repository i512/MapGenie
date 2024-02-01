package log

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInfof(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	ctx = Ctx(ctx, Info, b)

	Infof(ctx, "%s", "processing")
	assert.Equal(t, "processing\n", b.String())
}

func TestSuppress(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	ctx = Ctx(ctx, Info, b)

	Suppress(ctx, Error, func(ctx context.Context) {
		Infof(ctx, "%s", "processing")
		assert.Equal(t, "", b.String())
		Errorf(ctx, "%s", "boom")
		assert.Equal(t, "processing\nboom\n", b.String())
		Debugf(ctx, "%s", "oof")
		assert.Equal(t, "processing\nboom\noof\n", b.String())
	})
}

func TestLogsAroundSuppress(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	ctx = Ctx(ctx, Info, b)

	Infof(ctx, "before")
	assert.Equal(t, "before\n", b.String())

	Suppress(ctx, Error, func(ctx context.Context) {
		Infof(ctx, "%s", "processing")
		assert.Equal(t, "before\n", b.String())
	})

	Warnf(ctx, "after")
	assert.Equal(t, "before\nafter\n", b.String())
}

//func TestNestedSuppress(t *testing.T) {
//	b := &bytes.Buffer{}
//	ctx := context.Background()
//	ctx = Ctx(ctx, Info, b)
//
//	Suppress(ctx, Error, func(ctx context.Context) {
//		Infof(ctx, "%s", "before 1")
//		assert.Equal(t, "", b.String())
//	})
//}
