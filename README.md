# np

np is a tool for parsing, combining, deduplicating, and querying the output from multiple different tools.

Supported scan types:

- Nmap XML output
- Masscan XML output
- Naabu JSON output
- DNSx JSON output
- `np` JSON output

## Usage

```text
Usage of ./np:
  -debug
     Display debug output
  -exclude string
     Exclude these hosts from output
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
  -timeline
     Display timeline output
```

The `-json` option will display all hosts with no filtering. Other formats omit ports that are likely false positives (e.g. `tcpwrapped`).

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

Parse all scans, but ignore a host:

```sh
$ np -exclude bishopfox.com
scanme.nmap.org (45.33.32.156)
PORT      SERVICE    PRODUCT      VERSION
80/tcp    http       Apache httpd 2.4.7
22/tcp    ssh        OpenSSH      6.6.1p1 Ubuntu 2ubuntu2.13
9929/tcp  nping-echo Nping echo
31337/tcp Elite
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
scanme.nmap.org (45.33.32.156)
bishopfox.com (159.223.119.162)
```

Print all services:

```sh
$ np -services [-path /path/to/scans]
scanme.nmap.org:9929 nping-echo
scanme.nmap.org:31337 Elite
scanme.nmap.org:22 ssh
scanme.nmap.org:80 http
bishopfox.com:80 http
bishopfox.com:443 https
```

Show instances of the given service:

```sh
$ np -service https [-path /path/to/scans]
bishopfox.com:443
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

Make sure that `$GOPATH/bin` is part of `$PATH`, then:

```sh
go install github.com/leesoh/np/cmd/np@latest
```

## Similar Tools

In case `np` isn't quite what you're looking for, here are a few similar tools:

- [nmap-parse-output](https://github.com/ernw/nmap-parse-output) - Parse single Nmap scan output with extremely flexible output
- [nmap-bootstrap-xsl](https://github.com/honze-net/nmap-bootstrap-xsl/) - Create nice looking HTML reports from Nmap scan output

## Thanks

- [go-nmap](https://github.com/lair-framework/go-nmap) - For making the XML parsing less awful
- [naabu](https://github.com/projectdiscovery/naabu) - For the beautiful source code

