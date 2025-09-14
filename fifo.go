package main

func NewFifoChan(bufferSize int) Queue {
	return &FifoChan{
		c: make(chan Message, bufferSize),
	}
}

type FifoChan struct {
	c chan Message
}

func (fcr *FifoChan) Next() (tenant string, obj int) {
	v := <-fcr.c
	return v.tenant, v.message
}

func (fcr *FifoChan) Put(tenant string, obj int) error {
	m := Message{tenant, obj}
	select {
	case fcr.c <- m:
	default:
		return ChanWasFull
	}
	return nil
}

func (fcr *FifoChan) IsEmpty() bool {
	return len(fcr.c) == 0
}
