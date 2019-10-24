package common

import (
	"bufio"
	"bytes"
	"strings"
)

// LineParser is a Parser optimised for formats that list each track on its own
// line.
type LineParser struct {
	AllowNonSequential bool // When set to true, tracks don't have to be on subsequent lines.

	parseLine  func(line []byte) (Track, bool)
	ignoreLine func(line []byte) bool
	split      bufio.SplitFunc
}

// IgnoreLine sets the ignore line function. This function is called for each
// line and if it returns true, the line is skipped.
// By default it is nil which only ignores empty lines.
func (lp *LineParser) IgnoreLine(ignoreLine func(line []byte) bool) {
	lp.ignoreLine = ignoreLine
}

// Split sets the function that is used to split the text into lines.
func (lp *LineParser) Split(split bufio.SplitFunc) {
	lp.split = split
}

// Parse implements the Parser interface for LineParser.
// It takes each line in the given text and if it isn't ignored, passes it to
// the underlying parseLine function.
func (lp *LineParser) Parse(text string) (List, error) {
	if lp.parseLine == nil {
		panic("line parser with nil ParseLine func")
	}

	var tl List

	scanner := bufio.NewScanner(strings.NewReader(text))
	if lp.split != nil {
		scanner.Split(lp.split)
	}

	foundList := false
	expectEnd := false
	for lineIndex := 0; scanner.Scan(); lineIndex++ {
		line := bytes.TrimSpace(scanner.Bytes())
		if lp.ignoreLine != nil {
			if lp.ignoreLine(line) {
				continue
			}
		} else if len(line) == 0 {
			continue
		}

		track, ok := lp.parseLine(line)
		if !ok {
			if foundList {
				expectEnd = true
			}

			continue
		} else if track == (Track{}) {
			continue
		}

		// end offset must come after start offset
		if track.HasEndOffset() && track.StartOffsetMS > track.EndOffsetMS {
			return List{}, ErrInvalidFormat
		}

		if expectEnd {
			return List{}, ErrInvalidFormat
		}

		if tl.Len() > 0 {
			prevTrack := &tl.Tracks[tl.Len()-1]
			if !prevTrack.HasEndOffset() {
				prevTrack.EndOffsetMS = track.StartOffsetMS
			}

			// prev track has a later timestamp than the current one
			if prevTrack.StartOffsetMS >= track.StartOffsetMS {
				return List{}, ErrInvalidFormat
			}
		}

		tl.Tracks = append(tl.Tracks, track)
		foundList = true
	}

	if err := scanner.Err(); err != nil {
		return tl, err
	}

	if !foundList {
		return List{}, ErrNoTracklist
	}

	return tl, nil
}

// NewLineParser creates a new line parser which uses parseLine to parse
// each unignored line.
func NewLineParser(parseLine func(line []byte) (Track, bool)) *LineParser {
	if parseLine == nil {
		panic("passed nil parseLine func")
	}

	return &LineParser{parseLine: parseLine}
}
