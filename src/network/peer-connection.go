package network

import "net"

type PeerConnection interface {
	// Client
	Join(targetIp net.TCPAddr) error
	GetFile(fileHash uint64, depth uint32) error
	Search(targetIp *net.TCPAddr, searchId string, fileHash uint64, depth uint32) (*net.TCPAddr, error)

	// Server
	Connect(originIp net.TCPAddr) []net.TCPAddr
	ForwardSearch(searchId string, fileHash uint64, depth uint32) (*net.TCPAddr, error)
	ProvideFile(fileHash uint64) (file []byte, err error)
}
