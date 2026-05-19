package main

import (
	"fmt"
)

type Node struct {
	value int
	next  *Node
}

type LinkedList struct {
	head *Node
}

func New() *LinkedList {
	return &LinkedList{}
}

func (l *LinkedList) Add(value int) {
	newNode := &Node{value: value}
	if l.head == nil {
		l.head = newNode
	} else {
		current := l.head
		for current.next != nil {
			current = current.next
		}
		current.next = newNode
	}
}

func (l *LinkedList) Print() {
	current := l.head
	for current != nil {
		fmt.Printf("%d -> ", current.value)
		current = current.next
	}
}

func (l *LinkedList) Delete(value int) {
}

func (l *LinkedList) Remove(target int) {
	current := l.head

	for current.next != nil {

		if current.next.value == target {
			fmt.Printf("Value:  %d \n", current.value)
			current.next = current.next.next

		}
		current = current.next
	}
}

func main() {
	list := New()
	list.Add(1)
	list.Add(2)
	list.Add(3)
	list.Add(4)
	list.Add(5)
	list.Add(6)

	list.Print()
	list.Remove(4)
	fmt.Println("")
	list.Print()
	list.Remove(5)
	fmt.Println("")
	list.Print()
}
