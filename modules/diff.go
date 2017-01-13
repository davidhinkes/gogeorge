package modules

type Diff struct{}

func (_ Diff) Do(a, b float64) float64 {
	return a - b
}
