package main

import (
	"fmt"

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

func findMostConvenientMadeUp(target int, sourceSet []int) []Op {
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

	for q.Len() > 0 {
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
			// divide, add, subtract
			if curr.val%v == 0 {
				// divide
				newVal := curr.val / v
				if _, ok := nodes[newVal]; !ok {
					// if it's already in there it's at equal or lower step count, don't need to do anything
					n := node{
						stepCount: curr.stepCount + 1,
						ops:       append(curr.ops, MultiplicationOp{Num: v}),
						val:       newVal,
					}

					nodes[newVal] = &n
					q.PushBack(&n)
				}
			}

			// add and subtract
			newVal := curr.val - v
			if newVal > 0 {
				if _, ok := nodes[newVal]; !ok {
					n := node{
						stepCount: curr.stepCount + 1,
						ops:       append(curr.ops, AdditionOp{Num: v}),
						val:       newVal,
					}

					nodes[newVal] = &n
					q.PushBack(&n)
				}
			}

			newVal = curr.val + v
			if _, ok := nodes[newVal]; !ok {
				n := node{
					stepCount: curr.stepCount + 1,
					ops:       append(curr.ops, SubtractionOp{Num: v}),
					val:       newVal,
				}

				nodes[newVal] = &n
				q.PushBack(&n)
			}
		}
	}

	panic("no solution found")
}
