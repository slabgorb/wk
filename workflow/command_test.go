package workflow_test

import (
	"math"
	"testing"

	. "github.com/slabgorb/wk/workflow"
)

func TestMedian(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5}
	if Median(values) != 3.0 {
		t.Errorf("Expected 3 got %f", Median(values))
	}
	values = []float64{2, 4}
	if Median(values) != 3.0 {
		t.Errorf("Expected 3 got %f", Median(values))
	}
	values = []float64{3}
	if Median(values) != 3.0 {
		t.Errorf("Expected 3 got %f", Median(values))
	}
	if !math.IsNaN(Median([]float64{})) {
		t.Errorf("Expected NaN got %f", Median([]float64{}))
	}
}
