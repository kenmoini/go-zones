package main

import (
	"strconv"
	"strings"

	"github.com/brotherpowers/ipsubnet"
)

// reverseIntSlice reverses an int slice, keeping the order of the elements
func reverseIntSlice(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// splitV6AddressIntoParts splits an address into its parts:
//  IP, CIDR, reversedNetworkPortion, reversedHostPortion
func splitV6AddressIntoParts(address string) (string, string, string, string) {

	addressArray := strings.Split(address, "/")
	addressPart := addressArray[0]

	// If the address is a full CIDR address, calculate the netmask
	if len(addressArray) > 1 {
		cidrPartStr := addressArray[1]
		//cidrPart, _ := strconv.Atoi(addressArray[1])

		return addressPart, cidrPartStr, "", ""
	}

	return addressPart, "", "", ""
}

// splitV4AddressIntoParts splits an address into its parts:
//  IP, CIDR, reversedNetworkPortion, reversedHostPortion
func splitV4AddressIntoParts(address string) (string, string, string, string) {

	addressArray := strings.Split(address, "/")
	addressPart := addressArray[0]

	// If the address is a full CIDR address, calculate the netmask
	if len(addressArray) > 1 {
		cidrPart, _ := strconv.Atoi(addressArray[1])
		sub := ipsubnet.SubnetCalculator(addressPart, cidrPart)

		networkPortion := sub.GetNetworkPortionQuards()
		hostPortion := sub.GetHostPortionQuards()
		var networkArray []int
		var hostArray []int

		var reversedNetworkArray []string
		var reversedHostArray []string

		var reversedNetworkPortion string
		var reversedHostPortion string

		//log.Printf("networkPortion: %v", networkPortion)
		//log.Printf("hostPortion: %v", hostPortion)

		if cidrPart <= 8 {
			networkArray = networkPortion[0:1]
			hostArray = hostPortion[1:]
		} else if cidrPart <= 16 {
			networkArray = networkPortion[0:2]
			hostArray = hostPortion[2:]
		} else if cidrPart <= 24 {
			networkArray = networkPortion[0:3]
			hostArray = hostPortion[3:]
		} else if cidrPart <= 32 {
			networkArray = networkPortion[0:4]
			hostArray = hostPortion[4:]
		}

		//log.Printf("networkArray: %v", networkArray)
		//log.Printf("hostArray: %v", hostArray)

		revNetArr := reverseIntSlice(networkArray)
		revHostArr := reverseIntSlice(hostArray)
		//log.Printf("networkArray reversed: %v", revNetArr)
		//log.Printf("hostArray reversed: %v", revHostArr)

		// Join the reversed network and host portions
		for _, v := range revNetArr {
			reversedNetworkArray = append(reversedNetworkArray, strconv.Itoa(v))
		}
		for _, v := range revHostArr {
			reversedHostArray = append(reversedHostArray, strconv.Itoa(v))
		}
		reversedHostPortion = strings.Join(reversedHostArray, ".")
		reversedNetworkPortion = strings.Join(reversedNetworkArray, ".")

		//log.Printf("reversedHostPortion: %v", reversedHostPortion)
		//log.Printf("reversedNetworkPortion: %v", reversedNetworkPortion)

		return addressPart, addressArray[1], reversedNetworkPortion, reversedHostPortion
	}

	return addressPart, "", "", ""
}

// GetPadded extends the MaxLengths struct with a function that takes a record component and returns a padded record component
func (ml MaxLengths) GetPadded(recordType string, recordKey string, recordValue string, recordTTL int) string {
	direction := "RIGHT"
	val := recordValue
	length := len(val)

	if recordValue == "" {
		val = strconv.Itoa(recordTTL)
		length = len(val)
	}

	if recordKey == "Value" {
		direction = "LEFT"
	}

	switch recordType {
	case "A":
		length = ml.A[recordKey]
	case "AAAA":
		length = ml.AAAA[recordKey]
	case "CNAME":
		length = ml.CNAME[recordKey]
	case "MX":
		length = ml.MX[recordKey]
	case "NS":
		length = ml.NS[recordKey]
	case "TXT":
		length = ml.TXT[recordKey]
	case "SRV":
		length = ml.SRV[recordKey]
	case "PTR":
		length = ml.PTR[recordKey]
	}

	pVal := StrPad(val, length, " ", direction)
	return pVal

}
