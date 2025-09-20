package tcpclient

import (
	"fmt"
	"net"
)

func Search(addr *net.TCPAddr, searchId string, fileName string, depth uint32) (*net.TCPAddr, error) {
	req := Request{
		Action: "search",
		Params: map[string]interface{}{
			"searchId": searchId,
			"fileName": fileName,
			"depth":    depth,
		},
	}

	resp, err := sendRequest(addr, req)
	if err != nil {
		return nil, err
	}
	if resp.Status != "ok" {
		return nil, fmt.Errorf("[CLIENT] server error: %s", resp.Error)
	}

	// Cast response data
	ipWithFile, ok := resp.Data.(*net.TCPAddr)
	if !ok {
		return nil, fmt.Errorf("[CLIENT] invalid response type")
	}

	return ipWithFile, nil
}

// Add more requests here
