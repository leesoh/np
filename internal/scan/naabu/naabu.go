package naabu

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
)

type NaabuScan struct {
	Hosts []*Host
}

type Host struct {
	Name      string `json:"host"`
	IPAddress string `json:"ip"`
	Port      int    `json:"port"`
}

func Parse(scan []byte) (*NaabuScan, error) {
	s := &NaabuScan{}
	// Naabu scan is JSONL, so we have to read it line by line
	sc := bufio.NewScanner(bytes.NewReader(scan))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		h := &Host{}
		err := json.Unmarshal(sc.Bytes(), &h)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling Naabu: %v", err)
		}
		if isMissingKeys(h) {
			return nil, fmt.Errorf("missing Naabu keys")
		}
		s.Hosts = append(s.Hosts, h)
	}
	return s, nil
}

// isMissingKeys attempts to filter non-Naabu JSON files (e.g. Amass output)
func isMissingKeys(h *Host) bool {
	// Every host should have one of these
	if h.IPAddress == "" || h.Port == 0 {
		return true
	}
	return false
}
