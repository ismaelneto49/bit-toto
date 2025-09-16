package impl

import (
	"fmt"
	"net"
)

// "Attributes"
type PeerConnectionImpl struct {
	ip        net.TCPAddr
	knownIps  []*net.TCPAddr
	idCounter uint32
	visited   map[string]bool
}

// "Constructor"
func newPeerConnection() (*PeerConnectionImpl, error) {
	return &PeerConnectionImpl{
		ip:        net.TCPAddr{},
		knownIps:  []*net.TCPAddr{},
		idCounter: 0,
		visited:   make(map[string]bool),
	}, nil
}

// ========================CLIENT===========================

func (conn *PeerConnectionImpl) Join(targetIp net.TCPAddr) error {
	// connect with targetIp
	append(conn.knownIps, response.knownIps)
	return nil
}

func (conn *PeerConnectionImpl) GetFile(fileHash uint64, timeToLive uint32) error {
	searchId := fmt.Sprintf("%s%d", conn.ip.IP.To16().String(), conn.idCounter+1)
	hasFileMap := make(map[*net.TCPAddr]bool)

	conn.visited[searchId] = true
	for _, ip := range conn.knownIps {
		go conn.Search(ip, searchId, fileHash, timeToLive)
		hasFileMap[ip] = true
	}

	// search returns it

	fileSize := 1000

	ipsWithFile := []*net.TCPAddr{}
	for ip, hasFile := range hasFileMap {
		if hasFile {
			ipsWithFile = append(ipsWithFile, ip)
		}
	}
	go download(fileHash, uint32(fileSize), ipsWithFile)

	for ip := range hasFileMap {
		go cleanup(ip, searchId)
	}
	return nil
}

func download(fileHash uint64, fileSize uint32, ipsWithFile []*net.TCPAddr) []byte {
	chunkSize := fileSize / uint32(len(ipsWithFile))
	chunks := calculateChunks(chunkSize, fileSize)
	for i, ip := range ipsWithFile {
		start, finish := chunks[i][0], chunks[i][1]
		go downloadChunk(ip, fileHash, start, finish)
	}
	// assembleFile -> file
	// downloadChunk writes in a chan (chunk, ip)
	// chunk number is ipsWithFile.indexOf(response.targetIP)
	return nil
}

func calculateChunks(chunkSize uint32, fileSize uint32) [][2]uint32 {
	chunks := [][2]uint32{}
	start := uint32(0)
	for start < fileSize {
		end := start + chunkSize
		if end > fileSize {
			end = fileSize
		}
		chunks = append(chunks, [2]uint32{start, end})
		start = end
	}
	return chunks
}

func downloadChunk(ip *net.TCPAddr, fileHash uint64, start uint32, finish uint32) []byte {
	// connect with ip
	// request chunk
	// return response.chunk
	return nil
}

func assembleFile(chunks map[int][]byte) []byte {
	// order chunks by index
	// join chunks
	return nil
}

func cleanup(ip *net.TCPAddr, searchId string) {
	// connect with ip
	// send cleanup
}

func (conn *PeerConnectionImpl) Search(ip *net.TCPAddr, searchId string, fileHash uint64, timeToLive uint32) (fileSize uint32, hasFileMap map[*net.TCPAddr]bool, err error) {
	// connect with ip
	return response.fileSize, response.hasFileMap, nil
}

// ========================SERVER===========================

func (conn *PeerConnectionImpl) Connect(originIp *net.TCPAddr) []*net.TCPAddr {
	conn.knownIps = append(conn.knownIps, originIp)
	return conn.knownIps
}

func (conn *PeerConnectionImpl) ForwardSearch(searchId string, fileHash uint64, timeToLive uint32) (map[*net.TCPAddr]bool, error) {
	if conn.visited[searchId] {
		return nil, nil
	}
	conn.visited[searchId] = true

	hasFileMap := make(map[*net.TCPAddr]bool)
	hasFileMap[&conn.ip] = true // check if I have the file

	if timeToLive == 1 {
		return hasFileMap, nil
	}

	decTTL := timeToLive - 1
	for _, ip := range conn.knownIps {
		go conn.Search(ip, searchId, fileHash, decTTL)
	}
	// read from channel and aggregate maps
	return hasFileMap, nil
}

func (conn *PeerConnectionImpl) ProvideFileChunk(fileHash uint64, startByte uint32, endByte uint32) ([]byte, error) {
	// read file from disk
	// split chunk
	// send chunk
	fileChunk := []byte("file chunk")
	return fileChunk, nil
}

func (conn *PeerConnectionImpl) Clean(searchId string) error {
	delete(conn.visited, searchId)
	return nil
}
