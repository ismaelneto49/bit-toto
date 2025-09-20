package peerconnection

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ismaelneto49/bit-toto/src/helpers"
	"github.com/ismaelneto49/bit-toto/src/tcpclient"
)

// "Attributes"
type PeerConnectionImpl struct {
	ip        net.TCPAddr
	knownIps  []*net.TCPAddr
	idCounter uint32
	visitedBy map[string]bool
}

// "Constructor"
func NewPeerConnection(knownIps []*net.TCPAddr) (*PeerConnectionImpl, error) {
	return &PeerConnectionImpl{
		ip:        net.TCPAddr{},
		knownIps:  knownIps,
		idCounter: 0,
		visitedBy: make(map[string]bool),
	}, nil
}

// ========================CLIENT===========================

// func (conn *PeerConnectionImpl) Join(targetIp *net.TCPAddr) error {
// 	// connect with targetIp

// 	knownIps, err := tcpclient.Join(targetIp)
// 	helpers.Treat(err)
// 	conn.knownIps = append(conn.knownIps, knownIps...)
// 	return nil
// }

func (conn *PeerConnectionImpl) GetFile(fileName string, depth uint32) error {
	searchId := fmt.Sprintf("%s%d", conn.ip.IP.To16().String(), conn.idCounter+1)
	var ipWithFile *net.TCPAddr = nil

	conn.visitedBy[searchId] = true
	for _, ip := range conn.knownIps {
		foundIp, err := conn.Search(ip, searchId, fileName, depth)
		if err != nil {
			log.Println("[CLIENT] Error searching:", err)
			continue
		}
		if foundIp != nil {
			ipWithFile = foundIp
			log.Println("[CLIENT] Found file at the ip ", ipWithFile)
			break
		}
	}

	if ipWithFile == nil {
		return fmt.Errorf("file with name %s not found in network", fileName)
	}

	go download(ipWithFile, fileName)

	return nil
}

func download(targetIp *net.TCPAddr, fileName string) error {
	// connect with targetIp
	// save file to disk at files/fileName.txt
	// LOG: file saved at files/fileName.txt
	return nil
}

func (conn *PeerConnectionImpl) Search(targetIp *net.TCPAddr, searchId string, fileName string, depth uint32) (*net.TCPAddr, error) {
	// connect with targetIp
	ip, err := tcpclient.Search(targetIp, searchId, fileName, depth)
	return ip, err
}

// ========================SERVER===========================

func (conn *PeerConnectionImpl) Connect(originIp *net.TCPAddr) []*net.TCPAddr {
	conn.knownIps = append(conn.knownIps, originIp)
	return conn.knownIps
}

func (conn *PeerConnectionImpl) ForwardSearch(searchId string, fileName string, depth uint32) (*net.TCPAddr, error) {
	if conn.visitedBy[searchId] {
		return nil, nil
	}
	conn.visitedBy[searchId] = true

	filePath := fmt.Sprintf("files/%s/%s", conn.ip.IP.String(), fileName)
	_, err := os.Stat(filePath)
	fileExists := err == nil
	if fileExists {
		return &conn.ip, nil
	}

	if depth == 1 {
		return nil, nil
	}

	decDepth := depth - 1
	for _, targetIp := range conn.knownIps {
		ipWithFile, err := conn.Search(targetIp, searchId, fileName, decDepth)
		helpers.Treat(err)
		if ipWithFile != nil {
			log.Println("[CLIENT] Found file at the ip ", ipWithFile)
			return ipWithFile, nil
		}
	}
	return nil, nil
}

func (conn *PeerConnectionImpl) ProvideFile(fileName string) ([]byte, error) {
	// read file from disk
	file := []byte("file")
	return file, nil
}
