package balancer

import (
	"llm_server/socket"
	"net/http"
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

func (b *Balancer) Send(msg []byte, w *http.ResponseWriter) (bool, string) {
	sender := b.availables.Pop()
	success, result := sender.SendMessage(msg, w)
	if !success {
		return false, ""
	}
	b.availables.Add(sender)
	return true, result
}
