package tracklist

import (
	"errors"
	"github.com/gieseladev/tracklist/tlparser"
	"github.com/gieseladev/tracklist/tlparser/common"
)

type Track = common.Track
type List = common.List

var (
	ErrNoParser = errors.New("no parser could parse the format")
)

func Parse(text string) (List, error) {
	formatInvalid := false
	noTracklistCount := 0

	parsers := tlparser.All()
	for _, parser := range parsers {
		tl, err := parser.Parse(text)
		switch err {
		case common.ErrNoTracklist:
			noTracklistCount++
			continue
		case common.ErrInvalidFormat:
			formatInvalid = true
			continue
		case nil:
		default:
			return List{}, err
		}

		return tl, nil
	}

	if formatInvalid {
		return List{}, common.ErrInvalidFormat
	}

	if noTracklistCount == len(parsers) {
		return List{}, common.ErrNoTracklist
	}

	return List{}, ErrNoParser
}
