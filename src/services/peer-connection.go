package network

import "net"

type PeerConnection interface {
	// Client
	Join(targetIp net.TCPAddr) error
	GetFile(fileName string, depth uint32) error
	Search(targetIp *net.TCPAddr, searchId string, fileName string, depth uint32) (*net.TCPAddr, error)

	// Server
	Connect(originIp *net.TCPAddr) []*net.TCPAddr
	ForwardSearch(searchId string, fileName string, depth uint32) (*net.TCPAddr, error)
	ProvideFile(fileName string) (file []byte, err error)
}
