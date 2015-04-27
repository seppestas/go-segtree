package segtree

import "fmt"

func ExampleTree_QueryIndex() {
	tree := new(Tree)
	tree.Push(1, 10, "hello, world")
	tree.BuildTree()

	results, err := tree.QueryIndex(4)
	if err != nil {
		panic(fmt.Sprintf("Failed to query tree: %s", err.Error()))
	}

	result := <-results
	fmt.Println("Found:", result)

	// Output: Found: hello, world
}

func ExampleTree_QueryIndex_multiple_elements() {
	tree := new(Tree)
	tree.Push(1, 10, "hello, world")
	tree.Push(5, 15, 3.14)
	tree.BuildTree()

	results, err := tree.QueryIndex(6)
	if err != nil {
		panic(fmt.Sprintf("Failed to query tree: %s", err.Error()))
	}

	for result := range results {
		fmt.Println("Found:", result)
	}

	// Output:
	// Found: 3.14
	// Found: hello, world
}

func ExampleTree_QueryIndex_2_dimensional() {
	inner := new(Tree)
	inner.Push(1, 10, "hello, world")
	inner.BuildTree()

	outer := new(Tree)
	outer.Push(0, 99, inner)
	outer.BuildTree()

	resultsOuter, err := outer.QueryIndex(10)
	if err != nil {
		panic(fmt.Sprintf("Failed to query outer tree: %s", err.Error()))
	}

	result := <-resultsOuter
	resultsInner, err := result.(*Tree).QueryIndex(4)
	if err != nil {
		panic(fmt.Sprintf("Failed to query inner tree: %s", err.Error()))
	}

	result = <-resultsInner
	fmt.Println("Found:", result.(string))

	// Output: Found: hello, world
}

func ExampleTree_Clear() {
	tree := new(Tree)
	tree.Push(1, 10, "hello, world")
	tree.BuildTree()

	tree.Clear()

	_, err := tree.QueryIndex(4)
	if err != nil {
		// The tree is empty. Trying to Query it will fail
		fmt.Println(fmt.Sprintf("Failed to query tree: %s", err.Error()))
	}

	tree.Push(1, 10, "I destroyed the world")
	tree.BuildTree()

	results, err := tree.QueryIndex(4)
	if err != nil {
		panic(fmt.Sprintf("Failed to query tree: %s", err.Error()))
	}

	result := <-results
	// "hello world" will not be found
	fmt.Println("Found:", result)

	// Output:
	// Failed to query tree: Tree is empty. Build the tree first
	// Found: I destroyed the world
}

func ExampleTree_Print() {
	tree := new(Tree)
	tree.Push(1, 10, "hello, world")
	tree.Push(5, 6, "how are you today?")
	tree.Push(9, 45, "test")
	tree.BuildTree()

	tree.Print()

	// Output is a pretty tree (note that the leafs are not always placed properly)
}
