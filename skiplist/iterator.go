package skiplist

type Iterator struct {
	list *SkipList
	node *Node
}

func (it *Iterator) Valid() bool {
	return it.node != nil
}
func (it *Iterator) Key() interface{} {
	return it.node.key
}
func (it *Iterator) Next() {
	it.list.mu.Lock()
	defer it.list.mu.Unlock()
	it.node = it.node.getNext(0)
}
func (it *Iterator) Prev() {
	it.list.mu.Lock()
	defer it.list.mu.Unlock()
	it.node = it.list.findLessThan(it.node.key)
	if it.node == it.list.head {
		it.node = nil
	}
}
func (it *Iterator) Seek(key interface{}) {
	it.list.mu.Lock()
	defer it.list.mu.Unlock()
	it.node, _ = it.list.findGreaterOrEqual(key)
}
func (it *Iterator) SeekToFirst() {
	it.list.mu.Lock()
	defer it.list.mu.Unlock()
	it.node = it.list.head.getNext(0)
}
func (it *Iterator) SeekToLast() {
	it.list.mu.Lock()
	defer it.list.mu.Unlock()
	it.node = it.list.findLast()
	if it.node == it.list.head {
		it.node = nil
	}
}
