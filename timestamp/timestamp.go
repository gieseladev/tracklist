package timestamp

import (
	"bytes"
	"errors"
	"strconv"
)

var toSecs = [4]uint32{1, 60, 60 * 60, 24 * 60 * 60}

// ParseTimestamp parses a timestamp in the form of hh:mm:ss where all parts are
// optional. It will only parse up to dd:hh:mm:ss, if there are more parts, an
// error is returned.
func ParseTimestamp(ts []byte) (uint32, error) {
	parts := bytes.Split(ts, []byte{':'})
	if len(parts) >= len(toSecs) {
		return 0, errors.New("too many parts")
	}

	var secs uint32
	for i, part := range parts {
		n, err := strconv.ParseUint(string(part), 10, 0)
		if err != nil {
			return secs, err
		}

		secs += toSecs[len(parts)-i-1] * uint32(n)
	}

	return secs, nil
}
