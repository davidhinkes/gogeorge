package nodeto

import (
	"fmt"
	"testing"
)

type doubler struct{}

func (_ doubler) Do(ctx Context, x int) int {
	return 2 * x
}

type incrementer struct {
	increment int
}

func (i *incrementer) Do(ctx Context) int {
	i.increment++
	return i.increment
}

type printer struct{}

func (_ printer) Do(ctx Context, x interface{}) {
	fmt.Printf("\r%v -> %v", ctx.Iteration(), x)
}

func TestExample(t *testing.T) {
	var board Board
	count := board.Bind(&incrementer{})[0]
	doubled := board.Bind(doubler{}, count)[0]
	board.Bind(printer{}, doubled)
}
