package nmap

import (
	"encoding/xml"
	"fmt"
)

// Thanks to lair-framework/go-nmap for making this simpler
type NmapScan struct {
	XMLName          xml.Name  `xml:"nmaprun"`
	Text             string    `xml:",chardata"`
	Scanner          string    `xml:"scanner,attr"`
	Args             string    `xml:"args,attr"`
	Start            string    `xml:"start,attr"`
	Startstr         string    `xml:"startstr,attr"`
	Version          string    `xml:"version,attr"`
	XMLOutputVersion string    `xml:"xmloutputversion,attr"`
	ScanInfo         ScanInfo  `xml:"scaninfo"`
	Verbose          Verbose   `xml:"verbose"`
	Debugging        Debugging `xml:"debugging"`
	Hosts            []Host    `xml:"host"`
	RunStats         RunStats  `xml:"runstats"`
}

type ScanInfo struct {
	Text        string `xml:",chardata"`
	Type        string `xml:"type,attr"`
	Protocol    string `xml:"protocol,attr"`
	NumServices string `xml:"numservices,attr"`
	Services    string `xml:"services,attr"`
}

type Verbose struct {
	Text  string `xml:",chardata"`
	Level string `xml:"level,attr"`
}

type Debugging struct {
	Text  string `xml:",chardata"`
	Level string `xml:"level,attr"`
}

type RunStats struct {
	Text     string   `xml:",chardata"`
	Finished Finished `xml:"finished"`
	Hosts    Hosts    `xml:"hosts"`
}

type Finished struct {
	Text    string `xml:",chardata"`
	Time    string `xml:"time,attr"`
	Timestr string `xml:"timestr,attr"`
	Elapsed string `xml:"elapsed,attr"`
	Summary string `xml:"summary,attr"`
	Exit    string `xml:"exit,attr"`
}

type Hosts struct {
	Text  string `xml:",chardata"`
	Up    string `xml:"up,attr"`
	Down  string `xml:"down,attr"`
	Total string `xml:"total,attr"`
}

type Host struct {
	Text      string     `xml:",chardata"`
	StartTime string     `xml:"starttime,attr"`
	EndTime   string     `xml:"endtime,attr"`
	Status    Status     `xml:"status"`
	Address   Address    `xml:"address"`
	Hostnames []Hostname `xml:"hostnames>hostname"`
	Ports     []Port     `xml:"ports>port"`
	Times     Times      `xml:"times"`
}

type Status struct {
	Text      string `xml:",chardata"`
	State     string `xml:"state,attr"`
	Reason    string `xml:"reason,attr"`
	ReasonTTL string `xml:"reason_ttl,attr"`
}

type Address struct {
	Text     string `xml:",chardata"`
	Addr     string `xml:"addr,attr"`
	Addrtype string `xml:"addrtype,attr"`
}

type Hostname struct {
	Text string `xml:",chardata"`
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type Port struct {
	Text     string  `xml:",chardata"`
	Protocol string  `xml:"protocol,attr"`
	PortID   string  `xml:"portid,attr"`
	State    State   `xml:"state"`
	Service  Service `xml:"service"`
}

type State struct {
	Text      string `xml:",chardata"`
	State     string `xml:"state,attr"`
	Reason    string `xml:"reason,attr"`
	ReasonTTL string `xml:"reason_ttl,attr"`
}

type Service struct {
	Conf      string `xml:"conf,attr"`
	Method    string `xml:"method,attr"`
	Name      string `xml:"name,attr"`
	Product   string `xml:"product,attr"`
	Version   string `xml:"version,attr"`
	ExtraInfo string `xml:"extrainfo,attr"`
	Text      string `xml:",chardata"`
}

type Times struct {
	Text   string `xml:",chardata"`
	SRTT   string `xml:"srtt,attr"`
	RTTVar string `xml:"rttvar,attr"`
	To     string `xml:"to,attr"`
}

func Parse(scan []byte) (*NmapScan, error) {
	s := &NmapScan{}
	err := xml.Unmarshal(scan, s)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling Nmap: %v", err)
	}
	return s, nil
}
