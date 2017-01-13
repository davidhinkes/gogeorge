package modules

import (
  "testing"
  "testing/quick"
)

func TestQuickMod(t *testing.T) {
  m := 180.
  f := func(x float64) bool {
    return Mod(m).Do(x) <= m && Mod(m).Do(x) >= -m
  }
  if err := quick.Check(f, nil); err != nil {
    t.Error(err)
  }
}
