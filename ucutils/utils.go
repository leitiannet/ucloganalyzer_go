package ucutils

import (
	"fmt"
	"net"
	"time"
)

func GetLocalIp() []string {
	localIps := make([]string, 0)
	interAddrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, addr := range interAddrs {
			if n, ok := addr.(*net.IPNet); ok && !n.IP.IsLoopback() {
				if n.IP.To4() != nil {
					localIps = append(localIps, n.IP.String())
				}
			}
		}
	}
	return localIps
}

func FormatTimestamp(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 03:04:05 PM")
}

func Timestamp() int64 {
	return time.Now().Unix()
}

func StrInSlice(val string, arr []string) bool {
	for _, v := range arr {
		if val == v {
			return true
		}
	}
	return false
}
