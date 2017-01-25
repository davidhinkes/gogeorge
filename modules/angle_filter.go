package modules

type AngleFilter struct {
	last float64
}

func (a *AngleFilter) Do(x float64) float64 {
	for ;(x - a.last) > 180; x -= 360 {}
	for ;(x - a.last) < -180; x += 360 {}
	a.last = x
	return x
}
