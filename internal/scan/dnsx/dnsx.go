package dnsx

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
)

//{
//  "host": "approval-bot.it.epicgames.com",
//  "resolver": [
//    "8.8.8.8:53"
//  ],
//  "a": [
//    "34.225.117.206"
//  ],
//  "status_code": "NOERROR",
//  "timestamp": "2022-05-10T10:54:41.294781206-06:00"
//}

type DNSxScan struct {
	Records []*Record
}

type Record struct {
	IPAddresses []string `json:"a"`
	Hostname    string   `json:"host"`
}

func Parse(scan []byte) (*DNSxScan, error) {
	s := &DNSxScan{}
	// DNSx scan is JSONL, so we have to read it line by line
	sc := bufio.NewScanner(bytes.NewReader(scan))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		r := &Record{}
		err := json.Unmarshal(sc.Bytes(), &r)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling DNSx: %v", err)
		}
		// Incomplete data, move on
		if isMissingKeys(r) {
			continue
		}
		s.Records = append(s.Records, r)
	}
	return s, nil
}

// isMissingKeys attempts to filter non-DNSx JSON files (e.g. Amass output)
func isMissingKeys(r *Record) bool {
	// Every host should have both of these
	if len(r.IPAddresses) == 0 || r.Hostname == "" {
		return true
	}
	return false
}
