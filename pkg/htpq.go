/*
Copyright © 2023 Martin Coder <martincoder1@gmail.com>

Use of this source code is governed by an MIT-style
license that can be found in the LICENSE file or at
https://opensource.org/licenses/MIT.
*/

package huffmyfile

import (
	"fmt"
)

type HTWrapper struct {
	ht *HuffTree
	// priority int
	// // The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

//PriorityQueue implemented using Heap interface in order to always be able to access the two trees with
//the lowest total frequencies of their leaves
type PriorityQueue []*HTWrapper

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].ht.root.freq > pq[j].ht.root.freq
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	htw := x.(*HTWrapper)
	htw.index = n
	*pq = append(*pq, htw)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	htw := old[n-1]
	old[n-1] = nil // avoid memory leak
	htw.index = -1 // for safety
	*pq = old[0 : n-1]
	return htw
}

// update modifies the priority and value of an Item in the queue.
// func (pq *PriorityQueue) update(htw *HTWrapper, ht *HuffTree) {
// 	htw.ht = ht
// 	heap.Fix(pq, htw.index)
// }

//Prints an array representation of the priority queue's contents
func (pq *PriorityQueue) Print() {
	fmt.Print("[")
	for i, s := range *pq {

		fmt.Printf("%d: ", i)
		s.ht.Print()
		fmt.Print(", ")
	}
	fmt.Print("]")
}

func (pq *PriorityQueue) Println() {
	pq.Print()
	fmt.Println()
}
