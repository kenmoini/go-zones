package main

import (
	"os"
	"strconv"
	"text/template"
	"time"
)

// LoopThroughZonesForBindZonesFiles creates the zone files
// Validated with https://bind.jamiewood.io/
func LoopThroughZonesForBindZonesFiles(zones *Zones, basePath string) (bool, error) {
	for _, zone := range zones.Zones {
		if (zone.Name != "") && (zone.Network != "") {
			// Check for defaults/overrides
			var zoneTTL int = defaultTTL
			if zone.TTL != 0 {
				zoneTTL = zone.TTL
			}
			//=================================================
			// Create Forward Zone Files
			longTime := strconv.FormatInt(time.Now().UnixNano(), 10)
			shortTimeSerial := longTime[len(longTime)-9:]

			PackagedZoneStructure := PackagedZone{
				Zone:                  zone,
				TTL:                   zoneTTL,
				SerialNumber:          shortTimeSerial,
				DefaultZoneSOARefresh: defaultZoneSOARefresh,
				DefaultZoneSOARetry:   defaultZoneSOARetry,
				DefaultZoneSOAExpire:  defaultZoneSOAExpire,
				DefaultZoneSOAMinTTL:  defaultZoneSOAMinTTL,
				Mode:                  "forward",
				Path:                  basePath + "/zones/" + zone.Name + "." + zone.Network + ".forward.zone"}

			// Parse template
			t, err := template.New("zones").Parse(bindZoneFileTemplate)
			check(err)
			// Create zone file
			f, err := os.Create(PackagedZoneStructure.Path)
			check(err)
			// Execute zone file templating
			err = t.Execute(f, PackagedZoneStructure)
			check(err)
			// Close and write file
			f.Close()
		} else {
			return false, Stoerr("Name and/or network not defined!")
		}
	}

	return true, nil
}

const bindZoneFileTemplate = `$ORIGIN {{ .Zone.Name }}.
$TTL {{ .TTL }}

@ IN  SOA	{{ .Zone.PrimaryDNSServer }}. hostmaster.{{ .Zone.Name }}. (
	{{ .SerialNumber }}
	{{ .DefaultZoneSOARefresh }}
	{{ .DefaultZoneSOARetry }}
	{{ .DefaultZoneSOAExpire }}
	{{ .DefaultZoneSOAMinTTL }} )

{{ with .Zone.Records.NS }}{{ range . }}
{{ .Anchor }} {{ .TTL }} IN NS {{ .Name }}.{{ .Domain }}{{ end }}{{ end }}

{{ with .Zone.Records.MX }}{{ range . }}
{{ .Name }} {{ .TTL }} IN MX {{ .Priority }} {{ .Value }}{{ end }}{{ end }}

{{ with .Zone.Records.A }}{{ range . }}
{{ .Name }} {{ .TTL }} IN A {{ .Value }}{{ end }}{{ end }}

{{ with .Zone.Records.AAAA }}{{ range . }}
{{ .Name }} {{ .TTL }} IN AAAA {{ .Value }}{{ end }}{{ end }}

{{ with .Zone.Records.CNAME }}{{ range . }}
{{ .Name }} {{ .TTL }} IN CNAME {{ .Value }}{{ end }}{{ end }}

{{ with .Zone.Records.TXT }}{{ range . }}
{{ .Name }} {{ .TTL }} IN TXT {{ .Value }}{{ end }}{{ end }}

{{ with .Zone.Records.SRV }}{{ range . }}
{{ .Name }} {{ .TTL }} IN SRV {{ .Priority }} {{ .Weight }} {{ .Port }} {{ .Value }}{{ end }}{{ end }}
`
