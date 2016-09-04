# APINATE [![Build Status](https://travis-ci.org/PhilipHarries/apinate.svg?branch=master)](https://travis-ci.org/PhilipHarries/apinate)

# Configuration-driven arbitrary APIs

## Installation

go get github.com/mattn/gom
gom install
gom build -o /usr/bin/apinate

## Development

go run main.go

## Testing

run_tests.sh
or
docker build .

# Configuration

*apinate* may be configured by toml, json or yaml, it will check in turn for:
 - ~/.apinate.toml
 - ~/.apinate.json
 - ~/.apinate.yaml
 - /etc/apinate.toml
 - /etc/apinate.json
 - /etc/apinate.yaml

Configuration is from resource (url) to command.  If additional parameters are to be passed to the command, the mapping should be passed a boolean called "params" set to true.

The api can be directed to output json, html, yaml or text.

### toml
```
contenttype = "json"
[[mappings]]
  resource = "echo"
  command = "echo"
  params = true
[[mapping]]
  resource = "ping"
  command = "ping"
  params = true
[[mapping]]
  resource = hostname
  command = "hostname -f"
  params = false
```

### json
```
{
  "contenttype": "html"
  "mappings": [
    {
      "resource": "echo",
      "command":  "echo",
      "params":   true
    },
    {
      "resource": "ping",
      "command":  "ping",
      "params":   true
    },
    {
      "resource": "hostname",
      "command":  "hostname -f",
      "params":   false
    }
  ]
}
```

### yaml
```
contenttype: txt
mappings:
  - resource: echo
    command:  echo
    params:   true
  - resource: ping
    command:  ping
    params:   true
  - resource: hostname
    command:  hostname -f
    params:   false
```

# Usage

The api may be accessed simply via HTTP requests:

- http://hostname:8080/command
- http://hostname:8080/command/parameters

