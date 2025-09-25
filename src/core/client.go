package core

import (
	"fmt"
	"time"

	"github.com/ismaelneto49/bit-toto/src/services/peerconnection"
)

func InitClient(peerConn *peerconnection.PeerConnectionImpl, depth uint32) {
	for i := 1; i < 100; i++ {
		time.Sleep(1 * time.Second)
		filename := "file" + fmt.Sprint(i) + ".txt"
		err := peerConn.GetFile(filename, depth)
		if err != nil {
			continue
		}
	}
}
