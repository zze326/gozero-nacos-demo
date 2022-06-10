package utils

import (
	"fmt"
	"net"
	"os"
	"strings"
)

/**
 * @Author: zze
 * @Date: 2022/6/2 16:21
 * @Desc: 主机操作
 */

// GetOutBoundIP
// desc: 获取本机的出网 IP
func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "223.5.5.5:53")
	if err != nil {
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}

func HostIdentity() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	ip, err := GetOutBoundIP()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", hostname, ip), nil
}
