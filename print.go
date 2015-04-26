package segtree

import "fmt"

func (n *node) print() {
	from := fmt.Sprintf("%d", n.segment.from)
	switch n.segment.from {
	case Inf:
		from = "+∞"
	case NegInf:
		from = "-∞"
	}
	to := fmt.Sprintf("%d", n.segment.to)
	switch n.segment.to {
	case Inf:
		to = "Inf"
	case NegInf:
		to = "NegInf"
	}
	fmt.Printf("(%s,%s)", from, to)
	if n.intervals != nil {
		fmt.Print("->[")
		for _, intrvl := range n.intervals {
			fmt.Printf("(%v,%v)=[%v]", intrvl.from, intrvl.to, intrvl.element)
		}
		fmt.Print("]")
	}

}

// Traverse tree recursively call enter when entering node, resp. leave
func traverse(node *node, depth int, enter, leave func(*node, int)) {
	if node == nil {
		return
	}
	if enter != nil {
		enter(node, depth)
	}
	traverse(node.left, depth+1, enter, leave)
	traverse(node.right, depth+1, enter, leave)
	if leave != nil {
		leave(node, depth)
	}
}

// Returs log with base 2 of an int.
func log2(num int) int {
	if num == 0 {
		return NegInf
	}
	i := -1
	for num > 0 {
		num = num >> 1
		i++
	}
	return i
}

func space(n int) {
	for i := 0; i < n; i++ {
		fmt.Print(" ")
	}
}

// Print prints a binary tree recursively to sdout
func (t *Tree) Print() {
	endpoints := len(t.base)*2 + 2
	leaves := endpoints*2 - 3
	height := 1 + log2(leaves)

	fmt.Println("Height:", height, ", leaves:", leaves)
	levels := make([][]*node, height+1)

	traverse(t.root, 0, func(n *node, depth int) {
		levels[depth] = append(levels[depth], n)
	}, nil)

	for i, level := range levels {
		for j, n := range level {
			space(12 * (len(levels) - 1 - i))
			n.print()
			space(1 * (height - i))

			if j-1%2 == 0 {
				space(2)
			}
		}
		fmt.Println()
	}
}
