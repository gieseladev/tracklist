/*
Package celeste provides a tracklist parser.

Format: <index> [<start timestamp>] <track name>
*/
package celeste

import (
	"fmt"
	"github.com/gieseladev/tracklist/common"
	"github.com/gieseladev/tracklist/timestamp"
	"regexp"
	"strconv"
)

var lineMatcher = regexp.MustCompile(fmt.Sprintf(`^(\d+)\s*\[(%s)]\s*(.+?)\s*$`, common.TimestampMatcher.String()))

type celesteParser struct{}

func (p celesteParser) Parse(text string) (common.List, error) {
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

		sec, err := timestamp.ParseTimestamp(match[2])
		if err != nil {
			return common.Track{}, false
		}

		return common.Track{
			StartOffsetMS: 1000 * sec,
			Name:          string(match[3]),
		}, true
	}).Parse(text)
}

var Parser = celesteParser{}
