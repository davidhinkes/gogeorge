package nodeto

type Context interface {
	Iteration() uint64
}

type context struct {
	iteration uint64
}

func (c context) Iteration() uint64 {
	return c.iteration
}
