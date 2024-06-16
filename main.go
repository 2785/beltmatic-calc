package main

import (
	"fmt"

	"github.com/2785/libgo/args"
	"github.com/spf13/cobra"
)

func main() {
	c(args.RegisterArgs(&calcArgs, calc.Flags()))
	c(calc.Execute())
}

func c(e error) {
	if e != nil {
		panic(e)
	}
}

var calcArgs = struct {
	Goal      int
	SourceSet []int
}{
	SourceSet: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11},
}

var calc = &cobra.Command{
	Use: "calc",
	Run: func(cmd *cobra.Command, args []string) {
		if calcArgs.Goal == 0 {
			c(cmd.Usage())
			return
		}

		ops := findMostConvenientMadeUp(calcArgs.Goal, calcArgs.SourceSet)

		fmt.Println(renderOps(ops))
	},
}
