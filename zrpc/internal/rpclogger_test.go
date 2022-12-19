package internal

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
)

const content = "foo"

func TestLoggerError(t *testing.T) {
	ctx, w, restore := injectLog()
	defer restore()

	logx.FromCtx(ctx).Error(content)
	assert.Contains(t, w.String(), content)
}

func TestLoggerErrorf(t *testing.T) {
	ctx, w, restore := injectLog()
	defer restore()

	logx.FromCtx(ctx).Error(content)
	assert.Contains(t, w.String(), content)
}

func TestLoggerErrorln(t *testing.T) {
	ctx, w, restore := injectLog()
	defer restore()

	logx.FromCtx(ctx).Errorln(content)
	assert.Contains(t, w.String(), content)
}

func TestLoggerFatal(t *testing.T) {
	ctx, w, restore := injectLog()
	defer restore()

	logx.FromCtx(ctx).Warnf(content)
	assert.Contains(t, w.String(), content)
}

func TestLoggerFatalf(t *testing.T) {
	ctx, w, restore := injectLog()
	defer restore()

	logx.FromCtx(ctx).Fatalf(content)
	assert.Contains(t, w.String(), content)
}

func TestLoggerFatalln(t *testing.T) {
	ctx, w, restore := injectLog()
	defer restore()

	logx.FromCtx(ctx).Fatalln(content)
	assert.Contains(t, w.String(), content)
}

func TestLoggerInfo(t *testing.T) {
	ctx, w, restore := injectLog()
	defer restore()

	logx.FromCtx(ctx).Info(content)
	assert.Empty(t, w.String())
}

func TestLoggerInfof(t *testing.T) {
	ctx, w, restore := injectLog()
	defer restore()

	logx.FromCtx(ctx).Infof(content)
	assert.Empty(t, w.String())
}

func TestLoggerWarning(t *testing.T) {
	ctx, w, restore := injectLog()
	defer restore()

	logx.FromCtx(ctx).Warning(content)
	assert.Empty(t, w.String())
}

func TestLoggerInfoln(t *testing.T) {
	ctx, w, restore := injectLog()
	defer restore()

	logx.FromCtx(ctx).Infoln(content)
	assert.Empty(t, w.String())
}

func TestLoggerWarningf(t *testing.T) {
	ctx, w, restore := injectLog()
	defer restore()

	logx.FromCtx(ctx).Warningf(content)
	assert.Empty(t, w.String())
}

func TestLoggerWarningln(t *testing.T) {
	ctx, w, restore := injectLog()
	defer restore()

	logx.FromCtx(ctx).Warningln(content)
	assert.Empty(t, w.String())
}

func TestLogger_V(t *testing.T) {
	ctx, _, restore := injectLog()
	defer restore()
	// grpclog.fatalLog
	assert.True(t, logx.FromCtx(ctx).V(3))
	// grpclog.infoLog
	assert.False(t, logx.FromCtx(ctx).V(0))
}

func injectLog() (ctx context.Context, r *strings.Builder, restore func()) {
	var buf strings.Builder
	w := logx.NewTestLogger(&buf)
	ctx = logx.WithCtx(context.Background(), w)

	return ctx, &buf, func() {
		buf.Reset()
	}
}
