# np

np is a tool for parsing multiple Nmap scans and querying the results.

## Usage

```sh
# Parse all scans in the current directory
$ np
TBD

# Print a specific host
$ np -host <IP|name>
TBD

# Print all alive hosts
$ np -hosts [-path /path/to/scans]
TBD

# Show instances of the given service
$ np -service http [-path /path/to/scans]

# Print all services
$ np -services [-path /path/to/scans]

# Print all hosts with the given port open
$ np -port 80 [-path /path/to/scans]
www.example.com:80 (80.80.21.0)

# Print all hosts with the given ports open
$ np -port 80,443 [-path /path/to/scans]
www.example.com:80 (80.80.21.0)
www.example.com:443 (80.80.21.0)

# Print all open ports
$ np -ports [-path /path/to/scans]
TCP: 80,443,8080,8443
UDP: 53, 161

# Print full JSON dump
$ np [-path /path/to/scans] -json
...
```

## Installation

```sh
go install github.com/leesoh/np/cmd/np@latest
```

## Thanks

- [go-nmap](https://github.com/lair-framework/go-nmap) - For making the XML parsing less awful
- [naabu](https://github.com/projectdiscovery/naabu) - For the beautiful source code
