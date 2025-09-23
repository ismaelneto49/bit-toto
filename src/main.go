package main

import (
	"fmt"
	"net"

	"encoding/json"
	"os"

	"github.com/ismaelneto49/bit-toto/src/core"
	"github.com/ismaelneto49/bit-toto/src/helpers"
	"github.com/ismaelneto49/bit-toto/src/services/peerconnection"
)

type Config struct {
	Depth       uint32              `json:"depth"`
	Seed        int                 `json:"seed"`
	Connections map[string][]string `json:"connections"`
}

func main() {
	// Get "port" param from command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <port>")
		return
	}
	port := os.Args[1]

	// Open config file
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}

	defer file.Close()

	// Decode JSON config
	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		fmt.Println("Error decoding config file:", err)
		return
	}

	key := "127.0.0.1:" + port

	neighbors, ok := config.Connections[key]
	if !ok {
		fmt.Printf("No entry found for %s\n", key)
		return
	}
	fmt.Printf("Found entry for %s: %v\n", key, neighbors)

	var knownIps []*net.TCPAddr
	for _, n := range neighbors {
		addr, err := net.ResolveTCPAddr("tcp", n)
		helpers.Treat(err)

		knownIps = append(knownIps, addr)
	}

	peerConn, err := peerconnection.NewPeerConnection(knownIps, uint16(config.Seed))
	helpers.Treat(err)
	fmt.Println("[NODE] Node initialized with neighbors:", neighbors)

	// Pass config to InitClient
	go core.InitClient(peerConn, config.Depth)
	core.InitServer(peerConn, port)

}
