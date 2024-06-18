package main

import (
	"testing"

	"github.com/samber/lo"
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
		opSet     *operatorSet
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
			expected: "(8x3)^2+8",
		},
		{
			target: 662,
			sourceSet: []int{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 11,
			},
			expected: "(6^3+5)x3-1",
		},
	}

	for _, tt := range tc {
		ops := findMostConvenientMadeUp(tt.target, tt.sourceSet, lo.FromPtrOr(tt.opSet, operatorSetAll))
		actual := renderOps(ops)

		assert.Equal(t, tt.expected, actual)
	}
}
