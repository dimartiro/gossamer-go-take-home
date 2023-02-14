package network

import (
	"container/list"
	"errors"
)

// MessageTracker tracks a configurable fixed amount of messages.
// Messages are stored first-in-first-out.  Duplicate messages should not be stored in the queue.
type MessageTracker interface {
	// Add will add a message to the tracker, deleting the oldest message if necessary
	Add(message *Message) (err error)
	// Delete will delete message from tracker
	Delete(id string) (err error)
	// Get returns a message for a given ID.  Message is retained in tracker
	Message(id string) (message *Message, err error)
	// Messages returns messages in FIFO order
	Messages() (messages []*Message)
}

// ErrMessageNotFound is an error returned by MessageTracker when a message with specified id is not found
var ErrMessageNotFound = errors.New("message not found")

// messageTracker implementation
type messageTracker struct {
	messages *list.List               //Used as dequeue
	mapping  map[string]*list.Element //To get elements in O(1)
	length   int                      //Max elements in tracker
}

// NewMessageTracker create a new MessageTracker instance with the given max length
func NewMessageTracker(length int) MessageTracker {
	return &messageTracker{
		messages: list.New(),
		mapping:  make(map[string]*list.Element),
		length:   length,
	}
}

// Add will add a message to the tracker, deleting the oldest message if necessary
func (mt *messageTracker) Add(message *Message) error {
	//Ignore duplicates
	_, found := mt.mapping[message.ID]
	if found {
		return nil
	}

	if mt.messages.Len() >= mt.length {
		//Remove oldest message if queue is full
		oldest := mt.messages.Front()
		delete(mt.mapping, oldest.Value.(*Message).ID)
		mt.messages.Remove(oldest)
	}

	elem := mt.messages.PushBack(message)
	mt.mapping[message.ID] = elem
	return nil
}

// Delete will delete message from tracker
func (mt *messageTracker) Delete(id string) error {
	elem, ok := mt.mapping[id]
	if !ok {
		return ErrMessageNotFound
	}

	delete(mt.mapping, id)
	mt.messages.Remove(elem)
	return nil
}

// Get returns a message for a given ID.  Message is retained in tracker
func (mt *messageTracker) Message(id string) (*Message, error) {
	elem, ok := mt.mapping[id]
	if !ok {
		return nil, ErrMessageNotFound
	}

	return elem.Value.(*Message), nil
}

// Messages returns messages in FIFO order
func (mt *messageTracker) Messages() []*Message {
	messages := make([]*Message, len(mt.mapping))
	i := 0
	for elem := mt.messages.Front(); elem != nil; elem = elem.Next() {
		messages[i] = elem.Value.(*Message)
		i++
	}

	return messages
}
