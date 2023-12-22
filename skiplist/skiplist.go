package skiplist

import (
	"gldb/utils"
	"math/rand"
	"sync"
)

const kMaxHeight = 12
const kBranching = 4

type SkipList struct {
	head       *Node
	maxHeight  int
	comparator utils.Comparator
	mu         sync.RWMutex
}

func New(comp utils.Comparator) *SkipList {
	var skiplist SkipList
	skiplist.head = newNode(nil, kMaxHeight)
	skiplist.maxHeight = 1
	skiplist.comparator = comp
	return &skiplist
}

func (list *SkipList) Insert(key interface{}) {
	list.mu.Lock()
	defer list.mu.Unlock()
	prevNode, prev := list.findGreaterOrEqual(key)
	if prevNode != nil && list.comparator(prevNode.key, key) == 0 {
		return
	}
	height := list.randomHeight()
	if height > list.maxHeight {
		for i := list.maxHeight; i < height; i++ {
			prev[i] = list.head
		}
		list.maxHeight = height
	}
	x := newNode(key, height)
	for i := 0; i < height; i++ {
		x.setNext(i, prev[i].getNext(i))
		prev[i].setNext(i, x)
	}
}
func (list *SkipList) Contains(key interface{}) bool {
	list.mu.RLock()
	defer list.mu.RUnlock()
	x, _ := list.findGreaterOrEqual(key)
	if x != nil && list.comparator(x.key, key) == 0 {
		return true
	}
	return false
}
func (list *SkipList) findGreaterOrEqual(key interface{}) (*Node, [kMaxHeight]*Node) {
	var prev [kMaxHeight]*Node
	x := list.head
	level := list.maxHeight - 1
	for {
		next := x.getNext(level)
		if list.keyIsAfterNode(key, next) {
			x = next
		} else {
			prev[level] = x
			if level == 0 {
				return next, prev
			} else {
				// Switch to next list
				level--
			}
		}
	}
	return nil, prev
}
func (list *SkipList) findLessThan(key interface{}) *Node {
	x := list.head
	level := list.maxHeight - 1
	for {
		next := x.getNext(level)
		if next == nil || list.comparator(next.key, key) >= 0 {
			if level == 0 {
				return x
			} else {
				level--
			}
		} else {
			x = next
		}
	}
}
func (list *SkipList) NewIterator() *Iterator {
	return &Iterator{
		list: list,
	}
}

func (list *SkipList) randomHeight() int {
	var height = 1
	for height < kMaxHeight && (rand.Intn(kBranching) == 0) {
		height++
	}
	return height
}

func (list *SkipList) keyIsAfterNode(key interface{}, n *Node) bool {
	return (n != nil) && (list.comparator(n.key, key) < 0)
}

func (list *SkipList) findLast() *Node {
	x := list.head
	level := list.maxHeight - 1
	for true {
		next := x.getNext(level)
		if next == nil {
			if level == 0 {
				return x
			} else {
				level--
			}
		} else {
			x = next
		}
	}
	return nil
}
