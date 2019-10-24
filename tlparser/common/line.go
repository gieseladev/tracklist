package common

import (
	"bufio"
	"bytes"
	"strings"
)

type LineParser struct {
	AllowNonSequential bool

	parseLine  func(line []byte) (Track, bool)
	ignoreLine func(line []byte) bool
	split      bufio.SplitFunc
}

func (lp *LineParser) IgnoreLine(ignoreLine func(line []byte) bool) {
	lp.ignoreLine = ignoreLine
}
func (lp *LineParser) Split(split bufio.SplitFunc) {
	lp.split = split
}

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

func NewLineParser(parseLine func(line []byte) (Track, bool)) *LineParser {
	if parseLine == nil {
		panic("passed nil parseLine func")
	}

	return &LineParser{parseLine: parseLine}
}
