package segtree

import "testing"

func TestBuildingEmptyTree(t *testing.T) {
	tree := new(Tree)

	if err := tree.BuildTree(); err == nil {
		t.Error("Building an empty tree did not result in an error")
	}
}

func TestBuildingFilledTree(t *testing.T) {
	tree := new(Tree)
	tree.Push(1, 10, "hello, world")

	if err := tree.BuildTree(); err != nil {
		t.Errorf("Building a tree failed: %s\n", err.Error())
	}
}

func TestClearingTree(t *testing.T) {
	tree := new(Tree)
	tree.Push(1, 10, "hello, world")
	tree.Clear()

	if err := tree.BuildTree(); err == nil {
		t.Fatal("Failed to clear tree, it was possible to build the tree")
	}

	if _, err := tree.QueryIndex(4); err == nil {
		t.Error("Failed to clear tree, it was possible to query the tree")
	}
}
