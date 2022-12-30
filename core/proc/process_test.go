package proc

import (
	"github.com/davecgh/go-spew/spew"
	"math/big"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessName(t *testing.T) {
	assert.True(t, len(ProcessName()) > 0)
}

func TestPid(t *testing.T) {
	assert.True(t, Pid() > 0)
}

func TestGetMac(t *testing.T) {
	s := Mac()
	parseMAC, _ := net.ParseMAC(s)
	spew.Dump(parseMAC[0])
	i := big.NewInt(0).SetBytes(parseMAC).Int64()
	t.Logf("mac = %s %d\n", s, i)
	t.Logf("mac = %02x:%02x:%02x:%02x:%02x:%02x",
		byte(i>>40), byte(i>>32), byte(i>>24), byte(i>>16), byte(i>>8), byte(i),
	)
	assert.True(t, len(s) > 0)
}
