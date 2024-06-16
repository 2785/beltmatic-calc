package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	calc.AddCommand(calc)
	c(calc.Execute())

}

func c(e error) {
	if e != nil {
		panic(e)
	}
}

var calcArgs = struct {
	goal      int
	sourceSet []int
}{
	sourceSet: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11},
}

var calc = &cobra.Command{
	Use: "calc",
	Run: func(cmd *cobra.Command, args []string) {
		if calcArgs.goal == 0 {
			c(cmd.Usage())
		}

		ops := findMostConvenientMadeUp(calcArgs.goal, calcArgs.sourceSet)

		fmt.Println(renderOps(ops))
	},
}
