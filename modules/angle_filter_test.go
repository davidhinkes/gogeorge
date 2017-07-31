package modules

import (
	"testing"
)

func TestAngleFilter(t *testing.T) {
	cases := []struct {
		in   []float64
		want float64
	}{
		{in: []float64{3}, want: 3},
		{in: []float64{359}, want: -1},
		{in: []float64{3, 4}, want: 4},
		{in: []float64{3, -3}, want: -3},
		{in: []float64{300}, want: -60},
		{in: []float64{300, 1}, want: 1}, // this is correct, see above
	}
	for _, c := range cases {
		f := AngleFilter{}
		var got float64
		for _, x := range c.in {
			got = f.Do(x)
		}
		if got != c.want {
			t.Errorf("got %v want %v", got, c.want)
		}
	}
}
