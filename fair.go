package main

import (
	"sync"
)

type FairChan struct {
	bufferSize int
	channels   map[string]chan Message
	tenants    []string
	index      int
	mu         sync.Mutex
}

func NewFairChan(bufferSize int) Queue {
	return &FairChan{
		bufferSize: bufferSize,
		channels:   make(map[string]chan Message),
		tenants:    make([]string, 0),
		index:      0,
	}
}

func (fc *FairChan) Put(tenant string, message int) error {
	fc.mu.Lock()
	if _, exists := fc.channels[tenant]; !exists {
		fc.channels[tenant] = make(chan Message, fc.bufferSize)
		fc.tenants = append(fc.tenants, tenant)
	}
	fc.mu.Unlock()

	msg := Message{tenant, message}
	select {
	case fc.channels[tenant] <- msg:
		return nil
	default:
		return ChanWasFull
	}
}

func (fc *FairChan) Next() (string, int) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if len(fc.tenants) == 0 {
		return "", 0
	}

	// Round-robin through tenants
	for i := 0; i < len(fc.tenants); i++ {
		tenant := fc.tenants[fc.index]

		// Try to get message from current tenant's channel
		select {
		case msg := <-fc.channels[tenant]:
			fc.index = (fc.index + 1) % len(fc.tenants)
			return msg.tenant, msg.message
		default:
			// No message from this tenant, move to next
			fc.index = (fc.index + 1) % len(fc.tenants)
		}
	}

	// No messages available from any tenant
	return "", 0
}

func (fc *FairChan) IsEmpty() bool {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if len(fc.tenants) == 0 {
		return true
	}

	// Check if all tenant channels are empty
	for _, ch := range fc.channels {
		if len(ch) > 0 {
			return false
		}
	}
	return true
}
