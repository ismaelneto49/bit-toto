package tcpclient

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
)

func Search(addr *net.TCPAddr, searchId string, fileName string, depth uint32) (*net.TCPAddr, error) {
	req := Request{
		Action: "search",
		Params: map[string]interface{}{
			"searchId": searchId,
			"fileName": fileName,
			"depth":    depth,
		},
		ClientPort: os.Args[1],
	}

	resp, err := sendRequest(addr, req)
	if err != nil {
		return nil, err
	}
	if resp.Status != "ok" {
		return nil, fmt.Errorf("[CLIENT] server error: %s", resp.Error)
	}

	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response type")
	}

	ipStr, ok1 := dataMap["ip"].(string)
	portF, ok2 := dataMap["port"].(float64) // JSON numbers are float64
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid IP/Port data")
	}

	ipWithFile := &net.TCPAddr{
		IP:   net.ParseIP(ipStr),
		Port: int(portF),
	}
	return ipWithFile, nil
}

func Download(addr *net.TCPAddr, fileName string) ([]byte, error) {
	req := Request{
		Action: "download",
		Params: map[string]interface{}{
			"fileName": fileName,
		},
		ClientPort: os.Args[1],
	}

	resp, err := sendRequest(addr, req)
	if err != nil {
		return nil, err
	}
	if resp.Status != "ok" {
		return nil, fmt.Errorf("[CLIENT] server error: %s", resp.Error)
	}

	fileString, ok := resp.Data.(string)
	if !ok {
		return nil, fmt.Errorf("[CLIENT] invalid file data")
	}

	// decode base64 file
	decodedFile, err := base64.StdEncoding.DecodeString(fileString)
	if err != nil {
		return nil, err
	}
	return decodedFile, nil
}

// Add more requests here
