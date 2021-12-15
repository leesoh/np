# np

np is a tool for parsing multiple Nmap scans and querying the results.

## Usage

The `-json` option will display all hosts, while the other formats will omit hosts with no open ports, or ports that are likely false positives (e.g. `tcpwrapped`).

```sh
# Parse all scans in the current directory
$ np
scanme.nmap.org (45.33.32.156)
PORT      SERVICE    PRODUCT      VERSION
80/tcp    http       Apache httpd 2.4.7
22/tcp    ssh        OpenSSH      6.6.1p1 Ubuntu 2ubuntu2.13
9929/tcp  nping-echo Nping echo   
31337/tcp Elite                   

bishopfox.com (159.223.119.162)
PORT    SERVICE PRODUCT VERSION
80/tcp  http            
443/tcp https

# Print a specific host
$ np -host scanme.nmap.org
scanme.nmap.org (45.33.32.156)
PORT      SERVICE    PRODUCT      VERSION
80/tcp    http       Apache httpd 2.4.7
22/tcp    ssh        OpenSSH      6.6.1p1 Ubuntu 2ubuntu2.13
9929/tcp  nping-echo Nping echo   
31337/tcp Elite                   

# Print all alive hosts
$ np -hosts [-path /path/to/scans]
45.33.32.156 (scanme.nmap.org)
159.223.119.162 (bishopfox.com)

# Print all services
$ np -services [-path /path/to/scans]
Elite
http
https
nping-echo
ssh

# Show instances of the given service
$ np -service https [-path /path/to/scans]
bishopfox.com (159.223.119.162)
PORT    SERVICE PRODUCT VERSION
80/tcp  http            
443/tcp https

# Print all open ports
$ np -ports [-path /path/to/scans]
22,80,443,9929,31337

# Print all hosts with the given port open
$ np -port 80 [-path /path/to/scans]
45.33.32.156:80
159.223.119.162:80

# Print all hosts with the given ports open
$ np -port 80,443 [-path /path/to/scans]
45.33.32.156:80
159.223.119.162:80
159.223.119.162:443

# Print full JSON dump
$ np [-path /path/to/scans] -json
[
  {
    "ip": "45.33.32.156",
    "hostname": "scanme.nmap.org",
    "tcp_ports": {
      "22": {
        "name": "ssh",
        "product": "OpenSSH",
        "version": "6.6.1p1 Ubuntu 2ubuntu2.13",
        "extra_info": "Ubuntu Linux; protocol 2.0"
      },
      "31337": {
        "name": "Elite"
      },
      "80": {
        "name": "http",
        "product": "Apache httpd",
        "version": "2.4.7",
        "extra_info": "(Ubuntu)"
      },
      "9929": {
        "name": "nping-echo",
        "product": "Nping echo"
      }
    }
  },
  {
    "ip": "159.223.119.162",
    "hostname": "bishopfox.com",
    "tcp_ports": {
      "443": {
        "name": "https"
      },
      "80": {
        "name": "http"
      }
    }
  }
]
```

## Installation

```sh
go install github.com/leesoh/np/cmd/np@latest
```

## Thanks

- [go-nmap](https://github.com/lair-framework/go-nmap) - For making the XML parsing less awful
- [naabu](https://github.com/projectdiscovery/naabu) - For the beautiful source code
