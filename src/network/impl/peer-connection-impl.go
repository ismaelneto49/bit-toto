package impl

import (
	"fmt"
	"github.com/ismaelneto49/bit-toto/src/helpers"
	"net"
)

// "Attributes"
type PeerConnectionImpl struct {
	ip        net.TCPAddr
	knownIps  []*net.TCPAddr
	idCounter uint32
	visitedBy map[string]bool
}

// "Constructor"
func newPeerConnection() (*PeerConnectionImpl, error) {
	return &PeerConnectionImpl{
		ip:        net.TCPAddr{},
		knownIps:  []*net.TCPAddr{},
		idCounter: 0,
		visitedBy: make(map[string]bool),
	}, nil
}

// ========================CLIENT===========================

func (conn *PeerConnectionImpl) Join(targetIp net.TCPAddr) error {
	// connect with targetIp
	append(conn.knownIps, response.knownIps)
	return nil
}

func (conn *PeerConnectionImpl) GetFile(fileHash uint64, depth uint32) error {
	searchId := fmt.Sprintf("%s%d", conn.ip.IP.To16().String(), conn.idCounter+1)
	var ipWithFile *net.TCPAddr = nil

	conn.visitedBy[searchId] = true

	// for each neighbor in adjacency_list[node]:
	// 	if neighbor is not visited:
	// 			DFS(neighbor, visited)

	for _, ip := range conn.knownIps {
		foundIp, err := conn.Search(ip, searchId, fileHash, depth)
		helpers.Treat(err)
		if foundIp != nil {
			ipWithFile = foundIp
			// LOG: found file at ipWithFile
			break
		}
	}

	if ipWithFile == nil {
		// LOG: file not found in network
		return fmt.Errorf("file with hash %d not found in network", fileHash)
	}

	go download(*ipWithFile, fileHash)

	return nil
}

func download(targetIp net.TCPAddr, fileHash uint64) error {
	// connect with targetIp
	// save file to disk at files/fileHash.txt
	// LOG: file saved at files/fileHash.txt
	return nil
}

func (conn *PeerConnectionImpl) Search(targetIp *net.TCPAddr, searchId string, fileHash uint64, depth uint32) (*net.TCPAddr, error) {
	// connect with targetIp
	return nil, nil // Replace with actual logic as needed
}

// ========================SERVER===========================

func (conn *PeerConnectionImpl) Connect(originIp *net.TCPAddr) []*net.TCPAddr {
	conn.knownIps = append(conn.knownIps, originIp)
	return conn.knownIps
}

func (conn *PeerConnectionImpl) ForwardSearch(searchId string, fileHash uint64, depth uint32) (*net.TCPAddr, error) {
	if conn.visitedBy[searchId] {
		return nil, nil
	}
	conn.visitedBy[searchId] = true

	// check if I have the file
	// if I have the file, return my ip
	fileExists := false
	if fileExists {
		return &conn.ip, nil
	}

	if depth == 1 {
		return nil, nil
	}

	decDepth := depth - 1
	for _, targetIp := range conn.knownIps {
		foundIp, err := conn.Search(targetIp, searchId, fileHash, decDepth)
		helpers.Treat(err)
		if foundIp != nil {
			// LOG: found file at ipWithFile
			return foundIp, nil
		}
	}
	return nil, nil
}

func (conn *PeerConnectionImpl) ProvideFile(fileHash uint64) ([]byte, error) {
	// read file from disk
	file := []byte("file")
	return file, nil
}
