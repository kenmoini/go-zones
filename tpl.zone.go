package main

import (
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/brotherpowers/ipsubnet"
)

func GenerateBindZoneFiles(dns *DNS, basePath string) (bool, error) {
	// Loop through the zones to calculate records and needs for reverse zone files
	for _, zone := range dns.Zones {

		var ARecords []ARecord
		//var PTRRecords []PTRRecord

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
				//address, cidr, netmask := splitAddressIntoParts(record.Value)
				address, _, _ := splitAddressIntoParts(record.Value)

				//sub := ipsubnet.SubnetCalculator(address, cidr)

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

		//=================================================
		// Build Reverse Zone Variables (if needed)
		// Setup serial number from the current unix time
		//r_longTime := strconv.FormatInt(time.Now().UnixNano(), 10)
		//r_shortTimeSerial := r_longTime[len(r_longTime)-9:]

		//=================================================
		// Build the Forward Zone variable back up with our processed A Records
		newForwardZone := Zone{
			Name:             zone.Name,
			Zone:             zone.Zone,
			PrimaryDNSServer: zone.PrimaryDNSServer,
			DefaultTTL:       zoneTTL,
			Records: Records{
				A:     ARecords,
				AAAA:  zone.Records.AAAA,
				CNAME: zone.Records.CNAME,
				MX:    zone.Records.MX,
				NS:    zone.Records.NS,
				TXT:   zone.Records.TXT,
				SRV:   zone.Records.SRV,
			}}

		//=================================================
		// Calculate the forward zones lengths
		var maxLengths MaxLengths
		var AMap = map[string]int{"Name": 0, "Value": 0, "TTL": 0}
		var AAAAMap = map[string]int{"Name": 0, "Value": 0, "TTL": 0}
		var CNAMEMap = map[string]int{"Name": 0, "Value": 0, "TTL": 0}
		var TXTMap = map[string]int{"Name": 0, "Value": 0, "TTL": 0}
		var NSMap = map[string]int{"Name": 0, "Value": 0, "TTL": 0, "Anchor": 0}
		var MXMap = map[string]int{"Name": 0, "Value": 0, "TTL": 0, "Priority": 0}
		var SRVMap = map[string]int{"Name": 0, "Value": 0, "TTL": 0, "Priority": 0, "Weight": 0, "Port": 0}

		// Loop through each NS record type and set the longest record component lengths to the MaxLengths record string map
		for _, record := range zone.Records.NS {
			if len(record.Name) > NSMap["Name"] {
				NSMap["Name"] = len(record.Name)
			}
			if len(strconv.Itoa(record.TTL)) > NSMap["TTL"] {
				NSMap["TTL"] = len(strconv.Itoa(record.TTL))
			}
			if len(record.Anchor) > NSMap["Anchor"] {
				NSMap["Anchor"] = len(record.Anchor)
			}
			if len(record.Domain) > NSMap["Domain"] {
				NSMap["Domain"] = len(record.Domain)
			}
		}

		// Loop through each A record type and set the longest record component lengths to the MaxLengths record string map
		for _, record := range zone.Records.A {
			if len(record.Name) > AMap["Name"] {
				AMap["Name"] = len(record.Name)
			}
			if len(strconv.Itoa(record.TTL)) > AMap["TTL"] {
				AMap["TTL"] = len(strconv.Itoa(record.TTL))
			}
			if len(record.Value) > AMap["Value"] {
				AMap["Value"] = len(record.Value)
			}
		}

		// Loop through each AAAA record type and set the longest record component lengths to the MaxLengths record string map
		for _, record := range zone.Records.AAAA {
			if len(record.Name) > AAAAMap["Name"] {
				AAAAMap["Name"] = len(record.Name)
			}
			if len(strconv.Itoa(record.TTL)) > AAAAMap["TTL"] {
				AAAAMap["TTL"] = len(strconv.Itoa(record.TTL))
			}
			if len(record.Value) > AAAAMap["Value"] {
				AAAAMap["Value"] = len(record.Value)
			}
		}

		// Loop through each CNAME record type and set the longest record component lengths to the MaxLengths record string map
		for _, record := range zone.Records.CNAME {
			if len(record.Name) > CNAMEMap["Name"] {
				CNAMEMap["Name"] = len(record.Name)
			}
			if len(strconv.Itoa(record.TTL)) > CNAMEMap["TTL"] {
				CNAMEMap["TTL"] = len(strconv.Itoa(record.TTL))
			}
			if len(record.Value) > CNAMEMap["Value"] {
				CNAMEMap["Value"] = len(record.Value)
			}
		}

		// Loop through each MX record type and set the longest record component lengths to the MaxLengths record string map
		for _, record := range zone.Records.MX {
			if len(record.Name) > MXMap["Name"] {
				MXMap["Name"] = len(record.Name)
			}
			if len(strconv.Itoa(record.TTL)) > MXMap["TTL"] {
				MXMap["TTL"] = len(strconv.Itoa(record.TTL))
			}
			if len(strconv.Itoa(record.Priority)) > MXMap["Priority"] {
				MXMap["Priority"] = len(strconv.Itoa(record.Priority))
			}
			if len(record.Value) > MXMap["Value"] {
				MXMap["Value"] = len(record.Value)
			}
		}

		// Loop through each SRV record type and set the longest record component lengths to the MaxLengths record string map
		for _, record := range zone.Records.SRV {
			if len(record.Name) > SRVMap["Name"] {
				SRVMap["Name"] = len(record.Name)
			}
			if len(strconv.Itoa(record.TTL)) > SRVMap["TTL"] {
				SRVMap["TTL"] = len(strconv.Itoa(record.TTL))
			}
			if len(strconv.Itoa(record.Priority)) > SRVMap["Priority"] {
				SRVMap["Priority"] = len(strconv.Itoa(record.Priority))
			}
			if len(strconv.Itoa(record.Weight)) > SRVMap["Weight"] {
				SRVMap["Weight"] = len(strconv.Itoa(record.Weight))
			}
			if len(strconv.Itoa(record.Port)) > SRVMap["Port"] {
				SRVMap["Port"] = len(strconv.Itoa(record.Port))
			}
			if len(record.Value) > SRVMap["Value"] {
				SRVMap["Value"] = len(record.Value)
			}
		}

		// Loop through each TXT record type and set the longest record component lengths to the MaxLengths record string map
		for _, record := range zone.Records.TXT {
			if len(record.Name) > TXTMap["Name"] {
				TXTMap["Name"] = len(record.Name)
			}
			if len(strconv.Itoa(record.TTL)) > TXTMap["TTL"] {
				TXTMap["TTL"] = len(strconv.Itoa(record.TTL))
			}
			if len(record.Value) > TXTMap["Value"] {
				TXTMap["Value"] = len(record.Value)
			}
		}

		maxLengths.NS = NSMap
		maxLengths.A = AMap
		maxLengths.AAAA = AAAAMap
		maxLengths.CNAME = CNAMEMap
		maxLengths.MX = MXMap
		maxLengths.SRV = SRVMap
		maxLengths.TXT = TXTMap

		// Loop through each A record and set the

		PackagedZoneStructure := PackagedZone{
			Zone:                  newForwardZone,
			TTL:                   zoneTTL,
			SerialNumber:          f_shortTimeSerial,
			DefaultZoneSOARefresh: defaultZoneSOARefresh,
			DefaultZoneSOARetry:   defaultZoneSOARetry,
			DefaultZoneSOAExpire:  defaultZoneSOAExpire,
			DefaultZoneSOAMinTTL:  defaultZoneSOAMinTTL,
			Mode:                  "forward",
			Path:                  basePath + "/zones/forward." + zone.Name + ".zone",
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
		// Execute zone file templating
		err = t.Execute(f, PackagedZoneStructure)
		check(err)
		// Close and write file
		f.Close()

	}
	return true, nil
}

// splitAddressIntoParts splits an address into its parts, IP, CIDR, and reversed host portion
func splitAddressIntoParts(address string) (string, string, string) {
	cidr := strings.Split(address, "/")
	addressPart := cidr[0]

	// If the address is a full CIDR address, calculate the netmask
	if len(cidr) > 1 {
		cidrPart, _ := strconv.Atoi(cidr[1])
		sub := ipsubnet.SubnetCalculator(addressPart, cidrPart)

		return addressPart, cidr[1], sub.GetSubnetMask()
	}

	return addressPart, "", ""
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
	}

	pVal := StrPad(val, length, " ", direction)
	return pVal

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

{{ with .Zone.Records.NS -}}
# === NS Records ====================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "NS" "Anchor" .Anchor 0 }} {{ $.MaxLengths.GetPadded "NS" "TTL" "" (ttlSwap .TTL) }} IN NS {{ $.MaxLengths.GetPadded "NS" "Name" .Name 0 }}.{{ $.MaxLengths.GetPadded "NS" "Domain" .Domain 0 }}
{{- end }}
{{- end }}

{{ with .Zone.Records.MX -}}
# === MX Records ====================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "MX" "Name" .Name 0 }} {{ $.MaxLengths.GetPadded "MX" "TTL" "" (ttlSwap .TTL) }} IN MX {{ $maxLengths.GetPadded "MX" "Priority" "" .Priority }} {{ $maxLengths.GetPadded "MX" "Value" .Value 0 }}
{{- end }}
{{- end }}

{{ with .Zone.Records.A -}}
# === A Records =====================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "A" "Name" .Name 0 }} {{ $.MaxLengths.GetPadded "A" "TTL" "" (ttlSwap .TTL) }} IN A {{ $maxLengths.GetPadded "A" "Value" .Value 0 }}
{{- end }}
{{- end }}

{{ with .Zone.Records.AAAA -}}
# === AAAA Records ==================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "AAAA" "Name" .Name 0 }} {{ $.MaxLengths.GetPadded "AAAA" "TTL" "" (ttlSwap .TTL) }} IN AAAA {{ $.MaxLengths.GetPadded "AAAA" "Value" .Value 0 }}
{{- end }}
{{- end }}

{{ with .Zone.Records.CNAME -}}
# === CNAME Records =================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "CNAME" "Name" .Name 0 }} {{ $.MaxLengths.GetPadded "CNAME" "TTL" "" (ttlSwap .TTL) }} IN CNAME {{ $.MaxLengths.GetPadded "CNAME" "Value" .Value 0 }}
{{- end }}
{{- end }}

{{ with .Zone.Records.TXT -}}
# === TXT Records ===================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "TXT" "Name" .Name 0 }} {{ $.MaxLengths.GetPadded "TXT" "TTL" "" (ttlSwap .TTL) }} IN TXT {{ $.MaxLengths.GetPadded "TXT" "Value" .Value 0 }}
{{- end }}
{{- end }}

{{ with .Zone.Records.SRV -}}
# === SRV Records ===================================================
{{- range . }}
{{ $.MaxLengths.GetPadded "TXT" "Name" .Name 0 }} {{ $.MaxLengths.GetPadded "TXT" "TTL" "" (ttlSwap .TTL) }} IN SRV {{ $.MaxLengths.GetPadded "SRV" "Priority" "" .Priority }} {{ $.MaxLengths.GetPadded "SRV" "Weight" "" .Weight }} {{ $.MaxLengths.GetPadded "SRV" "Port" "" .Port }} {{ $.MaxLengths.GetPadded "SRV" "Value" .Value 0 }}
{{- end }}
{{- end }}

{{ with .Zone.Records.PTR -}}
# === PTR Records ===================================================
{{- range . }}

{{- end }}
{{- end }}
`
