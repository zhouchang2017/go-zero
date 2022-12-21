package logx

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestZapLogger_WithRequestId(t *testing.T) {
	var buf strings.Builder
	logger := NewTestLogger(&buf)
	SetGlobalLogger(logger)
	defer buf.Reset()

	logger.Infof("hello")
	assert.NotContains(t, buf.String(), "xxxxx1111.2222")

	l := logger.WithRequestId("xxxxx1111.2222")
	l.Infof("hello")
	assert.Contains(t, buf.String(), "xxxxx1111.2222")
}
