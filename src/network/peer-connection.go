package network

import "net"

type PeerConnection interface {
	// Client
	Join(targetIp net.TCPAddr) ([]net.TCPAddr, error)
	GetFile(fileHash uint64, timeToLive uint32) error
	Search(ip *net.TCPAddr, searchId string, fileHash uint64, timeToLive uint32) (fileSize uint32, hasFileMap map[*net.TCPAddr]bool, err error)

	// Server
	Connect(originIp net.TCPAddr) []net.TCPAddr
	ForwardSearch(searchId string, fileHash uint64, timeToLive uint32) (hasFileMap map[*net.TCPAddr]bool, err error)
	ProvideFileChunk(fileHash uint64, startByte uint32, endByte uint32) (fileChunk []byte, err error)
	Clean(searchId string) error
}
