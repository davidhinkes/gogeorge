package modules

import (
	"math"
)

type Mod float64

func (m Mod) Do(x float64) float64 {
	return math.Mod(x, float64(m))
}
