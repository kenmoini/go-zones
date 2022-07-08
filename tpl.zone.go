package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func GenerateBindZoneFiles(dnsServer *DNS, basePath string) (bool, error) {

	for _, zone := range dnsServer.Zones {

		var ARecords []ARecord
		var AAAARecords []AAAARecord

		// Check for defaults/overrides
		// Set default TTLs
		var zoneTTL int = defaultTTL
		if zone.DefaultTTL != 0 {
			zoneTTL = zone.DefaultTTL
		}

		//=================================================
		// Build Forward Zone Variables
		// Setup serial number from the current unix time
		longTime := strconv.FormatInt(time.Now().UnixNano(), 10)
		f_shortTimeSerial := longTime[len(longTime)-9:]

		// Loop through the A Records - check to see if any are full CIDR addresses, if so we'll generate PTR records
		for _, record := range zone.Records.A {
			// Check for a TTL on the record, otherwise set the default
			var recordTTL int = zoneTTL
			if record.TTL != 0 {
				recordTTL = record.TTL
			}
			if strings.Contains(record.Value, "/") {
				address, _, _, _ := splitV4AddressIntoParts(record.Value)

				// Create a new ARecord variable with just the address
				ARecord := ARecord{
					Name:  record.Name,
					Value: address,
					TTL:   recordTTL,
				}
				ARecords = append(ARecords, ARecord)

			} else {
				// This is just a plain old A record
				ARecord := ARecord{
					Name:  record.Name,
					Value: record.Value,
					TTL:   recordTTL,
				}
				ARecords = append(ARecords, ARecord)
			}
		}

		// Loop through the AAAA Records - check to see if any are full CIDR addresses, if so we'll generate PTR records
		for _, record := range zone.Records.AAAA {
			// Check for a TTL on the record, otherwise set the default
			var recordTTL int = zoneTTL
			if record.TTL != 0 {
				recordTTL = record.TTL
			}

			if strings.Contains(record.Value, "/") {
				address, _, _, _ := splitV6AddressIntoParts(record.Value)

				// Create a new AAAARecord variable with just the address
				AAAARecord := AAAARecord{
					Name:  record.Name,
					Value: address,
					TTL:   recordTTL,
				}
				AAAARecords = append(AAAARecords, AAAARecord)

			} else {
				// This is just a plain old AAAA record
				AAAARecord := AAAARecord{
					Name:  record.Name,
					Value: record.Value,
					TTL:   recordTTL,
				}
				AAAARecords = append(AAAARecords, AAAARecord)
			}
		}

		//=================================================
		// Build the Forward Zone variable back up with our processed A Records
		newForwardZone := Zone{
			Name:             zone.Name,
			Zone:             zone.Zone,
			PrimaryDNSServer: zone.PrimaryDNSServer,
			DefaultTTL:       zoneTTL,
			Records: Records{
				A:     ARecords,
				AAAA:  AAAARecords,
				CNAME: zone.Records.CNAME,
				MX:    zone.Records.MX,
				NS:    zone.Records.NS,
				TXT:   zone.Records.TXT,
				SRV:   zone.Records.SRV,
			}}

		// calculate the max lengths for the zone records
		maxLengths := calculateMaxRecordComponentLength(newForwardZone)

		PackagedZoneStructure := PackagedZone{
			Zone:                  newForwardZone,
			TTL:                   zoneTTL,
			SerialNumber:          f_shortTimeSerial,
			DefaultZoneSOARefresh: defaultZoneSOARefresh,
			DefaultZoneSOARetry:   defaultZoneSOARetry,
			DefaultZoneSOAExpire:  defaultZoneSOAExpire,
			DefaultZoneSOAMinTTL:  defaultZoneSOAMinTTL,
			Mode:                  "forward",
			Path:                  basePath + "/zones/fwd." + zone.Name + ".zone",
			MaxLengths:            maxLengths}

		// Parse template
		t, err := template.New("zones").Funcs(template.FuncMap{
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
		log.Println("Creating forward zone file: " + PackagedZoneStructure.Path)

		// Execute zone file templating
		err = t.Execute(f, PackagedZoneStructure)
		check(err)

		// Close and write file
		f.Close()

	}
	return true, nil
}

const bindZoneFileTemplate = `$ORIGIN {{ .Zone.Zone }}.
$TTL {{ .TTL }}

{{- $maxLengths := .MaxLengths }}

@ IN  SOA	{{ .Zone.PrimaryDNSServer }}. hostmaster.{{ .Zone.Zone }}. (
	{{ .SerialNumber }}
	{{ .DefaultZoneSOARefresh }}
	{{ .DefaultZoneSOARetry }}
	{{ .DefaultZoneSOAExpire }}
	{{ .DefaultZoneSOAMinTTL }} )

{{- with .Zone.Records.NS }}

; === NS Records ====================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "NS" "Anchor" .Anchor 0 }} {{ $.MaxLengths.GetPadded "NS" "TTL" "" (ttlSwap .TTL) }} IN NS {{ $.MaxLengths.GetPadded "NS" "Name" .Name 0 }}.{{ $.MaxLengths.GetPadded "NS" "Domain" .Domain 0 }}
{{- end }}
{{- end }}

{{- with .Zone.Records.MX }}

; === MX Records ====================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "MX" "Name" .Name 0 }} {{ $.MaxLengths.GetPadded "MX" "TTL" "" (ttlSwap .TTL) }} IN MX {{ $maxLengths.GetPadded "MX" "Priority" "" .Priority }} {{ $maxLengths.GetPadded "MX" "Value" .Value 0 }}
{{- end }}
{{- end }}

{{- with .Zone.Records.A }}

; === A Records =====================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "A" "Name" .Name 0 }} {{ $.MaxLengths.GetPadded "A" "TTL" "" (ttlSwap .TTL) }} IN A {{ $maxLengths.GetPadded "A" "Value" .Value 0 }}
{{- end }}
{{- end }}

{{- with .Zone.Records.AAAA }}

; === AAAA Records ==================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "AAAA" "Name" .Name 0 }} {{ $.MaxLengths.GetPadded "AAAA" "TTL" "" (ttlSwap .TTL) }} IN AAAA {{ $.MaxLengths.GetPadded "AAAA" "Value" .Value 0 }}
{{- end }}
{{- end }}

{{- with .Zone.Records.CNAME }}

; === CNAME Records =================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "CNAME" "Name" .Name 0 }} {{ $.MaxLengths.GetPadded "CNAME" "TTL" "" (ttlSwap .TTL) }} IN CNAME {{ $.MaxLengths.GetPadded "CNAME" "Value" .Value 0 }}
{{- end }}
{{- end }}

{{- with .Zone.Records.TXT }}

; === TXT Records ===================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "TXT" "Name" .Name 0 }} {{ $.MaxLengths.GetPadded "TXT" "TTL" "" (ttlSwap .TTL) }} IN TXT {{ $.MaxLengths.GetPadded "TXT" "Value" .Value 0 }}
{{- end }}
{{- end }}

{{- with .Zone.Records.SRV }}

; === SRV Records ===================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "TXT" "Name" .Name 0 }} {{ $.MaxLengths.GetPadded "TXT" "TTL" "" (ttlSwap .TTL) }} IN SRV {{ $.MaxLengths.GetPadded "SRV" "Priority" "" .Priority }} {{ $.MaxLengths.GetPadded "SRV" "Weight" "" .Weight }} {{ $.MaxLengths.GetPadded "SRV" "Port" "" .Port }} {{ $.MaxLengths.GetPadded "SRV" "Value" .Value 0 }}
{{- end }}
{{- end }}

{{- with .Zone.Records.PTR }}

; === PTR Records ===================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "PTR" "Name" .Name 0 }} IN PTR {{ $.MaxLengths.GetPadded "PTR" "Value" .Value 0 }}
{{- end }}
{{- end }}
`
