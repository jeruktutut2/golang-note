package service

import (
	"context"
	"fmt"
)

type LinkedListService interface {
	Singly(ctx context.Context, requestId string)
}

type LinkedListServiceImplementation struct {
}

func NewLinkedListService() LinkedListService {
	return &LinkedListServiceImplementation{}
}

// https://www.youtube.com/watch?v=1S0_-VxPLJo and https://www.youtube.com/watch?v=8QoynPUY9_8
type singlyNodeLinkedList struct {
	data int
	next *singlyNodeLinkedList
}

type singlyLinkedList struct {
	head   *singlyNodeLinkedList
	length int
}

func (singlyLinkedList *singlyLinkedList) prepend(n *singlyNodeLinkedList) {
	second := singlyLinkedList.head
	singlyLinkedList.head = n
	singlyLinkedList.head.next = second
	singlyLinkedList.length++
}

func (singlyLinkedList singlyLinkedList) printSinglyLinkedList() {
	node := singlyLinkedList.head
	for singlyLinkedList.length != 0 {
		// fmt.Println(node.data)
		fmt.Printf("%d ", node.data)
		node = node.next
		singlyLinkedList.length--
	}
	fmt.Printf("\n")
}

func (singlyLinkedList *singlyLinkedList) deleteWithValue(value int) {
	previousToDelete := singlyLinkedList.head
	for previousToDelete.next.data != value {
		previousToDelete = previousToDelete.next
	}
	previousToDelete.next = previousToDelete.next.next
	singlyLinkedList.length--
}

func (service *LinkedListServiceImplementation) Singly(ctx context.Context, requestId string) {
	linkedList := singlyLinkedList{}
	node1 := &singlyNodeLinkedList{data: 48}
	node2 := &singlyNodeLinkedList{data: 18}
	node3 := &singlyNodeLinkedList{data: 16}
	node4 := &singlyNodeLinkedList{data: 11}
	node5 := &singlyNodeLinkedList{data: 7}
	node6 := &singlyNodeLinkedList{data: 2}
	linkedList.prepend(node1)
	linkedList.prepend(node2)
	linkedList.prepend(node3)
	linkedList.prepend(node4)
	linkedList.prepend(node5)
	linkedList.prepend(node6)
	linkedList.printSinglyLinkedList()
	linkedList.deleteWithValue(16)
	linkedList.printSinglyLinkedList()
}
