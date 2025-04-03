package appctl

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// FindAvailablePorts 查找指定数量的可用端口
func findAvailablePorts(count, startPort, endPort int) ([]int, error) {
	var availablePorts []int

	for port := startPort; port <= endPort; port++ {
		if isPortAvailable(port) {
			availablePorts = append(availablePorts, port)
			if len(availablePorts) >= count {
				return availablePorts, nil
			}
		}
	}

	if len(availablePorts) < count {
		return nil, fmt.Errorf("端口范围%d-%d内只找到%d个可用端口,需要 %d 个",
			startPort, endPort, len(availablePorts), count)
	}

	return availablePorts, nil
}

// isPortAvailable 双重验证端口是否可用
func isPortAvailable(port int) bool {
	// 首先尝试监听
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	defer listener.Close()

	// 然后尝试连接
	conn, err := net.DialTimeout("tcp", "localhost:"+strconv.Itoa(port), 100*time.Millisecond)
	if err == nil {
		conn.Close()
		return false
	}

	return true
}
