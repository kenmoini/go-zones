package main

import (
	"net"
	"net/http"
	"strings"

	"github.com/miekg/dns"
)

/*************************************************************************************
* IP Resolution
*************************************************************************************/

// ReadUserIP gets the requesting client's IP so you can do a reverse DNS lookup
func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

// ReadUserIPNoPort gets the requesting client's IP without the port so you can do a reverse DNS lookup
func ReadUserIPNoPort(r *http.Request) string {
	IPAddress := ReadUserIP(r)

	NoPort := strings.Split(IPAddress, ":")
	if len(NoPort) > 0 {
		NoPort = NoPort[:len(NoPort)-1]
	}
	JoinedAddress := strings.Join(NoPort[:], ":")
	return JoinedAddress
}

/*************************************************************************************
* Forward/Reverse DNS Translation
*************************************************************************************/

// ExtractAddressFromReverse turns a standard PTR reverse record name
// into an IP address. This works for ipv4 or ipv6.
//
// 54.119.58.176.in-addr.arpa. becomes 176.58.119.54. If the conversion
// fails the empty string is returned.
func ExtractAddressFromReverse(reverseName string) string {
	search := ""

	f := reverse

	switch {
	case strings.HasSuffix(reverseName, IP4arpa):
		search = strings.TrimSuffix(reverseName, IP4arpa)
	case strings.HasSuffix(reverseName, IP6arpa):
		search = strings.TrimSuffix(reverseName, IP6arpa)
		f = reverse6
	default:
		return ""
	}

	// Reverse the segments and then combine them.
	return f(strings.Split(search, "."))
}

// IsReverse returns 0 is name is not in a reverse zone. Anything > 0 indicates
// name is in a reverse zone. The returned integer will be 1 for in-addr.arpa. (IPv4)
// and 2 for ip6.arpa. (IPv6).
func IsReverse(name string) int {
	if strings.HasSuffix(name, IP4arpa) {
		return 1
	}
	if strings.HasSuffix(name, IP6arpa) {
		return 2
	}
	return 0
}

// reverseaddr returns the in-addr.arpa. or ip6.arpa. hostname of the IP
// address addr suitable for rDNS (PTR) record lookup or an error if it fails
// to parse the IP address.
func reverseaddr(addr string) (arpa string, err error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "", &net.DNSError{Err: "unrecognized address", Name: addr}
	}
	if ip.To4() != nil {
		return uitoa(uint(ip[15])) + "." + uitoa(uint(ip[14])) + "." + uitoa(uint(ip[13])) + "." + uitoa(uint(ip[12])) + ".in-addr.arpa.", nil
	}
	// Must be IPv6
	buf := make([]byte, 0, len(ip)*4+len("ip6.arpa."))
	// Add it, in reverse, to the buffer
	for i := len(ip) - 1; i >= 0; i-- {
		v := ip[i]
		buf = append(buf, hexDigit[v&0xF],
			'.',
			hexDigit[v>>4],
			'.')
	}
	// Append "ip6.arpa." and return (buf already has the final .)
	buf = append(buf, "ip6.arpa."...)
	return string(buf), nil
}

// reverse will take an IPv4 address and creates the reverse DNS equivilent value
func reverse(slice []string) string {
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - i - 1
		slice[i], slice[j] = slice[j], slice[i]
	}
	ip := net.ParseIP(strings.Join(slice, ".")).To4()
	if ip == nil {
		return ""
	}
	return ip.String()
}

// reverse6 reverse the segments and combine them according to RFC3596:
// b.a.9.8.7.6.5.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2
// is reversed to 2001:db8::567:89ab
func reverse6(slice []string) string {
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - i - 1
		slice[i], slice[j] = slice[j], slice[i]
	}
	slice6 := []string{}
	for i := 0; i < len(slice)/4; i++ {
		slice6 = append(slice6, strings.Join(slice[i*4:i*4+4], ""))
	}
	ip := net.ParseIP(strings.Join(slice6, ":")).To16()
	if ip == nil {
		return ""
	}
	return ip.String()
}

// reverseSlice flips a string slice around
func reverseSlice(input []string) []string {
	if len(input) == 0 {
		return input
	}
	return append(reverseSlice(input[1:]), input[0])
}

// reverse6Short provides a shorter reverse IPV6 addres
func reverse6Short(ip string) string {
	cidrv6 := strings.Split(ip, "/")
	networkv6 := strings.Split(cidrv6[0], "::")
	noColons := strings.ReplaceAll(networkv6[0], ":", "")
	netChars := strings.Split(noColons, "")
	revNetChars := reverseSlice(netChars)
	revNet := (strings.Join(revNetChars, ".") + IP6arpa)

	revAddrv6, err := dns.ReverseAddr(cidrv6[0])
	check(err)

	extraNodeJunk := strings.ReplaceAll(revAddrv6, revNet, "")

	finalNetworkTarget := strings.ReplaceAll(revAddrv6, extraNodeJunk, "")

	return finalNetworkTarget
}

const (
	// IP4arpa is the reverse tree suffix for v4 IP addresses.
	IP4arpa = ".in-addr.arpa."
	// IP6arpa is the reverse tree suffix for v6 IP addresses.
	IP6arpa = ".ip6.arpa."
)
