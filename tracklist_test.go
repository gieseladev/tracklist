package tracklist

import (
	"encoding/csv"
	"errors"
	"github.com/gieseladev/tracklist/timestamp"
	"github.com/gieseladev/tracklist/tlparser/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type TestCase struct {
	Name     string
	Text     string
	Expected List
}

func (c TestCase) Run(t *testing.T) {
	actual, err := Parse(c.Text)
	require.NoError(t, err)

	assert.EqualValues(t, c.Expected, actual)
}

func loadTracklist(r io.Reader) (List, error) {
	reader := csv.NewReader(r)
	reader.ReuseRecord = true

	tl := List{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return tl, err
		}

		hasEndOffset := len(record) == 3

		if !hasEndOffset && len(record) != 2 {
			return tl, errors.New("unexpected field count")
		}

		startOffset, err := timestamp.ParseTimestamp([]byte(strings.TrimSpace(record[0])))
		if err != nil {
			return tl, err
		}

		nameFieldIndex := 1
		if hasEndOffset {
			nameFieldIndex = 2
		}

		track := common.Track{
			StartOffsetMS: 1000 * startOffset,
			Name:          record[nameFieldIndex],
		}

		if hasEndOffset {
			endOffset, err := timestamp.ParseTimestamp([]byte(strings.TrimSpace(record[1])))
			if err != nil {
				return tl, err
			}

			track.EndOffsetMS = 1000 * endOffset
		} else if len(tl.Tracks) > 0 {
			tl.Tracks[len(tl.Tracks)-1].EndOffsetMS = track.StartOffsetMS
		}

		tl.Tracks = append(tl.Tracks, track)
	}

	return tl, nil
}

func LoadTestCase(textPath, tracklistPath string) (TestCase, error) {
	textFile, err := os.Open(textPath)
	if err != nil {
		return TestCase{}, err
	}
	defer func() { _ = textFile.Close() }()

	data, err := ioutil.ReadAll(textFile)
	if err != nil {
		return TestCase{}, err
	}

	testCase := TestCase{
		Name: strings.TrimSuffix(filepath.Base(textPath), filepath.Ext(textPath)),
		Text: string(data),
	}

	tracklistFile, err := os.Open(tracklistPath)
	if err != nil {
		return testCase, err
	}
	defer func() { _ = tracklistFile.Close() }()

	tracklist, err := loadTracklist(tracklistFile)
	if err != nil {
		return testCase, err
	}

	testCase.Expected = tracklist

	return testCase, nil
}

func LoadTestCases(dirname string) ([]TestCase, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	if len(files)%2 != 0 {
		return nil, errors.New("even amount of files required")
	}

	var cases []TestCase
	for i := 1; i < len(files); i += 2 {
		tracklistPath := filepath.Join(dirname, files[i-1].Name())
		textPath := filepath.Join(dirname, files[i].Name())

		testCase, err := LoadTestCase(textPath, tracklistPath)
		if err != nil {
			return cases, err
		}

		cases = append(cases, testCase)
	}

	return cases, nil
}

func TestParse(t *testing.T) {
	testCases, err := LoadTestCases("test/data/tracklists")
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.Name, testCase.Run)
	}
}
