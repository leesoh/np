package naabu

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type NaabuScan struct {
	V1Hosts []*HostV1
	V2Hosts []*HostV2
}

type HostV1 struct {
	Name      string `json:"host"`
	IPAddress string `json:"ip"`
	Port      int    `json:"port"`
}

type HostV2 struct {
	Name      string    `json:"host"`
	IPAddress string    `json:"ip"`
	Port      Port      `json:"port"`
	Timestamp time.Time `json:"timestamp"`
}

type Port struct {
	Port     int  `json:"Port"`
	Protocol int  `json:"Protocol"`
	TLS      bool `json:"TLS"`
}

func ParseV1(scan []byte) (*NaabuScan, error) {
	s := &NaabuScan{}
	// Naabu scan is JSONL, so we have to read it line by line
	sc := bufio.NewScanner(bytes.NewReader(scan))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		hv1 := &HostV1{}
		err := json.Unmarshal(sc.Bytes(), &hv1)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling naabu V1: %v", err)
		}
		if isMissingV1Keys(hv1) {
			return nil, fmt.Errorf("missing Naabu keys")
		}
		s.V1Hosts = append(s.V1Hosts, hv1)
	}
	return s, nil
}

func ParseV2(scan []byte) (*NaabuScan, error) {
	s := &NaabuScan{}
	// Naabu scan is JSONL, so we have to read it line by line
	sc := bufio.NewScanner(bytes.NewReader(scan))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		hv2 := &HostV2{}
		err := json.Unmarshal(sc.Bytes(), &hv2)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling naabu V2: %v", err)
		}
		if isMissingV2Keys(hv2) {
			return nil, fmt.Errorf("missing Naabu keys")
		}
		s.V2Hosts = append(s.V2Hosts, hv2)
	}
	return s, nil
}

// isMissingV1Keys attempts to filter non-Naabu JSON files (e.g. Amass output)
func isMissingV1Keys(h *HostV1) bool {
	// Every host should have one of these
	if h.IPAddress == "" || h.Port == 0 {
		return true
	}
	return false
}

// isMissingV2Keys attempts to filter non-Naabu JSON files (e.g. Amass output)
func isMissingV2Keys(h *HostV2) bool {
	// Every host should have one of these
	if h.IPAddress == "" || h.Port.Port == 0 {
		return true
	}
	return false
}
