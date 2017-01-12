package nodeto

import (
	"testing"
)

type testType struct{}

func (t testType) Do(ctx Context, i int) float64 {
	return float64(i)
}

func TestGetModuleTypes(t *testing.T) {
	if _, _, err := getModuleTypes(testType{}); err != nil {
		t.Error(err)
	}
	if _, _, err := getModuleTypes(""); err == nil {
		t.Errorf("Should have gotten an error for type string")
	}
}
