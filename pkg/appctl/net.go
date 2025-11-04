package appctl

//
//// FindAvailablePorts 查找指定数量的可用端口
//func findAvailablePorts(count, startPort, endPort int) ([]int, error) {
//	var availablePorts []int
//
//	for port := startPort; port <= endPort; port++ {
//		if isPortAvailable(port) {
//			availablePorts = append(availablePorts, port)
//			if len(availablePorts) >= count {
//				return availablePorts, nil
//			}
//		}
//	}
//
//	if len(availablePorts) < count {
//		return nil, fmt.Errorf("端口范围%d-%d内只找到%d个可用端口,需要 %d 个",
//			startPort, endPort, len(availablePorts), count)
//	}
//
//	return availablePorts, nil
//}
//
//// isPortAvailable 双重验证端口是否可用
//func isPortAvailable(port int) bool {
//	// 双重检查机制
//	for i := 0; i < 2; i++ { // 检查两次减少竞态条件风险
//		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
//		if err != nil {
//			return false
//		}
//		_ = ln.Close()
//
//		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 500*time.Millisecond)
//		if err != nil {
//			return true // 连接失败说明端口可能可用
//		}
//		if conn != nil {
//			_ = conn.Close()
//			return false
//		}
//
//		time.Sleep(10 * time.Millisecond) // 短暂等待
//	}
//	return false
//}
