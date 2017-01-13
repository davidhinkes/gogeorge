package modules

import (
	"log"

	m "github.com/skelterjohn/go.matrix"
)

type PID struct {
	history  []float64
	integral float64
	p        float64
	i        float64
	d        float64
}

func (p *PID) Do(e float64) (command float64) {
	for i := 0; i < len(p.history)-1; i++ {
		p.history[i] = p.history[i+1]
	}
	p.history[len(p.history)-1] = e
	p.integral += e
	rate, _ := leastSquares(p.history)
	command = p.p*e + p.i*p.integral + p.d*rate
	return
}

func MakePID(p, i, d float64, historySize int) *PID {
	return &PID{
		history: make([]float64, historySize),
		p:       p,
		i:       i,
		d:       d,
	}
}

func leastSquares(x []float64) (float64, float64) {
	aMatrix := m.Zeros(len(x), 2)
	yVector := m.Zeros(len(x), 1)
	for i := range x {
		aMatrix.Set(i, 0, float64(i))
		aMatrix.Set(i, 1, 1.0)
		yVector.Set(i, 0, x[i])
	}
	qMatrix, rMatrix := aMatrix.QR()
	rMatrix = multiply(eye(2, len(x)), rMatrix)
	rInverseMatrix, err := rMatrix.Inverse()
	if err != nil {
		log.Fatalf("R inverse: %v", err)
	}
	qTransposeMatrix := qMatrix.Transpose()
	xMatrix := multiply(rInverseMatrix, qTransposeMatrix, yVector)
	return xMatrix.Get(0, 0), xMatrix.Get(1, 0)
}

func multiply(ms ...*m.DenseMatrix) *m.DenseMatrix {
	ret := ms[0]
	var err error
	for i := 1; i < len(ms); i++ {
		ret, err = ret.TimesDense(ms[i])
		if err != nil {
			log.Fatalf("multiply: %v", err)
		}
	}
	return ret
}

func eye(i, j int) *m.DenseMatrix {
	ret := m.Zeros(i, j)
	n := i
	if j < i {
		n = j
	}
	for k := 0; k < n; k++ {
		ret.Set(k, k, 1)
	}
	return ret
}
