package log

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestInfof(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	ctx = Ctx(ctx, Info, b)

	Infof(ctx, "%s", "processing")
	assert.Equal(t, Color(IColor, "processing")+"\n", b.String())
}

func TestFoldUnfoldsOnError(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	ctx = Ctx(ctx, Info, b)

	foldedCtx := Fold(ctx, "fold title")
	Debugf(foldedCtx, "debug")
	assert.Equal(t, "", b.String())
	Infof(foldedCtx, "info")
	assert.Equal(t, "", b.String())
	Warnf(foldedCtx, "warn")
	assert.Equal(t, "", b.String())
	Errorf(foldedCtx, "boom!")
	assert.Equal(t, strings.Join([]string{
		Color(IColor, "fold title"),
		Color(DColor, "  debug"),
		Color(IColor, "  info"),
		Color(WColor, "  warn"),
		Color(EColor, "  boom!"),
	}, "\n")+"\n", b.String())

	b.Reset()
	Infof(foldedCtx, "info2")
	assert.Equal(t, Color(IColor, "  info2")+"\n", b.String())
}

func TestUnfoldsParents(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	ctx = Ctx(ctx, Info, b)

	Infof(ctx, "top")
	assert.Equal(t, Color(IColor, "top")+"\n", b.String())
	b.Reset()

	fold1Ctx := Fold(ctx, "fold1")
	Infof(fold1Ctx, "fold1info")

	foldOtherCtx := Fold(fold1Ctx, "hidden")
	Infof(foldOtherCtx, "hidden")

	fold2Ctx := Fold(fold1Ctx, "fold2")
	Infof(fold2Ctx, "fold2info")

	fold3Ctx := Fold(fold2Ctx, "fold3")
	Infof(fold3Ctx, "fold3info")

	assert.Equal(t, "", b.String())
	Errorf(fold3Ctx, "fold3error")

	expected := strings.Join([]string{
		Color(IColor, "fold1"),
		Color(IColor, "  fold1info"),
		Color(IColor, "  fold2"),
		Color(IColor, "    fold2info"),
		Color(IColor, "    fold3"),
		Color(IColor, "      fold3info"),
		Color(EColor, "      fold3error"),
	}, "\n") + "\n"

	assert.Equal(t, expected, b.String())
	b.Reset()

	Infof(fold1Ctx, "fold1info2")
	Infof(fold3Ctx, "fold3info2")
	Infof(fold2Ctx, "fold2info2")
	Infof(ctx, "top2")
	Infof(foldOtherCtx, "hidden")

	expected = strings.Join([]string{
		Color(IColor, "  fold1info2"),
		Color(IColor, "      fold3info2"),
		Color(IColor, "    fold2info2"),
		Color(IColor, "top2"),
	}, "\n") + "\n"
	assert.Equal(t, expected, b.String())
}
