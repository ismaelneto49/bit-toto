package core

import (
	"fmt"
	"io"
	"log"

	"github.com/ismaelneto49/bit-toto/src/services/peerconnection"
)

func InitClient(peerConn *peerconnection.PeerConnectionImpl, depth uint32) {
	for {
		var cmd, filename string
		fmt.Println("[CLIENT] Enter command. Accepted commands:")
		fmt.Println("- request <filename>")
		fmt.Println("- quit")
		_, err := fmt.Scan(&cmd)
		if err != nil {
			if err == io.EOF {
				log.Println("[CLIENT] EOF received, exiting client REPL.")
				return
			}
			log.Println("[CLIENT] Error reading command:", err)
			continue
		}
		if cmd == "quit" {
			log.Println("[CLIENT] Exiting client REPL.")
			break
		}
		if cmd == "request" {
			_, err := fmt.Scan(&filename)
			if err != nil {
				log.Println("[CLIENT] Error reading filename:", err)
				continue
			}
			log.Printf("[CLIENT] Search Depth: %d\n", depth)
			err = peerConn.GetFile(filename, depth)
			if err != nil {
				log.Println("[CLIENT] Error searching file:", err)
			}
		} else {
			log.Println("[CLIENT] Unknown command.")
		}
	}
}
