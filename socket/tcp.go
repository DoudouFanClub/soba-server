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
	Ip   string
	Port string
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
	fmt.Println("endpt: ", message)
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
		//fmt.Println(string(buf[:n]))
		//http.ResponseWriter.Write(*w, chunk)

		//fmt.Printf("Iteration %s\n", string(chunk))
		// NEW
		bytesWritten, err := fmt.Fprintf(*w, "%s", chunk)
		//fmt.Println("Buffer to be sent:", buf[:n])
		flusher.Flush();
		
		if err == nil {
			fmt.Println("Number of bytes written to frontend:", bytesWritten)
		} else {
			fmt.Println("Error sending message to frontend: ", err)
		}
	}
	

	fmt.Println(result)
    
	return true, result
}

// func (e *Endpoint) SendMessage(message []byte, w *http.ResponseWriter) (bool, string) {
//     result := ""
//     conn, err1 := net.Dial("tcp", e.GetAddress())
//     if err1 != nil {
//         fmt.Println(err1)
//         return false, ""
//     }

//     _, err2 := conn.Write(message)
//     fmt.Println("endpt: ", message)
//     if err2 != nil {
//         fmt.Println(err2)
//         return false, ""
//     }

//     flusher, ok := (*w).(http.Flusher)
//     if !ok {
//         fmt.Println("Unable to construct a flusher")
//         return false, ""
//     }

//     reader := bufio.NewReader(conn)
//     resultChan := make(chan string)
//     errorChan := make(chan error)
//     doneChan := make(chan struct{})
//     var wg sync.WaitGroup

//     wg.Add(1)
//     go func() {
//         defer wg.Done()
//         var localResult string
//         for {
//             buf := make([]byte, 10)
//             n, err := reader.Read(buf)
//             if err != nil {
//                 if err == io.EOF {
//                     break
//                 }
//                 errorChan <- err
//                 return
//             }
//             chunk := buf[:n]
//             localResult += string(chunk)
//             resultChan <- string(chunk)
//         }
//         close(resultChan)
//         result = localResult
//     }()

//     go func() {
//         for {
//             select {
//             case chunk, ok := <-resultChan:
//                 if !ok {
//                     close(doneChan)
//                     return
//                 }
//                 if _, err := (*w).Write([]byte(chunk)); err != nil {
//                     fmt.Println("Error writing to frontend:", err)
//                     errorChan <- err
//                     return
//                 }
//                 flusher.Flush()
//             case err := <-errorChan:
//                 fmt.Println("Error reading from connection:", err)
//                 return
//             }
//         }
//     }()

//     go func() {
//         wg.Wait()
//         conn.Close() // Correctly close the connection
//         close(doneChan)
//     }()

//     <-doneChan

//     return true, result
// }

func (e *Endpoint) GetAddress() string {
	return e.Ip + ":" + e.Port
}
