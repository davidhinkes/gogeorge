package modules

import (
  "math"
  "testing"
)

func TestLeastSquares(t *testing.T) {
  m := 2.0
  b := 1.0
  n := 500
  xs := make([]float64, n)
  for i := range xs {
  xs[i] = m*float64(i)+b
  }
  gotM, gotB := leastSquares(xs)
  same(t, m, gotM, 1e-5, "m")
  same(t, b, gotB, 1e-5, "b")
}

func same(t *testing.T, w float64, g float64, delta float64, name string) {
  if math.Abs(w-g) <= delta {
    return
  }
  t.Errorf("%v want %v got %v", name, w, g)
}
