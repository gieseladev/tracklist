/*
Package mili provides a tracklist parser.

Format: <track number> - <track name> <start timestamp>
*/
package mili

import (
	"fmt"
	"github.com/gieseladev/tracklist/timestamp"
	"github.com/gieseladev/tracklist/tlparser/common"
	"regexp"
	"strconv"
)

var lineMatcher = regexp.MustCompile(fmt.Sprintf(`^(\d+)\s*\p{Pd}\s*(.+?)\s*(%s)\s*$`, common.TimestampMatcher.String()))

type miliParser struct{}

func (p miliParser) Parse(text string) (common.List, error) {
	prevTrackNumber := 0

	return common.NewLineParser(func(line []byte) (common.Track, bool) {
		match := lineMatcher.FindSubmatch(line)
		if match == nil {
			return common.Track{}, false
		}

		trackNumber, err := strconv.Atoi(string(match[1]))
		if err != nil || trackNumber != prevTrackNumber+1 {
			return common.Track{}, false
		}
		prevTrackNumber = trackNumber

		sec, err := timestamp.ParseTimestamp(match[3])
		if err != nil {
			return common.Track{}, false
		}

		return common.Track{
			StartOffsetMS: 1000 * sec,
			Name:          string(match[2]),
		}, true
	}).Parse(text)
}

var Parser = miliParser{}
