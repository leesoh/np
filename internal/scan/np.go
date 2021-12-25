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
		s.Logger.Debugf("not np results")
		return false
	}
	s.Logger.Debugf("found valid np results")
	return true
}

func (s *Scan) ParseNP() {
	err := json.Unmarshal(s.Bytes, &s.Result.Hosts)
	if err != nil {
		s.Logger.Errorf("error unmarshaling np results: %v", err)
	}
}
