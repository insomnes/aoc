package solution

// Heap with slice underneath, less function is used to compare elements
// The heap is a complete binary tree, which means:
// - all levels are filled except possibly for the last level
// - the last level is filled from left to right
//
// This binary tree also should satisfy the heap property (we use min-heap here):
// - the smallest element is at the root
// - the value of each node is less than or equal to the values of its children
//
// To implement a heap, we use a slice to store the elements with these rules:
// - Root is at index 0
// - Left child of node at index i is at 2*i+1
// - Right child of node at index i is at 2*i+2
// - Parent of node at index i is at (i-1)/2
// If we have (numbers are in insert order): 2, 0, 7, 4, 3, 6, 9
// The heap will look like this:
// i: 0  1  2  3  4  5  6
// [  0, 2, 6, 4, 3, 7, 9 ]
//	   0
//	  / \
//	 2   6
//	/ \ / \
// 4  3 7  9
// So the root is at index 0, and the children of i0 are at 2*0+1 = i1 and 2*0+2 = i2
// The parent of i1 is at (1-1)/2 = i0, and the parent of i2 is at (2-1)/2 = i0
// The children of i1 are at 2*1+1 = i3 and 2*1+2 = i4
// The parent of i3 is at (3-1)/2 = i1, and the parent of i4 is at (4-1)/2 = i1
// The children of i2 are at 2*2+1 = i5 and 2*2+2 = i6
// The parent of i5 is at (5-1)/2 = i2, and the parent of i6 is at (6-1)/2 = i2

type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

func NewHeap[T any](less func(a, b T) bool, initCapacity int) *Heap[T] {
	return &Heap[T]{data: make([]T, 0, initCapacity), less: less}
}

func (h *Heap[T]) Len() int {
	return len(h.data)
}

// Push adds a new element to the heap (in the end of slice representation),
// then restores the heap property
func (h *Heap[T]) Push(x T) {
	h.data = append(h.data, x)
	h.ascend(len(h.data) - 1)
}

// Pop removes the smallest element from the heap (first in slice ),
// under the hood, it swaps the first and last elements, pops the last element,
// and then restores the heap property by comparing new pseudo-root with children
func (h *Heap[T]) Pop() T {
	if len(h.data) == 0 {
		panic("heap is empty")
	}
	h.swap(0, len(h.data)-1)
	minVal := h.data[len(h.data)-1]
	h.data = h.data[:len(h.data)-1]
	h.descend(0)
	return minVal
}

// Just check the top element without removing it
func (h *Heap[T]) Peek() T {
	if len(h.data) == 0 {
		panic("heap is empty")
	}
	return h.data[0]
}

// Swap child index with parent index until the heap property is restored
// For example, if we have after pushing 2 to 0,4,7:
// [ 0, 4, 7, 2 ]
// Push will call ascend(3) to restore the heap property,
// and first we will check index parent:
// i = 3 => parent = (3-1)/2 = 1
// 2 < 4 => swap 2 and 4
// New i = 1 (current position of 2, previous position of 4),
// so the new parent index for 1 is (1-1)/2 = 0
// 2 > 0 => stop (by !h.less(h.data[1], h.data[0]))
func (h *Heap[T]) ascend(index int) {
	for {
		parent := (index - 1) / 2
		if index == 0 || !h.less(h.data[index], h.data[parent]) {
			break
		}
		h.swap(index, parent)
		index = parent
	}
}

// Swap parent index with the smallest child index until the heap property is restored
// For example, let's say we have original slice: [ 0, 2, 6, 4, 3, 7, 9 ]
// After popping the root, we swap 0 and 9 and remove 0 (at index 6):
// [ 9, 2, 6, 4, 3, 7 ]
// Pop will call descend(0) to restore the heap property,
// and first we will check index parent:
// i = 0 => left = 2*0+1 = 1, right = 2*0+2 = 2
// Values lv = 2, rv = 6, rv = 9
// 2 < 9 => smallest should be set to left index:
// smallest = 1, left = 1, right = 2
// 6 > 2, do not change smallest
// smallest index (1) != index (0), so we swap 0 and 1 and set index to 1
// [ 2, 9, 6, 4, 3, 7 ]
// i = 1 => left = 2*1+1 = 3, right = 2*1+2 = 4
// Values lv = 4, rv = 3, sv = 9
// 4 < 9 => smallest should be set to left index:
// smallest = 3, left = 3, right = 4
// 3 < 4 => smallest should be set to right index:
// smallest = 4, left = 3, right = 4
// smallest index (4) != index (1), so we swap 1 and 4 and set index to 4
// [ 2, 3, 6, 4, 9, 7 ]
// i = 4 => left = 2*4+1 = 9, right = 2*4+2 = 10
// meaning there are no children, check that smallest index (4) == index (4) and stop
func (h *Heap[T]) descend(index int) {
	size := len(h.data)
	for {
		leftIdx, rightIdx := 2*index+1, 2*index+2
		smallestIdx := index

		if leftIdx < size && h.less(h.data[leftIdx], h.data[smallestIdx]) {
			smallestIdx = leftIdx
		}
		if rightIdx < size && h.less(h.data[rightIdx], h.data[smallestIdx]) {
			smallestIdx = rightIdx
		}
		if smallestIdx == index {
			break
		}
		h.swap(index, smallestIdx)
		index = smallestIdx
	}
}

func (h *Heap[T]) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}
