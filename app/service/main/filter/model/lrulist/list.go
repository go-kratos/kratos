package lrulist

// 链表
type List struct {
	head   *Node // 头部
	tail   *Node // 尾部
	length int   // 长度
}

// new 链表
func New() *List {
	l := new(List)
	l.length = 0
	return l
}

// 节点
type Node struct {
	Key   string
	Value interface{}
	prev  *Node
	next  *Node
}

// new 节点
func NewNode(key string, value interface{}) *Node {
	return &Node{Key: key, Value: value}
}

// 长度
func (l *List) Len() int {
	return l.length
}

// 获取头部
func (l *List) Header() *Node {
	return l.head
}

// 获取尾部
func (l *List) Tailer() *Node {
	return l.tail
}

// 是否为空
func (l *List) IsEmpty() bool {
	return l.length == 0
}

// 头插
func (l *List) Prepend(key string, value interface{}) {
	node := NewNode(key, value)
	if l.Len() == 0 {
		l.head = node
		l.tail = l.head
	} else {
		formerHead := l.head
		formerHead.prev = node
		node.next = formerHead
		l.head = node
	}
	l.length++
}

// 尾插
func (l *List) Append(key string, value interface{}) {
	node := NewNode(key, value)
	if l.Len() == 0 {
		l.head = node
		l.tail = l.head
	} else {
		formerTail := l.tail
		formerTail.next = node
		node.prev = formerTail
		l.tail = node
	}
	l.length++
}

// 移动到头部
func (l *List) MoveToHead(node *Node) {
	if l.Len() == 0 {
		l.head = node
		l.tail = l.head
	} else {
		formerHead := l.head
		formerHead.prev = node
		node.next = formerHead
		l.head = node
	}
	l.length++
}

// 查找
func (l *List) Find(key string) *Node {
	for node := l.head; node != nil; node = node.next {
		if node.Key == key {
			return node
		}
	}
	return nil
}

// 删除
func (l *List) Remove(node *Node) *Node {
	if node == l.head {
		l.head = node.next
		node.next.prev = nil
	} else if node == l.tail {
		l.tail = node.prev
		node.prev.next = nil
	} else {
		node.prev.next = node.next
		node.next.prev = node.prev
	}
	node.prev = nil
	node.next = nil
	l.length--
	return node
}

// 遍历
func (l *List) Each(f func(node Node)) {
	for node := l.head; node != nil; node = node.next {
		f(*node)
	}
}

// 清空
func (l *List) Empty() {
	l.head = nil
	l.tail = nil
	l.length = 0
}
