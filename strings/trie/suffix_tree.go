package trie

import (
	"fmt"
	"sort"
)

/*
Suffix tree is a compressed trie containing all the suffixes of the given text as their keys and indices in the text as their values.

Using Ukkonen’s algorithm, the construction of such a tree for the string S takes time and space linear in the length of S.
Suffix trees also provide one of the first linear-time solutions for the longest common substring problem.

The suffix tree for the string S of length N is defined as a tree such that:
1. The tree has exactly n leaves numbered from 1 to N.
2. Except for the root, every internal node has at least two children.
3. Each edge is labelled with a non-empty substring of S.
4. No two edges starting out of a node can have string-labels beginning with the same character.
5. The string obtained by concatenating all the string-labels found on the path from the root to leaf i spells out suffix S[i..n] for i from 1 to n.

# of leaf nodes = M = the number of suffixes of a given string
# of non-leaf nodes <= M-1 (because an internal node of a suffix tree has at least 2 children)
# of total nodes <= 2M-1
So, space time complexity will be O(M).

Ideas:
1. This data structure starts from the idea that it reduce some paths without branching from a node to a leaf node
by each path which has only one edge.
	- Each label of edges, that are compressed, equals the concatenated edge labels from the node to the leaf node.
	e.g., Given a trie like "root --a--> node --b--> node --c--> node", it will be changed as "root --abc--> node"

2. Labling edges as suffixes is not efficient because of storing strings, so it will be better
to label edges as (start, end) pair which means the start index and the end index of each suffix.
	- The information about each edge will be stored in each node.
	e.g., If string S[x:y+1] is "abc", then "root --abc--> node" is gonna be "root --(x,y)--> node"

3. Each leaf node need to be the end of substring(=$) because the tree needs to distinguish overlapping substrings.
	e.g., Given a string "bcbc", substrings "bc" and "bc", where the first "bc" is in front of the second,
	are on the same edge if there is no use of $.

		(Before)
		root - bcbc$
		(After)
		root - bc - bc$
			  - $

How to build: https://www.geeksforgeeks.org/ukkonens-suffix-tree-construction-part-1/
Assume that we have a string S with the length N, where S has "$" at the end.
1. We will have N phases and, in each phase, we will process the prefix S[:i+1] (0 <= i < N).
2. Phase i will proceed at most i extentions, each of which will process the suffix of S[j:i+1] (0 <= j < i+1).
	- But, whenever one character is added, it doesn't have to visit all the characters of each suffix
	because of using information from the previous extentions in the current extention.
	- It will stores information like which node it should start from in each extention.
		- activeNode, activeEdge, and activeLength will be so useful that we can easily find it.
	- It will also stores information like where it should go to process remaining suffixes.
		- remainingSuffixCount, lastNewNode, and suffixLink will be so useful that we can easily move to another node.
	- There is also use a pointer variable to control all the leaf nodes at a time. => leafEnd
		- Once a leaf, always a leaf.
3. When S[j:i+1] is inserted into a suffix tree, we will follow as like below:
	- Extention Rule 1. if there is S[j:i] on the tree and S[i-1] is the last character of this edge, then just add.
		- Addition will be done by increasing the end index of all the leaf nodes by 1. So, it will take O(1).
	- Extention Rule 2. if there is S[j:i] on the tree and S[i] is not the same as the next character of this edge,
	then it will branch that node into two nodes, where one is a existing leaf node(A) representing a substring
	starting from S[i-1] of that edge and another is a new leaf node(B) representing only S[i].
	Also, that branched node will become a new internal node representing S[j:i].
		- In this case, lastNewNode will be changed as a new internal node.
	- Extention Rule 3. if there is S[j:i] on the tree and S[i] is the next character of this edge, there's nothing to do.
		- In this case, remainingSuffixCount won't be decreased.
4. Suffix Links: one edge with path-label xA, where x is a character and A denotes substring, connects root node and
ndoe W via node V and another edge with path-label A connects node X and Y.
The edge going from node V to node X is so-called a suffix link.

	root --x-->  V --A--> W
	             | <<<<<<<<<<<<< this is a suffix link
	      ...  - X --A--> Y

If A is empty string, then a suffix link of that node will go to root node.
This can be useful to find each internal node to process the next suffix by O(1).
5. Keep in mind: all the subtrees of internal nodes connecting with suffix links are the same tree structure.
	- each edge label of one subtree is the same as the label of another that the subtree is linking.
*/

type suffixTreeNode struct {
	start      int
	end        *int
	suffixLink *suffixTreeNode
	children   map[int]*suffixTreeNode
}

type SuffixTree struct {
	/* lastNewNode will point to newly created internal node, waiting for it's suffix link to be set, which might get
	a new suffix link (other than root) in next extension of same phase. lastNewNode will be set to NULL when last
	newly created internal node (if there is any) got it's suffix link reset to new internal node created in next
	extension of same phase. */
	root                 *suffixTreeNode
	lastNewNode          *suffixTreeNode
	activeNode           *suffixTreeNode
	activeEdge           int // Character index to be on the next node = the first letter of an edge to visit the next suffix
	activeLength         int
	remainingSuffixCount int // How many suffixes yet to be added in tree
	leafEnd              int
	count                int
	size                 int // Length of input string
	text                 string
}

func (this *suffixTreeNode) getSortedKeys() []int {
	keys := make([]int, 0, len(this.children))
	for k := range this.children {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

/* Returns the number of characters of a given edge label. */
func (this *suffixTreeNode) edgeLength() int {
	return (*this.end) - this.start + 1
}

func (this *SuffixTree) newNode(start int, end *int) *suffixTreeNode {
	this.count++
	node := &suffixTreeNode{
		start:      start,
		end:        end,
		suffixLink: this.root,
		children:   map[int]*suffixTreeNode{},
	}
	return node
}

/* activePoint change for walk down (APCFWD) using Skip/Count Trick (Trick 1).
If activeLength is greater than current edge length, set next internal node as activeNode
and adjust activeEdge and activeLength accordingly to represent same activePoint */
func (this *SuffixTree) walkDown(current *suffixTreeNode) bool {
	edgeLen := current.edgeLength()
	if this.activeLength >= edgeLen {
		this.activeEdge = int(this.text[(*current.end)+1])
		this.activeLength -= edgeLen
		this.activeNode = current
		return true
	}
	return false
}

func (this *SuffixTree) extend(pos int) {
	// Extension Rule 1, this takes care of extending all leaves created so far in tree.
	this.leafEnd = pos
	// Indicates that a new suffix added to the list of suffixes yet to be added in tree.
	this.remainingSuffixCount++
	// While starting a new phase, indicates there is no internal node waiting for it's suffix link reset in current phase.
	this.lastNewNode = nil

	// Add all suffixes (yet to be added) one by one in tree.
	for this.remainingSuffixCount > 0 {
		if this.activeLength == 0 {
			this.activeEdge = int(this.text[pos])
		}
		if this.activeNode.children[this.activeEdge] == nil {
			// If there is no outgoing edge starting with activeEdge from activeNode,
			// Extension Rule 2 (A new leaf edge gets created)
			this.activeNode.children[this.activeEdge] = this.newNode(pos, &this.leafEnd)

			// A new leaf edge is created in above line starting from an existing node (the current activeNode),
			// and if there is any internal node waiting for it's suffix link get reset,
			// point the suffix link from that last internal node to current activeNode.
			// Then set lastNewNode to NULL indicating no more node waiting for suffix link reset.
			if this.lastNewNode != nil {
				this.lastNewNode.suffixLink = this.activeNode
				this.lastNewNode = nil
			}
		} else {
			// If there is an outgoing edge starting with activeEdge from activeNode,
			// Get the next node at the end of edge starting with activeEdge.
			next := this.activeNode.children[this.activeEdge]
			if this.walkDown(next) {
				// If it is walked down, start from next node (the new activeNode).
				continue
			}

			// Extension Rule 3 (current character being processed is already on the edge)
			if this.text[next.start+this.activeLength] == this.text[pos] {
				// If a newly created node waiting for it's suffix link to be set,
				// then set suffix link of that waiting node to current active node.
				if this.lastNewNode != nil && this.activeNode != this.root {
					this.lastNewNode.suffixLink = this.activeNode
					this.lastNewNode = nil
				}
				this.activeLength++
				// Stop all further processing in this phase and move on to next phase.
				break
			}

			// Here is when activePoint is in middle of the edge being traversed
			// and current character being processed is not on the edge (we fall off the tree).
			// In this case, it has to add a new internal node and a new leaf edge going out of that new node.
			// This is Extension Rule 2, where a new leaf edge and a new internal node get created.
			splitEnd := new(int)
			*splitEnd = next.start + this.activeLength - 1
			internal := this.newNode(next.start, splitEnd)
			this.activeNode.children[this.activeEdge] = internal
			// New leaf coming out of new internal node
			internal.children[int(this.text[pos])] = this.newNode(pos, &this.leafEnd)
			// Existing leaf node out of new internal node
			next.start += this.activeLength
			internal.children[int(this.text[next.start])] = next

			// If there is any internal node created in last extensions of same phase
			// which is still waiting for it's suffix link reset, do it now.
			if this.lastNewNode != nil {
				this.lastNewNode.suffixLink = internal
			}

			// Make the current newly created internal node waiting for it's suffix link reset
			// (which is pointing to root at present).
			// If we come across any other internal node (existing or newly created) in next extension of same phase
			// when a new leaf edge gets added (i.e. when Extension Rule 2 applies is any of the next extension of
			// same phase) at that point, suffixLink of this node will point to that internal node.
			this.lastNewNode = internal
		}

		// One suffix got added in tree.
		this.remainingSuffixCount--
		if this.activeNode == this.root && this.activeLength > 0 {
			this.activeLength--
			this.activeEdge = int(this.text[pos-this.remainingSuffixCount+1])
		} else if this.activeNode != this.root {
			this.activeNode = this.activeNode.suffixLink
		}
	}
}

func (this *SuffixTree) FreeSuffixTreeByPostOrder(node *suffixTreeNode) {
	if node == nil {
		return
	}
	keys := node.getSortedKeys()
	for _, key := range keys {
		if node.children[key] != nil {
			this.FreeSuffixTreeByPostOrder(node.children[key])
			node.children[key] = nil
		}
	}
}

func (this *SuffixTree) PrintPretty(node *suffixTreeNode, spaces int) {
	s := "root "
	if node.start != -1 {
		s = fmt.Sprintf("--(%s)--> [%d,%d] ", this.text[node.start:*node.end+1], node.start, *node.end)
	}
	fmt.Print(s)
	keys := node.getSortedKeys()
	if len(keys) == 0 {
		fmt.Println()
	}
	for i, key := range keys {
		if node.children[key] != nil {
			if i == 0 {
				this.PrintPretty(node.children[key], spaces+len(s))
			} else {
				for j := 0; j < spaces+len(s); j++ {
					fmt.Print(" ")
				}
				this.PrintPretty(node.children[key], spaces+len(s))
			}
		}
	}
}

func (this *SuffixTree) Count() int {
	return this.count
}

func (this *SuffixTree) Root() *suffixTreeNode {
	return this.root
}

func NewSuffixTree(s string) *SuffixTree {
	tree := &SuffixTree{activeEdge: -1, leafEnd: -1, size: len(s), text: s}
	rootEnd := new(int)
	*rootEnd = -1
	tree.root = tree.newNode(-1, rootEnd)
	tree.activeNode = tree.root
	for i := 0; i < tree.size; i++ {
		tree.extend(i)
	}
	return tree
}
