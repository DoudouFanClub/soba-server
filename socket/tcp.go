package socket

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
)

// Holds the endpoint info incase there is multiple llm instances
type Endpoint struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

// passes in a reference buffer to read from
func (e *Endpoint) SendMessage(message []byte, w *http.ResponseWriter) (bool, string) {

	result := ""
	conn, err1 := net.Dial("tcp", e.GetAddress())
	if err1 != nil {
		fmt.Println(err1)
		return false, ""
	}

	//defer conn.Close()

	_, err2 := conn.Write(message)
	if err2 != nil {
		fmt.Println(err2)
		return false, ""
	}

	// NEW
	flusher, ok := (*w).(http.Flusher)
	if !ok {
		fmt.Println("Unable to construct a flusher")
		return false, ""
	}
	reader := bufio.NewReader(conn)
	buf := make([]byte, 1024)

	for {
		n, err := reader.Read(buf)
		//fmt.Println("Value of n from conn.readbuffer: ", n)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from connection:", err)
			break
		}

		chunk := make([]byte, n)
		copy(chunk, buf[:n])
		result += string(chunk)

		_, err2 := fmt.Fprintf(*w, "%s", chunk)
		//fmt.Println("Buffer to be sent:", buf[:n])
		flusher.Flush()

		if err2 != nil {
			fmt.Println("Error sending message to frontend: ", err)
		}
	}

	fmt.Println(result)

	return true, result
}

func (e *Endpoint) GetAddress() string {
	return e.Ip + ":" + e.Port
}
