package segtree

import "testing"

func TestFindingSingleElement(t *testing.T) {
	tree := new(Tree)
	test := "hello, world"
	tree.Push(1, 10, test)
	tree.BuildTree()

	results, err := tree.QueryIndex(4)
	if err != nil {
		queryFailed(t, err)
	}

	result := <-results

	if result != test {
		wrongElement(t, result, test)
	}

	if _, ok := <-results; ok != false {
		toManyElements(t)
	}
}

func TestFindingElementSizeZeroRange(t *testing.T) {
	tree := new(Tree)
	test := "hello, world"
	tree.Push(1, 1, test)
	tree.BuildTree()

	results, err := tree.QueryIndex(1)
	if err != nil {
		queryFailed(t, err)
	}

	result := <-results

	if result != test {
		wrongElement(t, result, test)
	}

	if _, ok := <-results; ok != false {
		toManyElements(t)
	}
}

func TestFindingElementPseudoEndlessRange(t *testing.T) {
	tree := new(Tree)
	test := "hello, world"
	tree.Push(1, Inf, test)
	tree.BuildTree()

	results, err := tree.QueryIndex(9999)
	if err != nil {
		queryFailed(t, err)
	}

	result := <-results

	if result != test {
		wrongElement(t, result, test)
	}

	if _, ok := <-results; ok != false {
		toManyElements(t)
	}
}

func find(element interface{}, elements []interface{}) (bool, int) {
	for i, e := range elements {
		if element == e {
			return true, i
		}
	}
	return false, -1
}

func TestFindingMultipleElements(t *testing.T) {
	tree := new(Tree)
	tests := make([]interface{}, 5)
	tests[0] = "one"
	tests[1] = "two"
	tests[2] = 3.14
	tests[3] = 4
	tests[4] = "stuff"

	for i, test := range tests {
		tree.Push(10+i, 100+i, test)
	}

	tree.BuildTree()
	results, err := tree.QueryIndex(15)

	if err != nil {
		queryFailed(t, err)
	}

	for result := range results {
		found, index := find(result, tests)
		if !found {
			receivedUnexpected(t, result)
		}
		// Remove element from tests
		tests[index], tests = tests[len(tests)-1], tests[:len(tests)-1]
	}

	if len(tests) > 0 {
		for _, test := range tests {
			t.Errorf("Did not find %v\n", test)
		}
	}

	if _, ok := <-results; ok != false {

	}
}

func TestFindingOverlappingElements(t *testing.T) {
	tree := new(Tree)
	tests := make([]interface{}, 2)
	tests[0] = "one"
	tests[1] = "two"
	tree.Push(1, 10, tests[0])
	tree.Push(5, 15, tests[1])
	tree.BuildTree()

	// Index only in first range
	results, err := tree.QueryIndex(4)
	if err != nil {
		queryFailed(t, err)
	}

	result := <-results
	if result != tests[0] {
		wrongElement(t, result, tests[0])
	}

	if _, ok := <-results; ok != false {
		toManyElements(t)
	}

	// Index only in second range
	results, err = tree.QueryIndex(11)
	if err != nil {
		queryFailed(t, err)
	}

	result = <-results
	if result != tests[1] {
		wrongElement(t, result, tests[1])
	}

	if _, ok := <-results; ok != false {
		toManyElements(t)
	}

	// Index in both ranges
	results, err = tree.QueryIndex(6)
	if err != nil {
		queryFailed(t, err)
	}

	for result := range results {
		found, index := find(result, tests)
		if !found {
			receivedUnexpected(t, result)
		}
		// Remove element from tests
		tests[index], tests = tests[len(tests)-1], tests[:len(tests)-1]

	}

	if len(tests) > 0 {
		for _, test := range tests {
			t.Errorf("Did not find %v\n", test)
		}
	}

	if _, ok := <-results; ok != false {
		toManyElements(t)
	}
}

func TestOutOfRangeNotFound(t *testing.T) {
	tree := new(Tree)
	test := "hello, world"

	tree.Push(1, 10, test)

	tree.BuildTree()
	results, err := tree.QueryIndex(20)

	if err != nil {
		queryFailed(t, err)
	}

	if result, ok := <-results; ok != false {
		receivedUnexpected(t, result)
	}
}

func TestAddingReverseDirection(t *testing.T) {
	tree := new(Tree)
	test := "hello, world"
	tree.Push(10, 1, test)
	tree.BuildTree()

	results, err := tree.QueryIndex(4)
	if err != nil {
		queryFailed(t, err)
	}

	result := <-results

	if result != test {
		wrongElement(t, result, test)
	}
}

func TestFindingSameElementTwice(t *testing.T) {
	tree := new(Tree)
	test := "hello, world"

	tree.Push(1, 10, test)
	tree.Push(5, 15, test)

	tree.BuildTree()
	results, err := tree.QueryIndex(10)

	if err != nil {
		queryFailed(t, err)
	}

	result := <-results

	if result != test {
		wrongElement(t, result, test)
	}

	// `test` should only be sent once
	if _, ok := <-results; ok != false {
		toManyElements(t)
	}
}

// Some frequently used errors

func queryFailed(t *testing.T, err error) {
	t.Fatalf("Failed to query tree: %s\n", err.Error())
}

func wrongElement(t *testing.T, result, expected interface{}) {
	t.Errorf("Found wrong element %v, expected %v\n", result, expected)
}

func toManyElements(t *testing.T) {
	t.Error("To many elements sent on channel")
}

func receivedUnexpected(t *testing.T, result interface{}) {
	t.Fatalf("Received unexpected %v\n", result)
}
