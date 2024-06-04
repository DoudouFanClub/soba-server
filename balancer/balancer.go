package balancer

import (
	"llm_server/socket"
)

type Balancer struct {
	availables Queue[socket.Endpoint]
}

func CreateBalancer() Balancer {
	q := NewQueue[socket.Endpoint]()
	return Balancer{availables: q}
}

func (b *Balancer) Available() bool {
	return b.availables.Empty()
}

func (b *Balancer) Send(msg []byte, buf *[]byte) {
	sender := b.availables.Pop()
	sender.SendMessage(msg, buf)
	b.availables.Add(sender)
}
