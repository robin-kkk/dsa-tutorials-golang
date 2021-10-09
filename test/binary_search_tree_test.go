package test

import (
	tree "awesome-golang/binary-tree"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	No   int
	Name string
	Age  int
}

func makeData() map[int]*Person {
	return map[int]*Person{
		8:  &Person{No: 8, Name: "Jack", Age: 12},
		4:  &Person{No: 4, Name: "Dan", Age: 18},
		12: &Person{No: 12, Name: "Lisa", Age: 24},
		2:  &Person{No: 2, Name: "Kim", Age: 9},
		1:  &Person{No: 1, Name: "Lee", Age: 14},
		3:  &Person{No: 3, Name: "Park", Age: 16},
		5:  &Person{No: 5, Name: "Moana", Age: 23},
		6:  &Person{No: 6, Name: "Raven", Age: 35},
		7:  &Person{No: 7, Name: "Liu", Age: 28},
		9:  &Person{No: 9, Name: "Jinny", Age: 27},
		23: &Person{No: 23, Name: "Jane", Age: 22},
		14: &Person{No: 14, Name: "Jim", Age: 30},
		27: &Person{No: 27, Name: "Ken", Age: 29},
	}
}

func TestBinarySearchTree(t *testing.T) {
	data := makeData()
	root := tree.NewBSTNode(8, nil, nil, nil)
	insertedOrder := []int{8, 4, 12, 2, 1, 3, 5, 6, 7, 9, 23, 14, 27}
	for _, key := range insertedOrder {
		root.Add(key, data[key])
	}

	nodes := make([]interface{}, 0)
	tree.InOrder(root, &nodes)
	ans := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 14, 23, 27}
	for i, node := range nodes {
		assert.Equal(t, data[ans[i]], node.(*Person))
		result := tree.Search(root, ans[i])
		assert.Equal(t, data[ans[i]].Name, result.(*Person).Name)
	}

	removed := []int{6, 27, 23, 4}
	for _, target := range removed {
		root = root.Remove(target).(*tree.BSTNode)
		assert.Equal(t, nil, tree.Search(root, target))
	}

	nodes = make([]interface{}, 0)
	tree.InOrder(root, &nodes)
	ans = []int{1, 2, 3, 5, 7, 8, 9, 12, 14}
	for i, node := range nodes {
		assert.Equal(t, ans[i], node.(*Person).No)
	}

	// Remove a root.
	root = root.Remove(8).(*tree.BSTNode)
	nodes = make([]interface{}, 0)
	tree.InOrder(root, &nodes)
	ans = []int{1, 2, 3, 5, 7, 9, 12, 14}
	for i, node := range nodes {
		assert.Equal(t, ans[i], node.(*Person).No)
	}

	// Check if it is a binary search tree.
	assert.Equal(t, true, tree.IsBST(root, 1, 14))
	assert.Equal(t, false, tree.IsBST(root, 2, 13))

	// Search interatively.
	ans = []int{3, 5, 14}
	for _, key := range ans {
		result := tree.SearchIteratively(root, key)
		assert.Equal(t, data[key].Name, result.(*Person).Name)
	}
	assert.Equal(t, nil, tree.SearchIteratively(root, 8))

	// Trim nodes out of [3, 10].
	root.Add(10, &Person{No: 10, Name: "Jay", Age: 24})
	root.Add(11, &Person{No: 11, Name: "Ray", Age: 31})
	root = tree.TrimBST(root, 3, 10)
	nodes = make([]interface{}, 0)
	tree.InOrder(root, &nodes)
	ans = []int{3, 5, 7, 9, 10}
	for i, node := range nodes {
		assert.Equal(t, ans[i], node.(*Person).No)
	}
}
