package libutil

import (
	"crypto/rand"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// generate mac address
func GenerateMacAddress() string {
	mac := make([]byte, 6)
	rand.Read(mac)
	mac[0] = (mac[0] | 2) & 0xfe
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

// check ip in network
func CheckIpInNetwork(ip, network string) bool {
	_, ipNet, err := net.ParseCIDR(network)
	if err != nil {
		return false
	}
	return ipNet.Contains(net.ParseIP(ip))
}

// generate mac address by ip address
func GenerateMacAddressByIp(ipWithMask string) (string, error) {
	ip := strings.Split(ipWithMask, "/")[0]
	prefixLength, err := strconv.Atoi(strings.Split(ipWithMask, "/")[1])
	if err != nil {
		return "", err
	}

	// split ip
	ipSplit := strings.Split(ip, ".")

	hexMacSlice := []int{}
	hexMacSlice = append(hexMacSlice, 0x52)
	hexMacSlice = append(hexMacSlice, prefixLength)

	for _, ipSplitItem := range ipSplit {
		ipSplitItemInt, err := strconv.Atoi(ipSplitItem)
		if err != nil {
			return "", err
		}

		hexMacSlice = append(hexMacSlice, ipSplitItemInt)
	}

	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		hexMacSlice[0], hexMacSlice[1], hexMacSlice[2], hexMacSlice[3], hexMacSlice[4], hexMacSlice[5]), nil
}
