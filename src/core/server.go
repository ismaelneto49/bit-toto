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

type TCPAddrData struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

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
		log.Println("[SERVER] New connection from", req.ClientPort)

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
		depthFloat, ok3 := params["depth"].(float64)
		depth := uint32(depthFloat)

		if !ok || !ok2 || !ok3 {
			return tcpclient.Response{Status: "error", Error: "invalid param types"}
		}
		ip, err := peerConn.ForwardSearch(searchId, fileName, depth)
		if err != nil {
			return tcpclient.Response{Status: "error", Error: err.Error()}
		}
		return tcpclient.Response{
			Status: "ok",
			Data: TCPAddrData{
				IP:   ip.IP.String(),
				Port: ip.Port,
			},
		}
	case "download":
		params, ok := req.Params.(map[string]interface{})
		if !ok {
			return tcpclient.Response{Status: "error", Error: "invalid params"}
		}
		fileName, ok := params["fileName"].(string)
		if !ok {
			return tcpclient.Response{Status: "error", Error: "invalid param types"}
		}

		file, err := peerConn.ProvideFile(fileName)
		if err != nil {
			return tcpclient.Response{Status: "error", Error: err.Error()}
		}
		return tcpclient.Response{
			Status: "ok",
			Data:   file,
		}
	default:
		return tcpclient.Response{Status: "error", Error: "unknown action"}
	}
}
