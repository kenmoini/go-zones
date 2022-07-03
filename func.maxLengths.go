package main

import "strconv"

// calculateMaxRecordComponentLength loops through a Zone's records and calculates the maximum length of each record component
func calculateMaxRecordComponentLength(zone Zone) MaxLengths {

	//=================================================
	// Calculate the forward zones lengths
	var maxLengths MaxLengths
	var AMap = map[string]int{"Name": 0, "Value": 0, "TTL": 0}
	var AAAAMap = map[string]int{"Name": 0, "Value": 0, "TTL": 0}
	var CNAMEMap = map[string]int{"Name": 0, "Value": 0, "TTL": 0}
	var TXTMap = map[string]int{"Name": 0, "Value": 0, "TTL": 0}
	var PTRMap = map[string]int{"Name": 0, "Value": 0, "TTL": 0}
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

	// Loop through each PTR record type and set the longest record component lengths to the MaxLengths record string map
	for _, record := range zone.Records.PTR {
		if len(record.Name) > PTRMap["Name"] {
			PTRMap["Name"] = len(record.Name)
		}
		if len(strconv.Itoa(record.TTL)) > PTRMap["TTL"] {
			PTRMap["TTL"] = len(strconv.Itoa(record.TTL))
		}
		if len(record.Value) > PTRMap["Value"] {
			PTRMap["Value"] = len(record.Value)
		}
	}

	maxLengths.NS = NSMap
	maxLengths.A = AMap
	maxLengths.AAAA = AAAAMap
	maxLengths.CNAME = CNAMEMap
	maxLengths.MX = MXMap
	maxLengths.SRV = SRVMap
	maxLengths.TXT = TXTMap
	maxLengths.PTR = PTRMap

	return maxLengths
}
