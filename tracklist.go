/*
Package tracklist pares various tracklist formats.
*/
package tracklist

import (
	"errors"
	"github.com/gieseladev/tracklist/common"
)

type Track = common.Track
type List = common.List

var (
	ErrNoParser      = errors.New("no parser could parse the format")
	ErrNoTracklist   = common.ErrNoTracklist
	ErrInvalidFormat = common.ErrInvalidFormat
)

func Parse(text string) (List, error) {
	formatInvalid := false
	noTracklistCount := 0

	parsers := All()
	for _, parser := range parsers {
		tl, err := parser.Parse(text)
		switch err {
		case ErrNoTracklist:
			noTracklistCount++
			continue
		case ErrInvalidFormat:
			formatInvalid = true
			continue
		case nil:
		default:
			return List{}, err
		}

		return tl, nil
	}

	if formatInvalid {
		return List{}, ErrInvalidFormat
	}

	if noTracklistCount == len(parsers) {
		return List{}, ErrNoTracklist
	}

	return List{}, ErrNoParser
}
