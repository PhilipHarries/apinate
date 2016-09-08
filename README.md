# APINATE [![Build Status](https://travis-ci.org/PhilipHarries/apinate.svg?branch=master)](https://travis-ci.org/PhilipHarries/apinate)

# Configuration-driven arbitrary APIs

*apinate* is a simple HTTP API that can be configured to run any back end script or command and display the output in a variety of formats.  It runs as a single binary, and is driven by a single configuration file.

*apinate* can output JSON, YAML or raw text, and can output configurably templated HTML.

*apinate* is the simplest and easiest way of spinning up a functional API.

### Caution

Clearly *apinate* has the potential to pose a significant security risk as it can be used to run arbitrary commands with the privileges assigned to the user it runs as.  Do not, for example, link *apinate* to a command such as rm, or even cat, unless you have taken appropriate precautions.

*apinate* should not be used with long (or continually) running commands, as the API endpoint will simply not return.

## Installation

```
go get github.com/mattn/gom
gom install
gom build -o /usr/bin/apinate
```

## Development

```
go run main.go
```

## Testing

```
run_tests.sh
```
or
```
docker build .
```

# Configuration

*apinate* may be configured by toml, json or yaml, it will check in turn for:
 - ~/.apinate/apinate.toml
 - ~/.apinate/apinate.json
 - ~/.apinate/apinate.yaml
 - /etc/apinate/apinate.toml
 - /etc/apinate/apinate.json
 - /etc/apinate/apinate.yaml

Configuration is from resource (url) to command.  If additional parameters are to be passed to the command, the mapping should be passed a boolean called "params" set to true.

The api can be directed to output json, html, yaml or raw text.

address, port, params, logfile and template directives are optional and default to 0.0.0.0, 8080, false, stderr and plain.tmpl

Supported options for logfile are stderr or a filename.

#### HTML output

If HTML output is specified, you can create your own template files in ~/.apinate/templates or /usr/share/apinate/templates, and specify which to use in the mapping configuration with the template directive

### toml
```
address = "0.0.0.0"
port = 8080
contenttype = "json"
logfile = "~/.apinate/apinate.out"
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
  template = "new.tmpl"
```

### json
```
{
  "address": "0.0.0.0"
  "port": 8080
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
address: 0.0.0.0
port: 8080
contenttype: raw
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

# Examples

In order to report on the systems hostname, and the time taken to ping parameter-driven locations, configured by json, outputting in xml:

### ~/.apinate.json
```
{
  "contenttype": "html",
  "mappings": [
    {
      "resource": "system-name",
      "command":  "hostname",
      "params":   false
    },
    {
      "resource": "ping-time",
      "command":  "ping -c 1 ",
      "params":   true
    }
  ]
}
```
### output
```
curl http://localhost:8080/system-name
<html>
  <head></head>
  <body>
        <p>myhostname</p>
  </body>
</html>
curl http://localhost:8080/ping-time/www.google.com
<html>
  <head></head>
  <body>
        <p>PING www.google.com (216.58.208.132) 56(84) bytes of data.</p>
        <p>64 bytes from lhr25s08-in-f132.1e100.net (216.58.208.132): icmp_seq=1 ttl=57 time=18.4 ms</p>
        <p></p>
        <p>--- www.google.com ping statistics ---</p>
        <p>1 packets transmitted, 1 received, 0% packet loss, time 0ms</p>
        <p>rtt min/avg/max/mdev = 18.495/18.495/18.495/0.000 ms</p>
  </body>
</html>
```

# Build / Packaging

build.sh will use fpm to build deb and rpm packages, call it with `build.sh <version>` where version should be the semver release version of apinate - e.g. 1.0.0

build.sh relies on fpm: https://github.com/jordansissel/fpm/wiki

