package main

import "errors"

type Queue interface {
	Next() (tenant string, message int)
	Put(tenant string, message int) error
	IsEmpty() bool
}

var ChanWasFull = errors.New("channel full")

type Message struct {
	tenant  string
	message int // use in represent the message info and time index.
}

func Generator() func(tenant string) Message {
	m := map[string]int{}

	return func(tenant string) Message {
		if _, ok := m[tenant]; !ok {
			m[tenant] = 0
		}

		m[tenant]++
		return Message{tenant, m[tenant]}
	}
}
