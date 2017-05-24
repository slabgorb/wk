// Package workflow provides reading and parsing for a workflow format
// file.

// A workflow file is expected to be in the following format:

/*

command:name # comments are denoted by '#' and anything on a line after a # is ignored.
  argumentname:value
  argumentname:value
  argumentname:value


parallelcommand:name
 argumentname:value
 argumentname:value

parallelcommand:name
 argumentname:value
 argumentname:value

command:name  # this command will only be executed
  argumentname:value

*/
package workflow

import (
	"io"
	"os"
)

func LoadFile(filepath string) (io.Reader, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return file, nil
}
