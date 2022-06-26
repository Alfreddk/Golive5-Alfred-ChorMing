// This are the linked list functions customized to bookings
// Note the linked list position uses One-base addressing
package main

import (
	"errors"
	"sync"
)

var mutex sync.Mutex

// Error messages for linked list operation
var (
	HeadNullPointer = errors.New("Null Head Pointer")
	NullPointer     = errors.New("Null Pointer")
	IndexOutOfRange = errors.New("Index position out of range")
)

// node is the struct for the booking
type node struct {
	item Item  // to store the booking item
	next *node // pointer to point to next node
}

// linkedList is the struct for the linked list
type linkedList struct {
	head *node
	size int
}

// addNode adds a node into the linked list
func (p *linkedList) addNode(item Item) error {
	mutex.Lock()
	defer mutex.Unlock()

	// create a new node
	newNode := &node{
		item: item,
		next: nil,
	}

	if p.head == nil { // check if it is empty pointer
		p.head = newNode
		//		fmt.Println(newNode.item)
	} else {
		// start from the head
		currentNode := p.head
		// traversal until the last node
		for currentNode.next != nil {
			currentNode = currentNode.next
		}
		currentNode.next = newNode
		//		fmt.Println(newNode.item)
	}
	p.size++
	return nil
}

// printAllNodes prints all the items in the linked list
func (p *linkedList) printAllNodes() error {
	mutex.Lock()
	defer mutex.Unlock()

	if p.head == nil { // check if it is empty pointer
		//		fmt.Println("There are no nodes in the list.")
		return HeadNullPointer
	} else {
		currentNode := p.head
		//		fmt.Printf("%v\n", currentNode.item)
		//		currentNode.item.printItem()
		// traverse to next node
		for currentNode.next != nil {
			currentNode = currentNode.next
			//	fmt.Printf("%v\n", currentNode.item)
			//			currentNode.item.printItem()
		}
	}
	return nil
}

// removeNodes removed a node at position pos with the first node as position 1
func (p *linkedList) removeNode(pos int) (*Item, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if p.head == nil { // check if it is empty pointer
		//		fmt.Println("There are no nodes in the list.")
		return nil, HeadNullPointer
	}

	// use pos 1 as 1, so if you have n item, there are only n positions
	var item Item
	if pos > 0 && pos <= p.size { // check valid position index
		if pos == 1 { // special case where pos=1, where item 1 is to be removed
			item = p.head.item   // copy item before removal
			p.head = p.head.next // adjust the pointer in in p.head set to next item
		} else {
			currentNode := p.head // normal case, where the pointers are in the link list
			var prevNode *node
			for i := 1; i < pos; i++ { // traverse the link (pos-1) times to find the right node
				prevNode = currentNode         // keep a copy of previous node
				currentNode = currentNode.next // move to next node until the right node is found
			}
			item = currentNode.item          // capture item first before removal
			prevNode.next = currentNode.next // point previous node next to the node after the removed node
		}
		p.size--          // in all valid cases, reduce size by 1
		return &item, nil // return item and error
	} else {
		// invalid index position
		return nil, IndexOutOfRange
	}
}

// addNodeAt addes a node at position pos with the booking item
func (p *linkedList) addNodeAt(item Item, pos int) error {
	mutex.Lock()
	defer mutex.Unlock()

	// create a new node
	newNode := &node{
		item: item, // initialise the new node with booking item
		next: nil,
	}

	// if p.head is null, then it is am empty linked list
	// add node to head
	if p.head == nil {
		p.head = newNode
		return nil
	} else {
		// head has a node
		currentNode := p.head // keep this first node, since this points to the first item
		var prevNode *node
		// if pos is valid
		if pos > 0 && pos <= p.size {
			{
				if pos == 1 {
					// exception for pos=1, p.head to point to new node
					newNode.next = currentNode
					p.head = newNode
				} else {
					// traversal until the last node
					for i := 1; i < pos; i++ { // traverse pos-1 time
						prevNode = currentNode         // keep a copy of previous node before it is traversed
						currentNode = currentNode.next // move to next node until the right node is found
					}
					newNode.next = currentNode // always change the next node pointer first
					prevNode.next = newNode
				}
			}
			p.size++
			return nil
		} else {
			return IndexOutOfRange
		}
	}
}

// getItem return the booking item at pos int
func (p *linkedList) getItem(pos int) (*Item, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if p.head == nil { // check if it is empty pointer
		//		fmt.Println("There are no nodes in the list.")
		return nil, HeadNullPointer
	} else {
		if pos > 0 && pos <= p.size {
			currentNode := p.head
			var prevNode *node
			for i := 1; i <= pos; i++ { // traverse pos times
				prevNode = currentNode
				currentNode = currentNode.next
			}
			if prevNode != nil {
				return &prevNode.item, nil // return content
			} else {
				//				fmt.Println("Null pointer")
				return nil, NullPointer
			}
		} else {
			//			fmt.Println("invalid Index")
			return nil, IndexOutOfRange
		}
	}
}

// getAllItems collects all item pointers in a slice
func (p *linkedList) getAllItems(itemList *[]Item) error {
	mutex.Lock()
	defer mutex.Unlock()

	if p.head == nil { // check if it is empty pointer
		//		fmt.Println("There are no nodes in the list.")
		(*itemList) = nil // set to nil pointer is there is no item
		return HeadNullPointer
	} else {
		currentNode := p.head
		//	store address of the item
		*itemList = append(*itemList, (currentNode.item))
		// traverse to next node
		for currentNode.next != nil {
			currentNode = currentNode.next
			*itemList = append(*itemList, (currentNode.item))
		}
	}
	return nil
}

// getAllItemPtrs2 collects the addresses of item pointers in a slice
func (p *linkedList) getAllItemPtrs2(itemList *[]*Item) error {
	mutex.Lock()
	defer mutex.Unlock()

	if p.head == nil { // check if it is empty pointer
		//		fmt.Println("There are no nodes in the list.")
		(*itemList) = nil // set to nil pointer is there is no item
		return HeadNullPointer
	} else {
		currentNode := p.head
		//	store address of the item
		*itemList = append(*itemList, &(currentNode.item))
		// traverse to next node
		for currentNode.next != nil {
			currentNode = currentNode.next
			*itemList = append(*itemList, &(currentNode.item))
		}
	}
	return nil
}
