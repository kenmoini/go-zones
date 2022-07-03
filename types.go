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

type RootYAML struct {
	DNS DNS `yaml:"dns"`
}

type TemplatePair struct {
	DNS         DNS                 `yaml:"dns"`
	BasePath    string              `yaml:"base_path"`
	RevViewPair map[string][]string `yaml:"revViewPair,omitempty"`
}

// DNS contains the overall DNS configuration such as the different forwarders, views, and zones
type DNS struct {
	ACLs  []ACL  `yaml:"acls"`
	Views []View `yaml:"views,omitempty"`
	Zones []Zone `yaml:"zones"`
}

// ACL composes a name and network slice pair
type ACL struct {
	Name     string   `yaml:"name"`
	Networks []string `yaml:"networks"`
}

// ServerConfig provides overrides for the server configuration

// View contains the configuration for a DNS view that can set the networks that Zones are applied to
type View struct {
	Name           string          `yaml:"name"`
	ACLs           []string        `yaml:"acls"`                      // ACL is the list of the ACLs that this view is applied to
	Recursion      bool            `yaml:"recursion"`                 // Recursion is the recursion setting for this view
	Forwarders     []string        `yaml:"forwarders,omitempty"`      // Forwarders is the list of DNS forwarders to use
	IncludedZones  []string        `yaml:"zones,omitempty"`           // IncludedZones is a list of the named Zones that are set for this view
	ForwardedZones []ForwardedZone `yaml:"forwarded_zones,omitempty"` // Forwarded Zones is a list of the Zones that are being forwarded with this view
}

type ReverseViewPair struct {
	View         string   `yaml:"view"`
	ReverseZones []string `yaml:"reverse"`
}

// ForwardedZone contains the configuration for a DNS forwarder per view
type ForwardedZone struct {
	Zone       string   `yaml:"zone"`       // Zone is the domain of the zone to forward
	Forwarders []string `yaml:"forwarders"` // Forwarders is the list of DNS forwarders to send requests to
}

// Zone is what each Zone is set up as
type Zone struct {
	Name             string  `yaml:"name"`
	Zone             string  `yaml:"zone"`
	PrimaryDNSServer string  `yaml:"primary_dns_server"`
	DefaultTTL       int     `yaml:"default_ttl,omitempty"`
	Records          Records `yaml:"records,omitempty"`
	//Network          string  `yaml:"network"`
	//SubnetV4         string  `yaml:"subnet,omitempty"`
	//SubnetV6         string  `yaml:"subnet_v6,omitempty"`
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
	PTR   []PTRRecord   `yaml:"PTR,omitempty"`
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
	NoPTR bool   `yaml:"no_ptr,omitempty"`
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

// PTRRecord is an TXT Record definition
type PTRRecord struct {
	Name              string `yaml:"name"`
	Value             string `yaml:"value"`
	TargetReverseZone string `yaml:"target_reverse_zone"`
	TTL               int    `yaml:"ttl,omitempty"`
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
	Zone                  Zone
	TTL                   int
	Mode                  string
	SerialNumber          string
	Path                  string
	DefaultZoneSOARefresh string
	DefaultZoneSOARetry   string
	DefaultZoneSOAExpire  string
	DefaultZoneSOAMinTTL  int
	MaxLengths            MaxLengths
}

// PackagedReverseZone is what is fed to the Reverse Zone Template
type PackagedReverseZone struct {
	Zone                  Zone
	ReverseName           string
	TTL                   int
	Mode                  string
	SerialNumber          string
	Path                  string
	DefaultZoneSOARefresh string
	DefaultZoneSOARetry   string
	DefaultZoneSOAExpire  string
	DefaultZoneSOAMinTTL  int
}

// MaxLengths is a a map of the max lengths for each record type
type MaxLengths struct {
	NS    map[string]int
	A     map[string]int
	AAAA  map[string]int
	MX    map[string]int
	CNAME map[string]int
	TXT   map[string]int
	SRV   map[string]int
	PTR   map[string]int
}

type ReverseZoneRecords struct {
	PTR []PTRRecord
}
