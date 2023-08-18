/*
Copyright © 2023 Martin Coder <martincoder1@gmail.com>

Use of this source code is governed by an MIT-style
license that can be found in the LICENSE file or at
https://opensource.org/licenses/MIT.
*/

package huffmyfile

import (
	"container/heap"
	"fmt"
)

type HuffNode struct {
	asciiVal int
	freq     int
	left     *HuffNode
	right    *HuffNode
}
type HuffTree struct {
	root *HuffNode
}

func (a *HuffTree) MakeHuffmanTree(freqMap map[int]int) {
	pq := make(PriorityQueue, 0)

	for k, v := range freqMap {
		//Fills priority queue with individual, single-node huffman trees for each character.
		r := HuffNode{asciiVal: k, freq: v}
		h := HuffTree{root: &r}
		htw := HTWrapper{ht: &h}
		heap.Push(&pq, &htw)
	}

	// fmt.Print("\nIndividual trees, prior to consolidation: ")
	pq.Println()

	//Removes two smallest (lowest total frequency) trees, combines them, pushes it back onto queue.
	//Repeats until there is one large Huffman tree with each character as a leaf.
	for pq.Len() > 1 {
		hta := pq.Pop().(*HTWrapper).ht
		htb := pq.Pop().(*HTWrapper).ht
		htw := HTWrapper{ht: hta.Combine(htb)}
		heap.Push(&pq, &htw)
	}
	heap.Init(&pq)

	a.root = new(HuffNode)
	*a.root = *pq[0].ht.root
}

//Returns a new HuffTree with a new root node pointing at the two which were combined.
//The left child has lower frequency, and the freq value of the root node is the sum of its children.
func (a *HuffTree) Combine(b *HuffTree) *HuffTree {

	r := HuffNode{freq: a.root.freq + b.root.freq}
	h := HuffTree{root: &r}

	if a.Compare(b) < 0 {
		h.root.left = a.root
		h.root.right = b.root
	} else if a.Compare(b) > 0 {
		h.root.left = b.root
		h.root.right = a.root
	} else { //tie-breaker
		if a.root.asciiVal <= b.root.asciiVal {
			h.root.left = a.root
			h.root.right = b.root
		} else {
			h.root.left = b.root
			h.root.right = a.root
		}
	}
	return &h
}

func (a *HuffTree) Compare(b *HuffTree) int {
	if a.root.freq > b.root.freq {
		return 1
	}
	if a.root.freq < b.root.freq {
		return -1
	}
	return 0
}

//Returns a map of the unique codes which each character will be represented by in a particular encoding.
func (ht *HuffTree) CodeMap() map[int]string {
	code := ""
	cm := make(map[int]string)

	ht.root.generateCodes(code, cm)
	return cm
}

func (r *HuffNode) generateCodes(code string, codeMap map[int]string) {

	if r.left == nil && r.right == nil {
		codeMap[r.asciiVal] = code
		return
	}
	if r.left != nil {
		r.left.generateCodes(code+"0", codeMap)
	}
	if r.right != nil {
		r.right.generateCodes(code+"1", codeMap)
	}

}

//Prints out a pre-order representation of the HuffmanTree. Characters are printed but
//only their ascii values are stored in the nodes.
func (a *HuffTree) Print() {
	printRec(a.root)
}

func printRec(n *HuffNode) {
	if n.left != nil {
		printRec(n.left)
	}
	fmt.Print("{")
	av := n.asciiVal
	if av == 0 {
		fmt.Print("_")
	} else if av == 10 {
		fmt.Print("'\\n'")
	} else {
		fmt.Print("'", string(rune(n.asciiVal)), "'")
	}
	fmt.Print(", ", n.freq, "} ")
	if n.right != nil {
		printRec(n.right)
	}
}
