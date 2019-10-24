/*
Package inception provides a tracklist parser.

Format: <start timestamp> <track name>
*/
package inception

import (
	"fmt"
	"github.com/gieseladev/tracklist/common"
	"github.com/gieseladev/tracklist/timestamp"
	"regexp"
)

var lineMatcher = regexp.MustCompile(fmt.Sprintf(`^(%s)\s*(.+?)\s*$`, common.TimestampMatcher.String()))

type inceptionParser struct{}

func (p inceptionParser) Parse(text string) (common.List, error) {
	allowTitle := false

	lineParser := common.NewLineParser(func(line []byte) (common.Track, bool) {
		match := lineMatcher.FindSubmatch(line)
		if match == nil {
			if allowTitle {
				allowTitle = false
				return common.Track{}, true
			}

			return common.Track{}, false
		}

		startOffset, err := timestamp.ParseTimestamp(match[1])
		if err != nil {
			return common.Track{}, false
		}

		allowTitle = true
		return common.Track{
			StartOffsetMS: 1000 * startOffset,
			Name:          string(match[2]),
		}, true
	})

	return lineParser.Parse(text)
}

var Parser = inceptionParser{}
