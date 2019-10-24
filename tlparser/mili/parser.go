package mili

import (
	"fmt"
	"github.com/gieseladev/tracklist/timestamp"
	"github.com/gieseladev/tracklist/tlparser/common"
	"regexp"
)

var lineMatcher = regexp.MustCompile(fmt.Sprintf(`^(\d+)\s*\p{Pd}\s*(.+?)\s*(%s)$`, common.TimestampMatcher.String()))

var Parser = common.NewLineParser(func(line []byte) (common.Track, bool) {
	match := lineMatcher.FindSubmatch(line)
	if match == nil {
		return common.Track{}, false
	}

	// TODO enforce increasing track numbers

	sec, err := timestamp.ParseTimestamp(match[3])
	if err != nil {
		return common.Track{}, false
	}

	return common.Track{
		StartOffsetMS: 1000 * sec,
		Name:          string(match[2]),
	}, true
})
