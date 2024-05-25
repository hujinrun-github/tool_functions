package address

import (
	"fmt"
	"net"
)

func GetLocalAddress() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		return []string{"127.0.0.1"}, err
	}

	addresses := []string{}
	// Iterate through the addresses and print them.
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Skip loopback and link-local addresses.
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() {
			continue
		}

		addresses = append(addresses, ip.String())
	}
	return addresses, nil
}

func GetLocalAddressByInterface(interfaceName string) ([]string, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		fmt.Println("Error:", err)
		return []string{"127.0.0.1"}, err
	}

	addrs, error := iface.Addrs()
	if error != nil {
		fmt.Println("Error:", error)
		return []string{"127.0.0.1"}, err
	}

	addresses := []string{}
	// Iterate through the addresses and print them.
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Skip loopback and link-local addresses.
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() {
			continue
		}

		addresses = append(addresses, ip.String())
	}
	return addresses, nil
}
