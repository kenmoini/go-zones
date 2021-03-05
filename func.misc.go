package main

import (
	"github.com/gosimple/slug"
)

// slugger slugs a string
func slugger(textToSlug string) string {
	return slug.Make(textToSlug)
}
