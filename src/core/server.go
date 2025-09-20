package core

import (
	"bufio"
	"encoding/json"
	"log"
	"net"

	"github.com/ismaelneto49/bit-toto/src/helpers"
	"github.com/ismaelneto49/bit-toto/src/services/peerconnection"
	"github.com/ismaelneto49/bit-toto/src/tcpclient"
)

var peerConn *peerconnection.PeerConnectionImpl

func InitServer(pc *peerconnection.PeerConnectionImpl, port string) {
	peerConn = pc
	ln, err := net.Listen("tcp", ":"+port)
	helpers.Treat(err)
	log.Println("[SERVER] Server listening on :" + port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("[SERVER] Error when accepting the connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	log.Println("[SERVER] New connection from", conn.RemoteAddr())
	defer conn.Close()
	reader := bufio.NewReader(conn)
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(conn)

	for {
		var req tcpclient.Request
		if err := decoder.Decode(&req); err != nil {
			log.Println("[SERVER] Error decoding:", err)
			return
		}

		resp := handleRequest(req)

		if err := encoder.Encode(resp); err != nil {
			log.Println("[SERVER] Error encoding:", err)
			return
		}
	}
}

func handleRequest(req tcpclient.Request) tcpclient.Response {
	switch req.Action {
	case "search":
		// Params come in as map[string]interface{}
		params, ok := req.Params.(map[string]interface{})
		if !ok {
			return tcpclient.Response{Status: "error", Error: "invalid params"}
		}
		searchId, ok := params["searchId"].(string)
		fileName, ok2 := params["fileName"].(string)
		depth, ok3 := params["depth"].(uint32)

		if !ok || !ok2 || !ok3 {
			return tcpclient.Response{Status: "error", Error: "invalid param types"}
		}
		ip, err := peerConn.ForwardSearch(searchId, fileName, depth)
		return tcpclient.Response{Status: "ok", Data: ip, Error: err.Error()}
	default:
		return tcpclient.Response{Status: "error", Error: "unknown action"}
	}
}
