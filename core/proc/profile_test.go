package proc

import (
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"testing"
)

func TestProfile(t *testing.T) {
	var buf strings.Builder
	logx.SetGlobalLogger(logx.NewTestLogger(&buf))

	profiler := StartProfile()
	// start again should not work
	assert.NotNil(t, StartProfile())
	profiler.Stop()
	// stop twice
	profiler.Stop()
	assert.True(t, strings.Contains(buf.String(), ".pprof"))
}
