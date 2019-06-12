
package util

import (
	"io/ioutil"
	"net"
	"net/http"
)

func GetExternalAddress() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if body, err1 := ioutil.ReadAll(resp.Body); err1 == nil {
		return string(body[:len(body)-1])
	}
	return ""
}

func GetInternalAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

