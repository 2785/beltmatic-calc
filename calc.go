package main

import (
	"fmt"
	"math/big"

	"github.com/gammazero/deque"
	"github.com/samber/lo"
)

type Op interface {
	String() string
	Weight() int
}

type StaticOp struct {
	Num int
}

func (s StaticOp) String() string {
	return fmt.Sprintf("%d", s.Num)
}

func (s StaticOp) Weight() int {
	return 10
}

type MultiplicationOp struct {
	Num int
}

func (m MultiplicationOp) String() string {
	return fmt.Sprintf("x%d", m.Num)
}

func (m MultiplicationOp) Weight() int {
	return 2
}

type AdditionOp struct {
	Num int
}

func (a AdditionOp) String() string {
	return fmt.Sprintf("+%d", a.Num)
}

func (a AdditionOp) Weight() int {
	return 1
}

type SubtractionOp struct {
	Num int
}

func (s SubtractionOp) Weight() int {
	return 1
}

func (s SubtractionOp) String() string {
	return fmt.Sprintf("-%d", s.Num)
}

type ExponentiationOp struct {
	Num int
}

func (e ExponentiationOp) Weight() int {
	return 3
}

func (e ExponentiationOp) String() string {
	return fmt.Sprintf("^%d", e.Num)
}

func renderOps(ops []Op) string {
	if len(ops) == 0 {
		return ""
	}

	weight := ops[0].Weight()

	curr := ops[0].String()

	for i := 1; i < len(ops); i++ {
		w := ops[i].Weight()
		if w > weight {
			curr = fmt.Sprintf("(%s)", curr)
		}

		curr = fmt.Sprintf("%s%s", curr, ops[i].String())
		weight = w
	}

	return curr
}

type node struct {
	stepCount int
	ops       []Op
	val       int
}

type operatorSet struct {
	add bool
	sub bool
	mul bool
	exp bool
}

var operatorSetAll = operatorSet{
	add: true,
	sub: true,
	mul: true,
	exp: true,
}

func findMostConvenientMadeUp(target int, sourceSet []int, operatorSet operatorSet) []Op {
	// we're implementing a dijkstra with dynamic nodes as intermediates.

	if lo.Contains(sourceSet, target) {
		return []Op{StaticOp{Num: target}}
	}

	// origin node is our target
	origin := node{
		stepCount: 0,
		ops:       []Op{},
		val:       target,
	}

	nodes := make(map[int]*node)
	nodes[target] = &origin

	q := deque.New[*node]()
	q.PushBack(&origin)

	steps := 0

	for q.Len() > 0 {
		if q.Len() > 10000 {
			panic("too many steps, can't find the solution")
		}

		steps++

		curr := q.PopFront()

		if lo.Contains(sourceSet, curr.val) {
			// found it

			// need to reverse the ops

			staticOp := StaticOp{Num: curr.val}

			ops := lo.Reverse(curr.ops)

			ops = append([]Op{staticOp}, ops...)

			return ops
		}

		// didn't find it, generate new nodes
		for _, v := range sourceSet {
			// add
			if operatorSet.add {
				newVal := curr.val - v
				if newVal > 0 {
					if _, ok := nodes[newVal]; !ok {
						n := node{
							stepCount: curr.stepCount + 1,
							val:       newVal,
						}

						ops := make([]Op, len(curr.ops))
						copy(ops, curr.ops)

						ops = append(ops, AdditionOp{Num: v})

						n.ops = ops

						nodes[newVal] = &n
						q.PushBack(&n)
					}
				}
			}

			if operatorSet.sub {
				newVal := curr.val + v
				if _, ok := nodes[newVal]; !ok {
					n := node{
						stepCount: curr.stepCount + 1,
						val:       newVal,
					}

					ops := make([]Op, len(curr.ops))
					copy(ops, curr.ops)

					ops = append(ops, SubtractionOp{Num: v})

					n.ops = ops

					nodes[newVal] = &n
					q.PushBack(&n)
				}
			}

			if operatorSet.mul {
				if curr.val%v == 0 && v > 1 {
					// divide
					newVal := curr.val / v
					if _, ok := nodes[newVal]; !ok {
						ops := make([]Op, len(curr.ops))
						copy(ops, curr.ops)

						ops = append(ops, MultiplicationOp{Num: v})

						// if it's already in there it's at equal or lower step count, don't need to do anything
						n := node{
							stepCount: curr.stepCount + 1,
							ops:       ops,
							val:       newVal,
						}

						nodes[newVal] = &n
						q.PushBack(&n)
					}
				}
			}

			// exponents is a bit annoying, we need to mess with decimals
			if operatorSet.exp {
				for _, v := range sourceSet {
					if v == 1 || v >= 10 {
						// what if we just didn't bother with exponents above 10, that's absurd in
						// the premise of this game
						continue
					}

					newVal, ok := hasIntegerNthRoot(curr.val, v)
					if ok {
						if _, ok := nodes[newVal]; !ok {
							n := node{
								stepCount: curr.stepCount + 1,
								val:       newVal,
							}

							ops := make([]Op, len(curr.ops))
							copy(ops, curr.ops)

							ops = append(ops, ExponentiationOp{Num: v})

							n.ops = ops

							nodes[newVal] = &n
							q.PushBack(&n)
						}
					}
				}
			}
		}
	}

	panic("no solution found")
}

func hasIntegerNthRoot(target, n int) (int, bool) {
	if n == 1 {
		return target, true
	}

	if n < 1 {
		return 0, false
	}

	if target == 1 {
		return 1, true
	}

	if target < 1 {
		return 0, false
	}

	// actually why don't we do a binary search

	// check if lo itself is the answer

	// hi for sure cannot be the answer

	lo, hi := 2, target

	if big.NewInt(int64(lo)).Exp(big.NewInt(int64(lo)), big.NewInt(int64(n)), nil).Int64() == int64(target) {
		return lo, true
	}

	for lo < hi {
		mid := (lo + hi) / 2

		res := big.NewInt(int64(mid)).Exp(big.NewInt(int64(mid)), big.NewInt(int64(n)), nil)

		resInt := int(res.Int64())

		if resInt == target {
			return mid, true
		} else if resInt > target {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}

	if lo == hi {
		if big.NewInt(int64(lo)).Exp(big.NewInt(int64(lo)), big.NewInt(int64(n)), nil).Int64() == int64(target) {
			return lo, true
		}
	}

	return 0, false
}
