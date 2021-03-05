package main

var readConfig *Config

// BUFFERSIZE is for copying files
const BUFFERSIZE int64 = 4096 // 4096 bits = default page size on OSX

const appName string = "GoZones"
const appVersion string = "0.0.1"
const serverUA = appName + "/" + appVersion

const hexDigit = "0123456789abcdef"
const defaultTTL = 3600

// defaultZoneSOA_Refresh defines the default SOA refresh
const defaultZoneSOARefresh = "6h"

// defaultZoneSOA_Retry defines the default SOA retry
const defaultZoneSOARetry = "1h"

// defaultZoneSOA_Expire defines the default SOA expiration
const defaultZoneSOAExpire = "1w"

// defaultZoneSOA_Min_ttl defines the default SOA min TTL
const defaultZoneSOAMinTTL = 600
