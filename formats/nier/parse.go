/*
Package nier provides a tracklist parser.

Format: <track number>. <track name> <start timestamp>
*/
package nier

import (
	"fmt"
	"github.com/gieseladev/tracklist/common"
	"github.com/gieseladev/tracklist/timestamp"
	"regexp"
	"strconv"
)

var lineMatcher = regexp.MustCompile(fmt.Sprintf(`^(\d+)[.:]\s*(.+?)\s*(%s)\s*$`, common.TimestampMatcher.String()))

type nierParser struct{}

func (p nierParser) Parse(text string) (common.List, error) {
	prevTrackNumber := 0

	lineParser := common.NewLineParser(func(line []byte) (common.Track, bool) {
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
	})

	return lineParser.Parse(text)
}

var Parser = nierParser{}
