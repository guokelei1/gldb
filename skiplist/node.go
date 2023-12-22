package skiplist

type Node struct {
	key  interface{}
	next []*Node
}

func newNode(key interface{}, level int) *Node {
	return &Node{
		key:  key,
		next: make([]*Node, level),
	}
}
func (node *Node) getNext(level int) *Node {
	return node.next[level]
}

func (node *Node) setNext(level int, next *Node) {
	node.next[level] = next
}
