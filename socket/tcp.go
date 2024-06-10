package socket

import (
	"encoding/json"
	"fmt"
	"io"
	"llm_server/database"
	"net"
	"net/http"
)

// Holds the endpoint info incase there is multiple llm instances
type Endpoint struct {
	Ip   string
	Port string
}

/*
	read up on go channels for connections (fs read set?)
*/
// passes in a reference buffer to read from
func (e *Endpoint) SendMessage(message []byte, w *http.ResponseWriter) bool {

	conn, err1 := net.Dial("tcp", e.GetAddress())
	if err1 != nil {
		fmt.Println(err1)
		return false
	}

	defer conn.Close()

	_, err2 := conn.Write(message)
	if err2 != nil {
		fmt.Println(err2)
		return false
	}

	decoder := json.NewDecoder(conn)
	for {
		var msg database.Message
		err := decoder.Decode(&msg)
		http.ResponseWriter.Write(*w, []byte(msg.Content))
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return false
		}
	}
	return true
}

func (e *Endpoint) GetAddress() string {
	return e.Ip + ":" + e.Port
}
