package tracklist

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type TestCase struct {
	Name     string
	Text     string
	Expected Tracklist
}

func (c TestCase) Run(t *testing.T) {
	assert.Equal(t, c.Text, c.Expected)
}

func parseTimestamp(ts string) (uint32, error) {
	var s uint32
	fields := strings.Split(ts, ":")

	mult := uint32(1000)
	for i := len(fields) - 1; i >= 0; i-- {
		n, err := strconv.Atoi(fields[i])
		if err != nil {
			return s, err
		}

		s += uint32(n) * mult
		mult *= 60
	}

	return s, nil
}

func loadTracklist(r io.Reader) (Tracklist, error) {
	reader := csv.NewReader(r)
	reader.ReuseRecord = true

	tracklist := Tracklist{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return tracklist, err
		}

		if len(record) != 2 {
			return tracklist, errors.New("unexpected fields count")
		}

		offset, err := parseTimestamp(strings.TrimSpace(record[0]))
		if err != nil {
			return tracklist, err
		}

		tracklist.Tracks = append(tracklist.Tracks, Track{
			StartOffsetMS: offset,
			Name:          record[1],
		})
	}

	return tracklist, nil
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
		Name: strings.TrimSuffix(path.Base(textPath), path.Ext(textPath)),
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

		fmt.Println(textPath, tracklistPath)

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
