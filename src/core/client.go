package core

import (
	"fmt"
	"io"

	"github.com/ismaelneto49/bit-toto/src/services/peerconnection"
)

func InitClient(peerConn *peerconnection.PeerConnectionImpl) {
	for {
		var cmd, filename string
		fmt.Println("[CLIENT] Enter command. Accepted commands:")
		fmt.Println("- request <filename>")
		fmt.Println("- quit")
		_, err := fmt.Scan(&cmd)
		if err != nil {
			if err == io.EOF {
				fmt.Println("[CLIENT] EOF received, exiting client REPL.")
				return
			}
			fmt.Println("[CLIENT] Error reading command:", err)
			continue
		}
		if cmd == "quit" {
			fmt.Println("CLIENT] Exiting client REPL.")
			break
		}
		if cmd == "request" {
			_, err := fmt.Scan(&filename)
			if err != nil {
				fmt.Println("CLIENT] Error reading filename:", err)
				continue
			}
			DEPTH := uint32(3)
			fmt.Printf("[CLIENT] Search Depth: %d\n", DEPTH)
			err = peerConn.GetFile(filename, DEPTH)
			if err != nil {
				fmt.Println("CLIENT] Error requesting file:", err)
			} else {
				fmt.Printf("CLIENT] Requested file: %s\n", filename)
			}
		} else {
			fmt.Println("CLIENT] Unknown command.")
		}
	}
}
