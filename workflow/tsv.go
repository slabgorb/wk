package workflow

import (
	"encoding/csv"
	"io"
)

func ReadTSV(f io.Reader) ([]map[string]string, error) {
	reader := csv.NewReader(f)
	reader.Comma = '\t'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	headers := data[0]
	var output []map[string]string
	for _, d := range data[1:len(data)] {
		row := make(map[string]string)
		for i, h := range headers {
			row[h] = d[i]
		}
		output = append(output, row)
	}

	return output, nil
}
