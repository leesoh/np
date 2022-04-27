package scan

import (
	"bytes"
	"encoding/json"

	"github.com/leesoh/np/internal/result"
)

func (s *Scan) IsNP() bool {
	var h []*result.Host
	dec := json.NewDecoder(bytes.NewReader(s.Bytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&h); err != nil {
		s.Logger.Debugf("not np output")
		return false
	}
	s.Logger.Debugf("found valid np results")
	return true
}

func (s *Scan) ParseNP() {
	var hosts []*result.Host
	// We're importing an old np session so just unpack into results
	err := json.Unmarshal(s.Bytes, hosts)
	if err != nil {
		s.Logger.Errorf("error unmarshaling np results: %v", err)
	}
	for _, hh := range hosts {
		s.Result.AddHost(hh)
	}
}
