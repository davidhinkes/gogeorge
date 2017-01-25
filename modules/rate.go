package modules

type Rate []float64

func MakeRate(historySize int) Rate {
	return Rate(make([]float64, historySize))
}

func (r Rate) Do(x float64) float64 {
	n := len(r)
	for i := 0; i < n-1; i++ {
		r[i] = r[i+1]
	}
	r[n-1] = x
	m, _ := LeastSquares(r)
	return m
}
