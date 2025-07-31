package main

import (
	"fmt"

	"github.com/Pietot/Gonnect-4/grid"
)

func main() {
	testGrid, _ := grid.InitGrid("3446666575")
	eval, stats := testGrid.Negamax()
	fmt.Printf("%v\n%v", eval, stats)
}
