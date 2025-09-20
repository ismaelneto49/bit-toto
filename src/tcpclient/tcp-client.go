package tcpclient

import (
	"bufio"
	"encoding/json"
	"net"
)

// Request and Response formats
type Request struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
}

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error,omitempty"`
}

func sendRequest(addr *net.TCPAddr, req Request) (Response, error) {
	// Create a TCP connection
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()

	// Encode request as JSON
	enc := json.NewEncoder(conn)
	if err := enc.Encode(req); err != nil {
		return Response{}, err
	}

	// Read response
	reader := bufio.NewReader(conn)
	dec := json.NewDecoder(reader)

	var resp Response
	if err := dec.Decode(&resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}
