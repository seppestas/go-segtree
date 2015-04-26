# Go segtree

A [segment tree](http://en.wikipedia.org/wiki/Segment_tree) implementation written in Go.

This library allows storing and retrieving elements indexed by a range. It does this by storing the elements in a segment tree.

Based on:
- [go-stree](https://github.com/toberndo/go-stree) by @toberndo
- Chapter 10.3 of [Computational Geometry: Algorithms and Applications](http://www.cs.uu.nl/geobook/) rev. 3 by Mark de Berg, Otfried Cheong, Marc van Kreveld and Mark Overmars (ISBN 978-3-540-77973-5)

The elements are sent to a channel as soon as they are found in the tree. This allows efficient querying of e.g multi-dimensional trees (trees containing trees). The elements are not sent in any specific order, however each found element will only be sent once.

# Example usages:

```go
tree := new(segtree.Tree)
tree.Push(1, 10, "hello, world")
tree.BuildTree()

results, err := tree.QueryIndex(4)
if err != nil {
	panic(fmt.Sprintf("Failed to query tree: %s", err.Error()))
}

result := <-results
fmt.Println("Found:", result)
```

An exemplary 2-dimensional segment tree as a tree in a tree:
```go
inner := new(segtree.Tree)
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
```

The library also allows to pretty print the content of the tree for easy debugging.

Example:
```go
tree := new(segtree.Tree)
tree.Push(1, 10, "hello, world")
tree.Push(5, 6, "how are you today?")
tree.Push(9, 45, "test")
tree.BuildTree()

tree.Print()
```

## State of this package

I'm still experimenting with walking the tree concurrently using go routines, but at first sight this does not have any real benefits.
I will commit my benchmarks once I have them properly defined so that performance changes can be tracked consistently.

I'm also still experimenting with "real" multi-dimensional segment trees.

However there are no planned changes to the interface of this package, so it should be safe to use it in your project.
