package utils

import "net"

func GetLocalIP() string {
	ipAddr := "localhost"
	addrSlice, err := net.InterfaceAddrs()
	if err != nil {
		return ipAddr
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddr = ipnet.IP.String()
				break
			}
		}
	}
	return ipAddr
}
