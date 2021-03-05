package main

import (
	"time"
)

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// CLIOpts defines the CLI Arguements
type CLIOpts struct {
	Config string
	Mode   string
	Source string
	Dir    string
}

// Config struct for webapp config at the top level
type Config struct {
	Application ApplicationYaml `yaml:"app"`
}

// ApplicationYaml is what is defined for this Application when running as a server
type ApplicationYaml struct {
	ServerEnabled bool   `yaml:"server_enabled"`
	Server        Server `yaml:"server,omitempty"`
}

// Server configures the HTTP server
type Server struct {
	// Host is the local machine IP Address to bind the HTTP Server to
	Host string `yaml:"host"`

	BasePath string `yaml:"base_path"`

	// Port is the local machine TCP Port to bind the HTTP Server to
	Port    string `yaml:"port"`
	Timeout struct {
		// Server is the general server timeout to use
		// for graceful shutdowns
		Server time.Duration `yaml:"server"`

		// Write is the amount of time to wait until an HTTP server
		// write opperation is cancelled
		Write time.Duration `yaml:"write"`

		// Read is the amount of time to wait until an HTTP server
		// read operation is cancelled
		Read time.Duration `yaml:"read"`

		// Read is the amount of time to wait
		// until an IDLE HTTP session is closed
		Idle time.Duration `yaml:"idle"`
	} `yaml:"timeout"`
}

// Zones is the overall Zones struct
type Zones struct {
	Zones []ZonesYaml `yaml:"zones"`
}

// ZonesYaml is what each Zone is set up as
type ZonesYaml struct {
	Name             string  `yaml:"name"`
	Network          string  `yaml:"network"`
	PrimaryDNSServer string  `yaml:"primary_dns_server"`
	SubnetV4         string  `yaml:"subnet,omitempty"`
	SubnetV6         string  `yaml:"subnet_v6,omitempty"`
	TTL              int     `yaml:"ttl,omitempty"`
	Records          Records `yaml:"records,omitempty"`
}

// Records is a collection of different record types
type Records struct {
	NS    []NSRecord    `yaml:"NS,omitempty"`
	A     []ARecord     `yaml:"A,omitempty"`
	MX    []MXRecord    `yaml:"MX,omitempty"`
	AAAA  []AAAARecord  `yaml:"AAAA,omitempty"`
	CNAME []CNAMERecord `yaml:"CNAME,omitempty"`
	TXT   []TXTRecord   `yaml:"TXT,omitempty"`
	SRV   []SRVRecord   `yaml:"SRV,omitempty"`
}

// NSRecord is an NS Record definition
type NSRecord struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	Anchor string `yaml:"anchor"`
	TTL    int    `yaml:"ttl,omitempty"`
}

// ARecord is an A Record definition
type ARecord struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
	TTL   int    `yaml:"ttl,omitempty"`
}

// AAAARecord is an AAAA Record definition
type AAAARecord struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
	TTL   int    `yaml:"ttl,omitempty"`
}

// CNAMERecord is a CNAME Record definition
type CNAMERecord struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
	TTL   int    `yaml:"ttl,omitempty"`
}

// MXRecord is an MX Record definition
type MXRecord struct {
	Name     string `yaml:"name"`
	Value    string `yaml:"value"`
	Priority int    `yaml:"priority"`
	TTL      int    `yaml:"ttl,omitempty"`
}

// SRVRecord is an SRV Record definition
type SRVRecord struct {
	Name     string `yaml:"name"`
	Value    string `yaml:"value"`
	Priority int    `yaml:"priority"`
	Port     int    `yaml:"port"`
	Weight   int    `yaml:"weight"`
	TTL      int    `yaml:"ttl,omitempty"`
}

// TXTRecord is an TXT Record definition
type TXTRecord struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
	TTL   int    `yaml:"ttl,omitempty"`
}

// BindZoneConfig will setup the bind config for zones
type BindZoneConfig struct {
	Network string // Network type (internal/external/etc)
	Name    string // what zone is being served, example.com
	Path    string // Path to the Zones file
	Mode    string // Mode is forward or reverse zone
}

// PackagedZone is what is fed to the Zone Template
type PackagedZone struct {
	Zone                  ZonesYaml
	TTL                   int
	Mode                  string
	SerialNumber          string
	Path                  string
	DefaultZoneSOARefresh string
	DefaultZoneSOARetry   string
	DefaultZoneSOAExpire  string
	DefaultZoneSOAMinTTL  int
}
