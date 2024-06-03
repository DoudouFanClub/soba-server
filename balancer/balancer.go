package balancer

import (
	"llm_server/socket"
)

type Balancer struct {
	Availables Queue[socket.Endpoint]
}

func CreateBalancer() Balancer {
	q := NewQueue[socket.Endpoint]()
	return Balancer{Availables: q}
}

func (b *Balancer) CanSend() bool {
	return b.Availables.Empty()
}

func (b *Balancer) Send(msg []byte, buf *[]byte) {
	sender := b.Availables.Pop()
	sender.SendMessage(msg, buf)
	b.Availables.Add(sender)
}
