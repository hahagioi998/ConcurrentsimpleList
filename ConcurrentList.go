package main

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// 使用 Go 语言完成一个 并发安全的有序链表（数据严格有序并且没有重复元素）

// key point
//     一写多读的场景读只需要 atomic
//     一写多读的场景写需要   atomic+lock

const (
	isDeleted = 1
	unDeleted = 0
)



type IntList struct {
	head   *intNode
	length int64
}

type intNode struct {
	value     int
	next      *intNode
	mu        sync.Mutex
	isDeleted uint32
}

func newIntNode(value int) *intNode {
	return &intNode{value: value}
}


func NewInt() *IntList {
	return &IntList{head: newIntNode(0)}
}

func (l *IntList) Insert(value int) bool {
	var a, b *intNode
	for {
		a = l.head
		b = a.getNextNode()
		for b != nil && b.getValue() < value {
			a = b
			b = b.getNextNode()
		}
		if b != nil && b.value == value {
			return false
		}

		a.mu.Lock()
		// double check
		if a.getNextNode() == b {
			break
		} else {
			a.mu.Unlock()
			continue
		}
	}

	x := newIntNode(value)
	x.next = b
	a.setNextNode(x)

	atomic.AddInt64(&l.length, 1)
	a.mu.Unlock()
	return true
}

func (l *IntList) Delete(value int) bool {
	var a, b *intNode
	for {
		a = l.head
		b = a.getNextNode()
		for b != nil && b.getValue() < value {
			a = b
			b = b.getNextNode()
		}

		if b == nil || b.value != value {
			return false
		}

		b.mu.Lock()
		if atomic.LoadUint32(&b.isDeleted) == isDeleted {
			b.mu.Unlock()
			// 已经被删除
			return false
		}

		// 检查a有没有变化
		a.mu.Lock()
		if a.getNextNode() != b || atomic.LoadUint32(&a.isDeleted) == isDeleted {
			a.mu.Unlock()
			b.mu.Unlock()
			continue
		} else {
			// b为所寻找的值
			break
		}
	}
	atomic.StoreUint32(&b.isDeleted, isDeleted)
	a.setNextNode(b.getNextNode())
	atomic.AddInt64(&l.length, -1)
	a.mu.Unlock()
	b.mu.Unlock()
	return true
}

func (l *IntList) Contains(value int) bool {
	if l.Len() == 0 {
		return false
	}
	x := l.head.getNextNode()
	for x != nil && x.getValue() < value {
		x = x.getNextNode()
	}
	if x == nil || atomic.LoadUint32(&x.isDeleted) == isDeleted || x.value != value {
		return false
	}
	return true
}

func (l *IntList) Range(f func(value int) bool) {
	x := l.head.getNextNode()
	for x != nil {
		if !f(x.getValue()) {
			break
		}
		x = x.getNextNode()
	}
}

func (l *IntList) Len() int {
	return int(atomic.LoadInt64(&l.length))
}

// util
func (n *intNode) getValue() int {
	return int(atomic.LoadInt64((*int64)(unsafe.Pointer(&n.value))))
}

func (n *intNode) setValue(value int) {
	atomic.StoreInt64((*int64)(unsafe.Pointer(&n.value)), int64(value))
}

func (n *intNode) getNextNode() *intNode {
	return (*intNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&n.next))))
}

func (n *intNode) setNextNode(next *intNode) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&n.next)), unsafe.Pointer(next))
}