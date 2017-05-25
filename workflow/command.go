package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

// this is just for the purposes of this exercise, in an actual app
// this may be another argument in the workflow DSL
const InterestingColumn = "Volume 2015"

type commandFunc func(args map[string]string) (int, error)

var commandList = map[string]commandFunc{
	"fetch":       fetchCommand,
	"aggregation": aggregationCommand,
	"to_json":     toJsonCommand,
}

func checkArg(key string, args map[string]string) (string, error) {
	value, ok := args[key]
	if !ok {
		return "", fmt.Errorf("%s argument not passed", key)
	}
	return value, nil
}

// fetchCommand gets a file from the interwebs then saves it locally
// for further processing. Unlike other commands, we are uninterested
// in the byte count of the operation's output file.
func fetchCommand(args map[string]string) (int, error) {
	var url, file string
	url, err := checkArg("url", args)
	if err != nil {
		return 0, err
	}
	file, err = checkArg("file", args)
	if err != nil {
		return 0, err
	}
	output, err := os.Create(file)
	if err != nil {
		return 0, err
	}

	tries := 0
	for {
		response, err := http.Get(url)
		if err != nil {
			if tries > 2 {
				return 0, err
			} else {
				tries += 1
				continue
			}
		}
		defer response.Body.Close()
		_, err = io.Copy(output, response.Body)
		if err != nil {
			return 0, err
		}
		return 0, nil // bytes should not be counted towards total for purposes of this exercise
		break
	}
	return 0, nil
}

// Median returns the median value for a list. If the list is empty,
// Median returns NaN.
func Median(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()

	}
	sort.Float64s(values)
	mid := len(values) / 2
	median := values[mid]
	if len(values)%2 == 0 {
		median = (median + values[mid-1]) / 2
	}
	return median

}

func getData(tsvFile string) ([]map[string]string, error) {
	reader, err := LoadFile(tsvFile)
	if err != nil {
		return nil, err
	}
	return ReadTSV(reader)
}

// toJsonCommand converts a local TSV file to Json.
func toJsonCommand(args map[string]string) (int, error) {
	file, err := checkArg("file", args)
	if err != nil {
		return 0, err
	}
	tsvFile, err := checkArg("tsvfile", args)
	if err != nil {
		return 0, err
	}
	data, err := getData(tsvFile)
	if err != nil {
		return 0, err
	}
	fout, err := os.Create(file)
	if err != nil {
		return 0, err
	}
	defer fout.Close()
	j, err := json.MarshalIndent(data, " ", "    ")
	if err != nil {
		return 0, err
	}

	b, err := fout.Write(j)
	if err != nil {
		return 0, err
	}

	return b, nil
}

// aggregationCommand reads a local TSV file, gets min, max, and
// median values, and saves that to a local json file.
func aggregationCommand(args map[string]string) (int, error) {
	file, err := checkArg("file", args)
	if err != nil {
		return 0, err
	}
	tsvFile, err := checkArg("tsvfile", args)
	if err != nil {
		return 0, err
	}
	data, err := getData(tsvFile)
	if err != nil {
		return 0, err
	}
	if len(data) < 2 {
		return 0, fmt.Errorf("no data read")
	}

	column := []float64{}
	max := 0.0
	min := 1000000000.0

	for i, d := range data {
		datum := strings.Trim(d[InterestingColumn], " ")
		f, err := strconv.ParseFloat(datum, 64)
		if err != nil {
			return 0, fmt.Errorf("could not parse float from data in line %d, val is %s", i, datum)
		}
		if f > max {
			max = f
		}
		if f < min {
			min = f
		}
		column = append(column, f)
	}

	fout, err := os.Create(file)
	if err != nil {
		return 0, err
	}
	defer fout.Close()

	output := map[string]float64{
		"max":    max,
		"min":    min,
		"median": Median(column),
	}

	j, err := json.MarshalIndent(output, " ", "    ")
	if err != nil {
		return 0, err
	}

	b, err := fout.Write(j)
	if err != nil {
		return 0, err
	}

	return b, nil
}
