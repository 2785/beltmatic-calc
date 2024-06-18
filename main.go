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
	Goal        int
	SourceSet   []int
	OperatorSet string
}{
	SourceSet: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11},
	// div is stupid, let's not do div
	OperatorSet: "+-*^",
}

var calc = &cobra.Command{
	Use: "calc",
	Run: func(cmd *cobra.Command, args []string) {
		if calcArgs.Goal == 0 {
			c(cmd.Usage())
			return
		}

		opSet := operatorSet{}

		for _, c := range calcArgs.OperatorSet {
			switch c {
			case '+':
				opSet.add = true
			case '-':
				opSet.sub = true
			case '*':
				opSet.mul = true
			case '^':
				opSet.exp = true
			default:
				panic(fmt.Sprintf("unknown operator: %c", c))
			}
		}

		ops := findMostConvenientMadeUp(calcArgs.Goal, calcArgs.SourceSet, opSet)

		fmt.Println(renderOps(ops))
	},
}
