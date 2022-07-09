package ipv6calc

import (
	"strings"
)

// PadIPv6Octet pads an IPv6 octet with zeros to make it a full octet
func PadIPv6Octet(octet string) string {
	return StrPad(octet, 4, "0", "LEFT")
}

// ReverseString reverses a string
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IPv6Address is the structure of an IPv6 address with the different parts composed
type IPv6Address struct {
	IP                string
	CIDR              string
	AddressParts      []string
	AddressString     string
	NetworkParts      []string
	NetworkPrefix     string
	HostParts         []string
	HostAddress       string
	PaddedHostAddress string
	ReverseZone       string
	ReverseRecord     string
}

func NewIPv6Address(address string) *IPv6Address {

	addressPart := ""
	cidrPart := ""
	networkPrefix := ""
	hostAddress := ""
	paddedHostAddress := ""
	addressParts := []string{"0000", "0000", "0000", "0000", "0000", "0000", "0000", "0000"}
	networkParts := []string{"", "", "", "", "", "", "", ""}
	hostParts := []string{"", "", "", "", "", "", "", ""}

	// Check if this is a full CIDR address
	if strings.Contains(address, "/") {
		// Split the address into its primary parts
		addressArray := strings.Split(address, "/")
		addressPart = addressArray[0]
		cidrPart = addressArray[1]
	}

	// Split the address up across a double colon
	addressPartsArr := strings.Split(addressPart, "::")
	// If there is only one part, then it is a full address, split it across the single colons
	if len(addressPartsArr) == 1 {
		addressSplit := strings.Split(addressPartsArr[0], ":")
		// Loop through the address parts and pad them with zeros then insert them into the addressParts array
		for i, part := range addressSplit {
			addressParts[i] = PadIPv6Octet(part)
		}

		// TODO: Calculate the networkParts, networkPrefix, hostParts, hostAddress, paddedHostAddress from the binary translation of the CIDR slash

	} else if len(addressPartsArr) == 2 {
		// If there are two parts, then the first part is the network portion and the second part is the host portion

		// Split the network portion across the single colons
		networkSplit := strings.Split(addressPartsArr[0], ":")
		// Loop through the network parts and pad them with zeros then insert them into the networkParts array
		for i, part := range networkSplit {
			networkParts[i] = PadIPv6Octet(part)
			addressParts[i] = PadIPv6Octet(part)
		}

		// Join the network parts together and set them to the networkPrefix
		networkPrefix = strings.TrimRight(strings.Join(networkParts, ":"), ":")

		// Split the host portion across the single colons
		hostSplit := strings.Split(addressPartsArr[1], ":")
		offsetPosition := 8 - (len(hostSplit) + len(networkSplit))
		// Loop through the host parts and pad them with zeros then insert them into the hostParts array
		for i, part := range hostSplit {
			paddedI := i + len(networkSplit) + offsetPosition
			hostParts[paddedI] = PadIPv6Octet(part)
			addressParts[paddedI] = PadIPv6Octet(part)
		}

		// Join the host parts together and set them to the hostAddress
		hostAddress = strings.TrimLeft(strings.Join(hostParts, ":"), ":")

		// Remove the network prefix from the overall address string
		paddedHostAddress = strings.TrimLeft(strings.Replace(strings.Join(addressParts, ":"), networkPrefix, "", 1), ":")

	}

	// Remove the colons from the networkPrefix
	networkPrefixR := strings.Replace(networkPrefix, ":", "", -1)
	// Remove the colons from the paddedHostAddress
	paddedHostAddressR := strings.Replace(paddedHostAddress, ":", "", -1)

	// Reverse the networkPrefixR string
	revNetwork := strings.Join(strings.Split(ReverseString(networkPrefixR), ""), ".")
	// Reverse the paddedHostAddressR string
	revPaddedHost := strings.Join(strings.Split(ReverseString(paddedHostAddressR), ""), ".")

	ipv6Address := &IPv6Address{
		IP:                addressPart,
		CIDR:              cidrPart,
		AddressParts:      addressParts,
		AddressString:     strings.Join(addressParts, ":"),
		NetworkParts:      networkParts,
		NetworkPrefix:     networkPrefix,
		HostParts:         hostParts,
		HostAddress:       hostAddress,
		PaddedHostAddress: paddedHostAddress,
		ReverseZone:       revNetwork + ".ip6.arpa",
		ReverseRecord:     revPaddedHost,
	}

	return ipv6Address
}

// SplitV6AddressIntoParts splits an address into its parts:
//  IP, CIDR, reversedNetworkPortion, reversedHostPortion
func SplitV6AddressIntoParts(address string) (string, string, string, string) {

	computedAddress := NewIPv6Address(address)
	//log.Printf("----------------------------------")
	//log.Printf("computedAddress.IP: %v", computedAddress.IP)
	//log.Printf("computedAddress.CIDR: %v", computedAddress.CIDR)
	//log.Printf("computedAddress.AddressParts (%v): %v", len(computedAddress.AddressParts), computedAddress.AddressParts)
	//log.Printf("computedAddress.AddressString: (%v): %v", len(computedAddress.AddressString), computedAddress.AddressString)
	//log.Printf("computedAddress.NetworkParts (%v): %v", len(computedAddress.NetworkParts), computedAddress.NetworkParts)
	//log.Printf("computedAddress.NetworkPrefix (%v): %v", len(computedAddress.NetworkPrefix), computedAddress.NetworkPrefix)
	//log.Printf("computedAddress.HostParts (%v): %v", len(computedAddress.HostParts), computedAddress.HostParts)
	//log.Printf("computedAddress.HostAddress (%v): %v", len(computedAddress.HostAddress), computedAddress.HostAddress)
	//log.Printf("computedAddress.PaddedHostAddress (%v): %v", len(computedAddress.PaddedHostAddress), computedAddress.PaddedHostAddress)
	//log.Printf("computedAddress.ReverseZone (%v): %v", len(computedAddress.ReverseZone), computedAddress.ReverseZone)
	//log.Printf("computedAddress.ReverseRecord (%v): %v", len(computedAddress.ReverseRecord), computedAddress.ReverseRecord)

	return computedAddress.IP, computedAddress.CIDR, computedAddress.ReverseZone, computedAddress.ReverseRecord
}
