package network

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const (
	connectionTimeout = 1500 * time.Millisecond
	nodeTimeout       = 1500 * time.Millisecond
)

type TCPclient struct {
	address        string
	port           int
	ReceiveChannel chan string
}

func NewTcpClient(ip string, port int) *TCPclient {
	client := &TCPclient{
		address:        ip,
		port:           port,
		ReceiveChannel: make(chan string, 100),
	}
	return client
}

func (t *TCPclient) UpdateTarget(ip string, port int) {
	t.address = ip
	t.port = port
}

func (t *TCPclient) Send(message string) (res string, err error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", t.address, t.port), nodeTimeout)
	if err != nil {
		fmt.Println("Error connecting to node: \n", err)
		return "", err
	} else {
		fmt.Println("Connection stablished")
		// Start goroutines
		go t.read(conn)
		message = fmt.Sprintf("%s\n", message)
		fmt.Fprintf(conn, message)
		response := <-t.ReceiveChannel
		conn.Close()
		return response, nil
	}
}

func (t *TCPclient) read(conn net.Conn) {
	reader := bufio.NewScanner(conn)
	for {
		if ok := reader.Scan(); !ok {
			if reader.Err() != nil {
				fmt.Println("Error in connection: ", reader.Err())
			}
			continue
		}

		response := reader.Text()
		if response == "" {
			continue
		}

		t.ReceiveChannel <- response
		break
	}
	fmt.Println("GoRoutine for Reading is finished")
}
