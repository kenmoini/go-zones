package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func GenerateBindReverseZoneFiles(dnsServer *DNS, basePath string) (map[string][]string, error) {
	// Loop through the zones to calculate records and needs for reverse zone files
	var reverseZones = make(map[string][]PTRRecord)

	var reverseViewPair = make(map[string][]string)

	for _, zone := range dnsServer.Zones {
		var PTRRecords []PTRRecord
		var views []string

		// get the view(s) associated with this zone
		for _, view := range dnsServer.Views {
			if stringInSlice(zone.Name, view.IncludedZones) {
				views = append(views, view.Name)
			}
		}
		//log.Printf("debug: zone: %v", zone.Name)
		//log.Printf("debug: views: %v", views)

		// Check for defaults/overrides
		// Set default TTLs
		var zoneTTL int = defaultTTL
		if zone.DefaultTTL != 0 {
			zoneTTL = zone.DefaultTTL
		}

		// Setup serial number from the current unix time
		r_longTime := strconv.FormatInt(time.Now().UnixNano(), 10)
		r_shortTimeSerial := r_longTime[len(r_longTime)-9:]

		// Loop through the A Records - check to see if any are full CIDR addresses, if so we'll generate PTR records
		for _, record := range zone.Records.A {
			// Check for a TTL on the record, otherwise set the default
			var recordTTL int = zoneTTL
			if record.TTL != 0 {
				recordTTL = record.TTL
			}

			if strings.Contains(record.Value, "/") {
				_, _, r_networkPortion, r_hostPortion := splitV4AddressIntoParts(record.Value)
				//address, cidr, r_networkPortion, r_hostPortion := splitV4AddressIntoParts(record.Value)
				//log.Printf("address: %v", address)
				//log.Printf("cidr: %v", cidr)
				//log.Printf("r_hostPortion: %v", r_hostPortion)
				//log.Printf("r_networkPortion: %v", r_networkPortion)

				// Unless the record has NoPTR set then create a PTR record
				if !record.NoPTR {
					// Check to make sure this is not a wildcard record
					if !strings.Contains(record.Name, "*") {
						// Create a new PTRRecord variable with the reverse address
						recordValuePrefix := ""
						if record.Name != "@" {
							recordValuePrefix = record.Name + "."
						}
						PTRRecord := PTRRecord{
							Name:              r_hostPortion,
							Value:             recordValuePrefix + zone.Zone + ".",
							TTL:               recordTTL,
							TargetReverseZone: r_networkPortion + ".in-addr.arpa",
						}
						//log.Printf("PTRRecord: %v", PTRRecord)
						PTRRecords = append(PTRRecords, PTRRecord)

						// Loop through the assciated views and add the reverse zone to the list of reverse zones
						for _, view := range views {
							if !stringInSlice(r_networkPortion+".in-addr.arpa", reverseViewPair[view]) {
								reverseViewPair[view] = append(reverseViewPair[view], r_networkPortion+".in-addr.arpa")
							}
						}

					}
				}
			}

		}

		// If we created PTR records in this Zone from A or AAAA Records, then we'll need to create a reverse zone
		if len(PTRRecords) > 0 {

			var usedValues []string

			// Loop through the PTR Records and create a reverse zone
			for _, record := range PTRRecords {
				if !stringInSlice(record.Name, usedValues) {
					usedValues = append(usedValues, record.Name)

					PTRRecord := PTRRecord{
						Name:              record.Name,
						Value:             record.Value,
						TTL:               record.TTL,
						TargetReverseZone: record.TargetReverseZone,
					}

					reverseZones[record.TargetReverseZone] = append(reverseZones[record.TargetReverseZone], PTRRecord)
				}
			}

			//log.Printf("reverseZones: %v", reverseZones)
		}

		//=================================================
		// Build Reverse Zone Variables (if needed)
		if len(reverseZones) > 0 {

			// Loop through the reverse zones, set each up individually
			for reverseZone, records := range reverseZones {
				// Build the Forward Zone variable back up with our processed A Records
				newReverseZone := Zone{
					Name:             reverseZone,
					Zone:             reverseZone,
					PrimaryDNSServer: zone.PrimaryDNSServer,
					DefaultTTL:       zoneTTL,
					Records: Records{
						PTR: records,
					}}

				// Calculate the max lengths for the zone records
				maxLengths := calculateMaxRecordComponentLength(newReverseZone)

				PackagedZoneStructure := PackagedZone{
					Zone:                  newReverseZone,
					TTL:                   zoneTTL,
					SerialNumber:          r_shortTimeSerial,
					DefaultZoneSOARefresh: defaultZoneSOARefresh,
					DefaultZoneSOARetry:   defaultZoneSOARetry,
					DefaultZoneSOAExpire:  defaultZoneSOAExpire,
					DefaultZoneSOAMinTTL:  defaultZoneSOAMinTTL,
					Mode:                  "reverse",
					Path:                  basePath + "/zones/reverse." + reverseZone + ".zone",
					MaxLengths:            maxLengths}

				// Parse template
				t, err := template.New("revzones").Funcs(template.FuncMap{
					"ttlSwap": func(ttl int) int {
						if ttl == 0 {
							return zoneTTL
						}
						return ttl
					},
				}).Parse(bindZoneFileTemplate)
				check(err)
				// Create zone file
				f, err := os.Create(PackagedZoneStructure.Path)
				check(err)
				// Execute zone file templating
				err = t.Execute(f, PackagedZoneStructure)
				check(err)
				// Close and write file
				f.Close()
			}

			GenerateBindZoneReverseConfigFile(reverseZones, basePath)
		}

	}
	log.Printf("reverseViewPair: %v", reverseViewPair)
	return reverseViewPair, nil

}
