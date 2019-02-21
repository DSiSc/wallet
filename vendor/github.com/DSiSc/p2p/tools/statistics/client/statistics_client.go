package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/p2p/tools/common"
	"io/ioutil"
	"net/http"
)

// StatisticsClient statistics client
type StatisticsClient struct {
	Server string
}

// NewStatisticsClient create a statistics client instance
func NewStatisticsClient(serverAddr string) *StatisticsClient {
	return &StatisticsClient{
		Server: "http://" + serverAddr,
	}
}

// ReportNeighbors report local neighbor to server
func (this *StatisticsClient) ReportNeighbors(localAddr string, neighbors []*common.Neighbor) error {
	msgByte, err := json.Marshal(neighbors)
	if err != nil {
		log.Warn("failed to encode peer neighbors, as: %v", err)
		return fmt.Errorf("failed to encode peer neighbors, as: %v", err)
	}
	resp, err := http.Post(this.Server+"/neighbor/"+localAddr, "application/json", bytes.NewReader(msgByte))
	if err != nil {
		log.Warn("failed to send neighbor info to display server, as:%v", err)
		return fmt.Errorf("failed to send neighbor info to display server, as:%v", err)
	}
	defer resp.Body.Close()
	var respStr string
	return parseResp(&respStr, resp)
}

// ReportNeighbors report local neighbor to server
func (this *StatisticsClient) GetTopos() (map[string][]*common.Neighbor, error) {
	resp, err := http.Get(this.Server + "/topos")
	if err != nil {
		log.Warn("failed to get topo info from server as: %v", err)
		return nil, fmt.Errorf("failed to get topo info from server as: %v", err)
	}
	defer resp.Body.Close()
	var topos map[string][]*common.Neighbor
	err = parseResp(&topos, resp)
	if err != nil {
		return nil, err
	} else {
		return topos, nil
	}
}

// ReportMsg report received/broadcast message to server
func (this *StatisticsClient) ReportMsg(msgId string, msg *common.ReportMsg) error {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Warn("failed to encode message, as: %v", err)
		return fmt.Errorf("failed to encode message, as: %v", err)
	}
	resp, err := http.Post(this.Server+"/msgs/0x"+msgId, "application/json", bytes.NewReader(msgByte))
	if err != nil {
		log.Warn("failed to send message info to display server, as:%v", err)
		return fmt.Errorf("failed to send message info to display server, as:%v", err)
	}
	defer resp.Body.Close()
	var respStr string
	return parseResp(&respStr, resp)
}

// GetReportMessage get report message from server
func (this *StatisticsClient) GetReportMessage(msgId string) (map[string]*common.ReportMsg, error) {
	resp, err := http.Get(this.Server + "/msgs/" + msgId)
	if err != nil {
		log.Warn("failed to get msg info from server as: %v", err)
		return nil, fmt.Errorf("failed to get msg info from server as: %v", err)
	}
	defer resp.Body.Close()
	var repMsg map[string]*common.ReportMsg
	err = parseResp(&repMsg, resp)
	if err != nil {
		return nil, err
	} else {
		return repMsg, nil
	}
}

// GetReportMessage get report message from server
func (this *StatisticsClient) GetAllReportMessage() (map[string]map[string]*common.ReportMsg, error) {
	resp, err := http.Get(this.Server + "/msgs")
	if err != nil {
		log.Warn("failed to get msgs info from server as: %v", err)
		return nil, fmt.Errorf("failed to get msgs info from server as: %v", err)
	}
	defer resp.Body.Close()
	var repMsgs map[string]map[string]*common.ReportMsg
	err = parseResp(&repMsgs, resp)
	if err != nil {
		return nil, err
	} else {
		return repMsgs, nil
	}
}

// unmarshal response to specified type
func parseResp(v interface{}, resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response body, as: %v", err)
		return fmt.Errorf("Failed to read response body, as: %v", err)
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		log.Error("Failed to parse response body, as: %v", err)
		return fmt.Errorf("Failed to parse response body, as: %v", err)
	}
	return nil
}

// TopoReachability calculate the max reachability of this net
func TopoReachbility(topo map[string][]*common.Neighbor) int {
	maxReachableCount := 0
	for node, _ := range topo {
		count := reachableCount(node, topo)
		if count > maxReachableCount {
			maxReachableCount = count
		}
	}
	return maxReachableCount
}

// calculate the max reachability of this net from the specified node
func reachableCount(startNode string, topo map[string][]*common.Neighbor) int {
	reachableNodes := make(map[string]interface{})

	unVisitedNodes := make(map[string]interface{})
	unVisitedNodes[startNode] = struct{}{}
	for {
		if len(unVisitedNodes) <= 0 {
			break
		}
		neighbors := make(map[string]interface{})
		for node, _ := range unVisitedNodes {
			for _, neighbor := range topo[node] {
				if reachableNodes[neighbor.Address] == nil {
					reachableNodes[neighbor.Address] = struct{}{}
					neighbors[neighbor.Address] = struct{}{}
				}
			}
		}
		unVisitedNodes = neighbors
	}
	return len(reachableNodes)
}
