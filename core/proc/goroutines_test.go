package proc

import (
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"testing"
)

func TestDumpGoroutines(t *testing.T) {
	var buf strings.Builder
	logx.SetGlobalLogger(logx.NewTestLogger(&buf))

	dumpGoroutines()
	assert.True(t, strings.Contains(buf.String(), ".dump"))
}
