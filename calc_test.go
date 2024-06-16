package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRednerOps(t *testing.T) {
	tc := []struct {
		ops      []Op
		expected string
	}{
		{
			ops: []Op{
				StaticOp{Num: 1},
				MultiplicationOp{Num: 2},
				AdditionOp{Num: 3},
				SubtractionOp{Num: 4},
			},
			expected: "1x2+3-4",
		},
		{
			ops: []Op{
				StaticOp{Num: 1},
				AdditionOp{Num: 2},
				MultiplicationOp{Num: 3},
				SubtractionOp{Num: 4},
			},
			expected: "(1+2)x3-4",
		},
	}

	for _, tt := range tc {
		actual := renderOps(tt.ops)
		assert.Equal(t, tt.expected, actual)
	}
}

func TestMath(t *testing.T) {
	tc := []struct {
		target    int
		sourceSet []int
		expected  string
	}{
		{
			target:    10,
			sourceSet: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11},
			expected:  "9+1",
		},
		{
			target: 6,
			sourceSet: []int{
				2, 3, 5,
			},
			expected: "3x2",
		},
		{
			target: 584,
			sourceSet: []int{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 11,
			},
			expected: "(9x8+1)x8",
		},
		{
			target: 662,
			sourceSet: []int{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 11,
			},
			expected: "(11x6x5+1)x2",
		},
	}

	for _, tt := range tc {
		ops := findMostConvenientMadeUp(tt.target, tt.sourceSet)
		actual := renderOps(ops)

		assert.Equal(t, tt.expected, actual)
	}
}
