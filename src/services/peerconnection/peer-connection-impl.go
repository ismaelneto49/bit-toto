package peerconnection

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"math/rand"

	"github.com/ismaelneto49/bit-toto/src/helpers"
	"github.com/ismaelneto49/bit-toto/src/tcpclient"
)

// "Attributes"
type PeerConnectionImpl struct {
	ip        net.TCPAddr
	knownIps  []*net.TCPAddr
	idCounter uint32
	visitedBy map[string]bool
	seed      uint16
}

// "Constructor"
func NewPeerConnection(knownIps []*net.TCPAddr, seed uint16) (*PeerConnectionImpl, error) {
	port, err := strconv.Atoi(os.Args[1])
	helpers.Treat(err)
	return &PeerConnectionImpl{
		ip:        net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: port},
		knownIps:  knownIps,
		idCounter: 0,
		visitedBy: make(map[string]bool),
		seed:      uint16(42),
	}, nil
}

// ========================CLIENT===========================

func (conn *PeerConnectionImpl) GetFile(fileName string, depth uint32) error {
	port := os.Args[1]
	searchId := fmt.Sprintf("%s:%s-%d", conn.ip.IP.String(), port, conn.idCounter+1)
	conn.idCounter++
	var ipWithFile *net.TCPAddr = nil

	conn.visitedBy[searchId] = true
	for _, ip := range conn.knownIps {
		log.Println("[CLIENT] Searching for file at the ip: ", ip)
		foundIp, err := conn.Search(ip, searchId, fileName, depth)
		if err != nil {
			log.Println("[CLIENT] Error searching:", err)
			continue
		}
		if foundIp != nil {
			ipWithFile = foundIp
			log.Println("[CLIENT] Found file at the ip: ", ipWithFile)
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
	// file, err := tcpclient.Download(targetIp, fileName)
	_, err := tcpclient.Download(targetIp, fileName)
	if err != nil {
		return err
	}

	// save file to disk at files/fileName.txt
	port := os.Args[1]
	filesFolderName := targetIp.IP.String() + ":" + port
	filesFolderPath := fmt.Sprintf("files/%s/%s", filesFolderName, fileName)
	// if err := os.WriteFile(filesFolderPath, file, 0644); err != nil {
	// 	return err
	// }
	log.Println("[CLIENT] File saved at: ", filesFolderPath)
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
		return nil, errors.New("[SERVER] already visited by the searchId: " + searchId)
	}
	conn.visitedBy[searchId] = true

	port := os.Args[1]
	filesFolderName := conn.ip.IP.String() + ":" + port
	filePath := fmt.Sprintf("files/%s/%s", filesFolderName, fileName)
	_, err := os.Stat(filePath)
	fileExists := err == nil
	if fileExists {
		log.Println("[SERVER] File exists at ip: ", &conn.ip)
		fundura := 1000 - depth + 1
		log.Printf("[SAPATO] %s,%d\n", fileName, fundura)
		return &conn.ip, nil
	}
	log.Println("[SERVER] File does not exist at ip: ", &conn.ip)

	if depth == 1 {
		log.Println("[SERVER] Max search depth reached")
		return nil, errors.New("max search depth reached")

	}

	// Shuffle knownIps using the seed
	rng := rand.New(rand.NewSource(int64(conn.seed)))
	rng.Shuffle(len(conn.knownIps), func(i, j int) {
		conn.knownIps[i], conn.knownIps[j] = conn.knownIps[j], conn.knownIps[i]
	})

	decDepth := depth - 1
	for _, targetIp := range conn.knownIps {
		ipWithFile, err := conn.Search(targetIp, searchId, fileName, decDepth)
		if err != nil {
			if err == io.EOF {
				// Client disconnected — not fatal, just log it
				log.Println("[CLIENT] connection closed while searching", targetIp)
				continue
			}
			log.Println("[CLIENT] search error:", err)
			continue
		}
		if ipWithFile != nil {
			log.Println("[CLIENT] Found file at the ip ", ipWithFile)
			return ipWithFile, nil
		}
	}
	return nil, errors.New("[SERVER] file not found in neighbors")
}

func (conn *PeerConnectionImpl) ProvideFile(fileName string) ([]byte, error) {
	// read file from disk
	// TODO this code snipped repeats a lot. Refactor it later
	port := os.Args[1]
	filesFolderName := conn.ip.IP.String() + ":" + port
	filesFolderPath := fmt.Sprintf("files/%s/%s", filesFolderName, fileName)

	file, err := os.ReadFile(filesFolderPath)
	log.Println("[SERVER] File read at the ip ", &conn.ip)
	if err != nil {
		return nil, err
	}
	return file, nil
}
