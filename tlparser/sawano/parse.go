/*
Package sawano provides a tracklist parser.

Format: <start timestamp> - <end timestamp> <track name>
*/
package sawano

import (
	"fmt"
	"github.com/gieseladev/tracklist/timestamp"
	"github.com/gieseladev/tracklist/tlparser/common"
	"regexp"
)

var lineMatcher = regexp.MustCompile(fmt.Sprintf(`^(%s)\s*\p{Pd}\s*(%[1]s)\s*(.+?)\s*$`, common.TimestampMatcher.String()))

var Parser = common.NewLineParser(func(line []byte) (common.Track, bool) {
	match := lineMatcher.FindSubmatch(line)
	if match == nil {
		return common.Track{}, false
	}

	startOffset, err := timestamp.ParseTimestamp(match[1])
	if err != nil {
		return common.Track{}, false
	}
	endOffset, err := timestamp.ParseTimestamp(match[2])
	if err != nil {
		return common.Track{}, false
	}

	return common.Track{
		StartOffsetMS: 1000 * startOffset,
		EndOffsetMS:   1000 * endOffset,
		Name:          string(match[3]),
	}, true
})
