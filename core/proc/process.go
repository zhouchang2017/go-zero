package proc

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"time"
)

var (
	procName string
	pid      int
	mac      string
)

func init() {
	procName = filepath.Base(os.Args[0])
	pid = os.Getpid()
	mac = getMac()
}

//getMac 获取本机MAC地址
func getMac() string {
	interfaces, _ := net.Interfaces()
	for _, netInterface := range interfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		return macAddr
	}
	return randomMac()
}

func randomMac() string {
	var m [6]byte
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		macByte := rand.Intn(256)
		m[i] = byte(macByte)

		rand.Seed(int64(macByte))
	}
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", m[0], m[1], m[2], m[3], m[4], m[5])
}

// Pid returns pid of current process.
func Pid() int {
	return pid
}

// ProcessName returns the processname, same as the command name.
func ProcessName() string {
	return procName
}

//Mac returns the mac address
func Mac() string {
	return mac
}
