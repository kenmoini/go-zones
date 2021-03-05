package main

import (
	"github.com/gosimple/slug"
)

// slugger slugs a string
func slugger(textToSlug string) string {
	return slug.Make(textToSlug)
}

// Convert unsigned integer to decimal string.
func uitoa(val uint) string {
	if val == 0 { // avoid string allocation
		return "0"
	}
	var buf [20]byte // big enough for 64bit value base 10
	i := len(buf) - 1
	for val >= 10 {
		q := val / 10
		buf[i] = byte('0' + val - q*10)
		i--
		val = q
	}
	// val < 10
	buf[i] = byte('0' + val)
	return string(buf[i:])
}
