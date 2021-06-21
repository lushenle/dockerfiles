package main

import (
	"net"
	"strings"
	"log"
	"fmt"
)

func GetLocalIPAddrs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	ips := make([]string, 0)
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}

func main() {
	if localIps, err := GetLocalIPAddrs(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(strings.Join(localIps, ", "))
	}

}
