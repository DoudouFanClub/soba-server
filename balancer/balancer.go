package balancer

import (
	"fmt"
	"llm_server/socket"
	"net/http"
)

type Balancer struct {
	availables Queue[socket.Endpoint]
}

func CreateBalancer() *Balancer {
	q := NewQueue[socket.Endpoint]()
	return &Balancer{availables: q}
}

func (b *Balancer) Available() bool {
	return b.availables.Empty()
}

func (b *Balancer) Add(endpt socket.Endpoint) {
	b.availables.Add(endpt)
}

func (b *Balancer) Send(msg []byte, w *http.ResponseWriter) (bool, string) {
	sender := b.availables.Pop()
	success, result := sender.SendMessage(msg, w)
	if !success {
		fmt.Println("Unable to find the inferer endpoint")
		// add the endpoint back so that the server can continue to function when there's no endpoint
		b.availables.Add(sender)
		return false, ""
	}
	b.availables.Add(sender)
	return true, result
}
