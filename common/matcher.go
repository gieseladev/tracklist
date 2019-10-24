package common

import "regexp"

//var TimestampMatcher = regexp.MustCompile(`\d+(?::\d+)*`)
var TimestampMatcher = regexp.MustCompile(`\d+(?::\d+)+`)
