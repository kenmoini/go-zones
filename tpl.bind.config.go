package main

import (
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/brotherpowers/ipsubnet"
)

// BindZoneConfig will setup the bind config for zones
type BindZoneConfig struct {
	Network string // Network type (internal/external/etc)
	Name    string // what zone is being served, example.com
	Path    string // Path to the Zones file
	Mode    string // Mode is forward or reverse zone
}

// LoopThroughZonesForBindConfig creates the zone files
func LoopThroughZonesForBindConfig(zones *Zones, basePath string) (bool, error) {
	for _, zone := range zones.Zones {
		if (zone.Name != "") && (zone.Network != "") {
			//=================================================
			// Create Forward Zone Files
			BindConfigStructure := BindZoneConfig{
				Name:    zone.Name,
				Mode:    "forward",
				Network: zone.Network,
				Path:    basePath + "/zones/" + zone.Name + "." + zone.Network + ".forward.zone"}

			// Parse template
			t, err := template.New("zones").Parse(bindZoneConfigTemplate)
			check(err)
			// Create zone file
			f, err := os.Create(BindConfigStructure.Path)
			check(err)
			// Execute zone file templating
			err = t.Execute(f, BindConfigStructure)
			check(err)
			// Close and write file
			f.Close()

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

					ReverseBindConfigStructure := BindZoneConfig{
						Name:    reverseZone,
						Mode:    "reverse",
						Network: zone.Network,
						Path:    basePath + "/zones/" + reverseZone + zone.Network + ".reverse.zone"}

					// Parse template
					t, err := template.New("zones").Parse(bindZoneConfigTemplate)
					check(err)
					// Create zone file
					f, err := os.Create(ReverseBindConfigStructure.Path)
					check(err)
					// Execute zone file templating
					err = t.Execute(f, ReverseBindConfigStructure)
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

					ReverseBindConfigStructure := BindZoneConfig{
						Name:    shortReverse,
						Mode:    "reverse",
						Network: zone.Network,
						Path:    basePath + "/zones/" + shortReverse + zone.Network + ".reverse.zone"}

					// Parse template
					t, err := template.New("zones").Parse(bindZoneConfigTemplate)
					check(err)
					// Create zone file
					f, err := os.Create(ReverseBindConfigStructure.Path)
					check(err)
					// Execute zone file templating
					err = t.Execute(f, ReverseBindConfigStructure)
					check(err)
					// Close and write file
					f.Close()
				}
			}
		} else {
			logStdOut("Name and/or network not defined!")
		}
	}

	return true, nil
}

const bindZoneConfigTemplate = `zone "{{ .Name }}" {
	type master;
	file "{{ .Path }}";
};`
