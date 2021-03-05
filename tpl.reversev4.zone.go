package main

import (
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/brotherpowers/ipsubnet"
)

// LoopThroughZonesForBindReverseV4ZonesFiles creates the zone files
// Validated with https://bind.jamiewood.io/
func LoopThroughZonesForBindReverseV4ZonesFiles(zones *Zones, basePath string) (bool, error) {
	for _, zone := range zones.Zones {
		if (zone.Name != "") && (zone.Network != "") {
			// Check for defaults/overrides
			var zoneTTL int = defaultTTL
			if zone.TTL != 0 {
				zoneTTL = zone.TTL
			}
			//=================================================
			// Create Reverse Zone files
			if IsReverse(zone.Name) == 0 {
				if zone.SubnetV4 != "" {
					cidr := strings.Split(zone.SubnetV4, "/")
					subnet := cidr[0]
					netblock, _ := strconv.Atoi(cidr[1])

					sub := ipsubnet.SubnetCalculator(subnet, netblock)
					networkPortion := sub.GetNetworkPortion()
					revAddr, err := reverseaddr(networkPortion)
					check(err)
					reverseZone := strings.ReplaceAll(revAddr, "0.", "")

					//=================================================
					// Create Forward Zone Files
					longTime := strconv.FormatInt(time.Now().UnixNano(), 10)
					shortTimeSerial := longTime[len(longTime)-9:]

					PackagedRevZoneStructure := PackagedReverseZone{
						Zone:                  zone,
						ReverseName:           reverseZone,
						TTL:                   zoneTTL,
						SerialNumber:          shortTimeSerial,
						DefaultZoneSOARefresh: defaultZoneSOARefresh,
						DefaultZoneSOARetry:   defaultZoneSOARetry,
						DefaultZoneSOAExpire:  defaultZoneSOAExpire,
						DefaultZoneSOAMinTTL:  defaultZoneSOAMinTTL,
						Mode:                  "reverse",
						Path:                  basePath + "/zones/" + reverseZone + "" + zone.Network + ".reverse.zone"}

					// Parse template
					t, err := template.New("revzones").Parse(bindReverseV4ZoneFileTemplate)
					check(err)
					// Create zone file
					f, err := os.Create(PackagedRevZoneStructure.Path)
					check(err)
					// Execute zone file templating
					err = t.Execute(f, PackagedRevZoneStructure)
					check(err)
					// Close and write file
					f.Close()
				}
				if zone.SubnetV6 != "" {

					cidrv6 := strings.Split(zone.SubnetV6, "/")
					subnetv6 := cidrv6[0]
					endSubnetv6 := subnetv6[len(subnetv6)-2:]
					var wholeSubnetv6 string

					if endSubnetv6 == "::" {
						wholeSubnetv6 = (subnetv6 + "0")
					} else {
						wholeSubnetv6 = (subnetv6)
					}

					shortReverse := reverse6Short(wholeSubnetv6)
					logStdOut("shortReverse: " + shortReverse)

				}
			}
		} else {
			return false, Stoerr("Name and/or network not defined!")
		}
	}
	return true, nil
}

// ReverseName just wraps a god damned string
type ReverseName string

// IPV4ToPortions Just wraps the functions needed by IPv4 Decoding for DNS
func IPV4ToPortions(network string, recordIP string) (subnet string, netblock string, networkPortion string, reverseAddress string, reverseZone string, networkPrefix string, revNetworkAddr string) {
	// Take whole network, split at the CIDR designation
	cidr := strings.Split(network, "/")
	subnet = cidr[0]
	netblock = cidr[1]

	netblockInt, err := strconv.Atoi(netblock)
	check(err)

	// Get the network portion of the IPV4 block
	sub := ipsubnet.SubnetCalculator(subnet, netblockInt)
	networkPortion = sub.GetNetworkPortion()

	revAddr, err := reverseaddr(recordIP)
	check(err)
	reverseZone = strings.TrimPrefix(revAddr, "0.")
	reverseZone = strings.TrimPrefix(reverseZone, "0.")
	reverseZone = strings.TrimPrefix(reverseZone, "0.")
	// Do this 3 times because...well...
	networkPrefix = strings.TrimSuffix(networkPortion, ".0")
	networkPrefix = strings.TrimSuffix(networkPrefix, ".0")
	networkPrefix = strings.TrimSuffix(networkPrefix, ".0")

	revNetworkAddr = strings.TrimPrefix(recordIP, networkPrefix)
	revNetworkAddr = strings.TrimPrefix(revNetworkAddr, ".")

	return subnet, netblock, networkPortion, revAddr, reverseZone, networkPrefix, revNetworkAddr
}

// RevValue takes a value and
func (r ARecord) RevValue(network string, recordIP string) string {
	_, _, _, _, _, _, revNetworkAddr := IPV4ToPortions(network, recordIP)
	return revNetworkAddr
}

const bindReverseV4ZoneFileTemplate = `$ORIGIN {{ .ReverseName }}
$TTL {{ .TTL }}

@ IN  SOA	{{ .Zone.PrimaryDNSServer }}. hostmaster.{{ .Zone.Name }}. (
	{{ .SerialNumber }}
	{{ .DefaultZoneSOARefresh }}
	{{ .DefaultZoneSOARetry }}
	{{ .DefaultZoneSOAExpire }}
	{{ .DefaultZoneSOAMinTTL }} )

{{ with .Zone.Records.NS }}{{ range . }}
{{ .Anchor }} {{ .TTL }} IN NS {{ .Name }}.{{ .Domain }}{{ end }}{{ end }}

{{ with .Zone.Records.A }}{{ range . }}
{{ .RevValue $.Zone.SubnetV4 .Value }} {{ .TTL }} IN PTR {{ .Name }}.{{ $.Zone.Name }}.{{ end }}{{ end }}
`
