package workflow_test

import (
	"testing"

	. "github.com/slabgorb/wk/workflow"
)

var exampleBadStep = &Step{
	Command: "foo",
	Arguments: map[string]string{
		"bar": "baz",
	},
	Parallel: false,
	LineNo:   1,
}

func TestStepError(t *testing.T) {
	ts := &*exampleBadStep
	ts.SetError("WHOOOPS")
	expected := "error running step foo line 1: WHOOOPS"
	if ts.Error() != expected {
		t.Errorf("expected %s got %s", expected, ts.Error())
	}
}

func TestBadCommand(t *testing.T) {
	ts := &*exampleBadStep
	_, err := ts.Run()
	expected := "error running step foo line 1: unknown command 'foo'"
	if err.Error() != expected {
		t.Errorf("expected %s got %s", expected, err.Error())
	}
	if ts.Error() != expected {
		t.Errorf("expected %s got %s", expected, err.Error())
	}

}
