package modules

import (
	"log"

	"github.com/davidhinkes/gogeorge/internal/io"
	m "github.com/skelterjohn/go.matrix"
)

type getter interface {
	Get() float64
}

type PID struct {
	rate     Rate
	integral float64
	p        getter
	i        getter
	d        getter
}

func (p *PID) Do(e float64) float64 {
	p.integral += e
	d := p.rate.Do(e)
	return p.p.Get()*e + p.i.Get()*p.integral + p.d.Get()*d
}

func NewPID(p, i, d *io.Float64, historySize int) *PID {
	return &PID{
		rate: MakeRate(historySize),
		p:    p,
		i:    i,
		d:    d,
	}
}

func LeastSquares(x []float64) (float64, float64) {
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
