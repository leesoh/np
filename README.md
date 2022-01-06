# np

np is a tool for parsing multiple Nmap scans and querying the results. It can also import its own JSON results.

## Usage

```text
Usage of np:
  -host string
     Show results for specified host
  -hosts
     Print alive hosts
  -json
     Display JSON output
  -path string
     Path to scan file (default ".")
  -port string
     Display hosts with matching port(s)
  -ports
     Print all ports
  -service string
     Display hosts with matching service name
  -services
     Print all services
  -verbose
     Display verbose output
```

The `-json` option will display all hosts with at least one open port, while the other formats will omit ports that are likely false positives (e.g. `tcpwrapped`).

## Examples

Parse all scans in the current directory:

```sh
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
```

Print a specific host:

```sh
$ np -host scanme.nmap.org
scanme.nmap.org (45.33.32.156)
PORT      SERVICE    PRODUCT      VERSION
80/tcp    http       Apache httpd 2.4.7
22/tcp    ssh        OpenSSH      6.6.1p1 Ubuntu 2ubuntu2.13
9929/tcp  nping-echo Nping echo   
31337/tcp Elite                   
```

Print all alive hosts:

```sh
$ np -hosts [-path /path/to/scans]
45.33.32.156 (scanme.nmap.org)
159.223.119.162 (bishopfox.com)
```

Print all services:

```sh
$ np -services [-path /path/to/scans]
Elite
http
https
nping-echo
ssh
```

Show instances of the given service:

```sh
$ np -service https [-path /path/to/scans]
bishopfox.com (159.223.119.162)
PORT    SERVICE PRODUCT VERSION
80/tcp  http            
443/tcp https
```

Print all open ports:

```sh
$ np -ports [-path /path/to/scans]
22,80,443,9929,31337
```

Print all hosts with the given port open:

```sh
$ np -port 80 [-path /path/to/scans]
45.33.32.156:80
159.223.119.162:80
```

Print all hosts with the given ports open:

```sh
$ np -port 80,443 [-path /path/to/scans]
45.33.32.156:80
159.223.119.162:80
159.223.119.162:443
```

Print full JSON dump:

```sh
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

