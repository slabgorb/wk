package workflow_test

import (
	"bytes"
	"testing"

	. "github.com/slabgorb/wk/workflow"
)

var exampleFile = `
command:fetch # this is a comment
  url:http://www.google.com
  tries:3
parallel:median
  column:2
`
var badFile1 = `
command: # this is a badly formatted command, and should throw
  url:http://www.google.com
command:median
  column:2
`
var badFile2 = `
command:median
  column
`

func TestHappyPathParse(t *testing.T) {
	buf := bytes.NewBufferString(exampleFile)
	steps, err := ParseSteps(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(steps.List) != 2 {
		t.Errorf("expected 2 steps, got %d", len(steps.List))
	}
	if steps.List[0].Command != "fetch" {
		t.Errorf("expected first command to be 'fetch', got %s", steps.List[0].Command)
	}
	url, ok := steps.List[0].Arguments["url"]
	if !ok || url == "" {
		t.Errorf("did not properly parse 'url' argument to command")
	}
	tries, ok := steps.List[0].Arguments["tries"]
	if !ok || tries != "3" {
		t.Errorf("did not properly parse 'tries' argument to command")
	}
}

func TestBadFileParse(t *testing.T) {
	buf := bytes.NewBufferString(badFile1)
	_, err := ParseSteps(buf)
	expected := "parsing error at line 2, 'command: # this is a badly formatted command, and should throw' malformed command"
	if err.Error() != expected {
		t.Errorf("expected error %s, got %s", expected, err)
	}
	buf = bytes.NewBufferString(badFile2)
	_, err = ParseSteps(buf)
	expected = "parsing error for at line 3, '  column' malformed argument"
	if err.Error() != expected {
		t.Errorf("expected error %s, got %s", expected, err)
	}
}
