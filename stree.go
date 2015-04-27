// Package segmentTree implements a segment tree to solve stabbing problems
// This allows to store and retrieve elements indexed by a range.

package segtree

import (
	"errors"
	"sort"
)

const (
	// Inf is defined as the max value of an integer, used as +∞
	Inf = int(^uint(0) >> 1)
	// NegInf is defined as the min value of an integer, used as -∞
	NegInf = -Inf - 1
)

type interval struct {
	segment
	element interface{}
}

type segment struct {
	from int
	to   int
}

type node struct {
	segment     segment
	left, right *node
	intervals   []*interval
}

type Tree struct {
	base []interval
	root *node
}

// Push pushes an interval to the interval stack
func (t *Tree) Push(from, to int, element interface{}) {
	if to < from {
		from, to = to, from
	}
	t.base = append(t.base, interval{segment{from, to}, element})
}

// Clear clears the interval stack
func (t *Tree) Clear() {
	t.root = nil
	t.base = nil
}

// BuildTree builds the segment tree from the interval stack
func (t *Tree) BuildTree() error {
	if len(t.base) == 0 {
		return errors.New("No intervals in stack. Push intervals on the stack of the tree first.")
	}

	leaves := elementaryIntervals(t.endpoints())
	t.root = t.insertNodes(leaves)

	for i := range t.base {
		t.root.insertInterval(&t.base[i])
	}

	return nil
}

// Removes duplicate entries from a sorted slice
func removedups(sorted []int) (unique []int) {
	unique = make([]int, 0, len(sorted))
	unique = append(unique, sorted[0])
	prev := sorted[0]
	for _, val := range sorted[1:] {
		if val != prev {
			unique = append(unique, val)
			prev = val
		}
	}
	return
}

// Creates a sorted slice of unique endpoints from a tree's base
func (t *Tree) endpoints() []int {
	baseLen := len(t.base)
	endpoints := make([]int, baseLen*2)

	// When there are a lot of intervals, there is a big chance of big overlaps
	// Try to have the endpoints sorted as much as possible when putting them
	// in the slice to reduce the final sort time.
	// endpoints[0] = NegInf
	for i, interval := range t.base {
		endpoints[i] = interval.from
		endpoints[i+baseLen] = interval.to
	}
	// endpoints[baseLen*2+1] = Inf

	sort.Sort(sort.IntSlice(endpoints))

	return removedups(endpoints)
}

// Creates a slice of elementary intervals from a slice of (sorted) endpoints
// Input: [p1, p2, ..., pn]
// Output: [{p1 : p1}, {p1 : p2}, {p2 : p2},... , {pn : pn}
func elementaryIntervals(endpoints []int) []segment {
	if len(endpoints) == 1 {
		return []segment{segment{endpoints[0], endpoints[0]}}
	}

	intervals := make([]segment, len(endpoints)*2-1)

	for i := 0; i < len(endpoints); i++ {
		intervals[i*2] = segment{endpoints[i], endpoints[i]}
		if i < len(endpoints)-1 {
			intervals[i*2+1] = segment{endpoints[i], endpoints[i+1]}
		}
	}
	return intervals
}

// insertNodes builds the tree structure from the elementary intervals
func (t *Tree) insertNodes(leaves []segment) *node {
	var n *node
	if len(leaves) == 1 {
		n = &node{segment: leaves[0]}
		n.left = nil
		n.right = nil
	} else {
		n = &node{segment: segment{leaves[0].from, leaves[len(leaves)-1].to}}
		center := len(leaves) / 2
		n.left = t.insertNodes(leaves[:center])
		n.right = t.insertNodes(leaves[center:])
	}

	return n
}

func (s *segment) subsetOf(other *segment) bool {
	return other.from <= s.from && other.to >= s.to
}

func (s *segment) intersectsWith(other *segment) bool {
	return other.from <= s.to && s.from <= other.to ||
		s.from <= other.to && other.from <= s.to
}

// Inserts interval into given tree structure
func (n *node) insertInterval(i *interval) {
	if n.segment.subsetOf(&i.segment) {
		if n.intervals == nil {
			n.intervals = make([]*interval, 0, 1)
		}
		n.intervals = append(n.intervals, i)
	} else {
		if n.left != nil && n.left.segment.intersectsWith(&i.segment) {
			n.left.insertInterval(i)
		}
		if n.right != nil && n.right.segment.intersectsWith(&i.segment) {
			n.right.insertInterval(i)
		}
	}
}

// QueryIndex looks for all segments in the interval tree that contain
// a given index. The elements associated with the segments will be sent
// on the returned channel. No element will be sent twice.
// The elements will not be sent in any specific order.
func (t *Tree) QueryIndex(index int) (<-chan interface{}, error) {
	if t.root == nil {
		return nil, errors.New("Tree is empty. Build the tree first")
	}

	intervals := make(chan *interval)

	go func(t *Tree, index int, intervals chan *interval) {
		query(t.root, index, intervals)
		close(intervals)
	}(t, index, intervals)

	elements := make(chan interface{})

	go func(intervals chan *interval, elements chan interface{}) {
		results := make(map[interface{}]struct{})
		for interval := range intervals {
			_, alreadyFound := results[interval.element]
			if !alreadyFound {
				// Store an empty struct in the map to minimize memory footprint
				results[interval.element] = struct{}{}
				elements <- interval.element
			}
		}
		close(elements)
	}(intervals, elements)

	return elements, nil
}

func (s segment) contains(index int) bool {
	return s.from <= index && index <= s.to
}

func query(node *node, index int, results chan<- *interval) {
	if node.segment.contains(index) {
		for _, interval := range node.intervals {
			results <- interval
		}
		if node.left != nil {
			query(node.left, index, results)
		}
		if node.right != nil {
			query(node.right, index, results)
		}
	}

}
