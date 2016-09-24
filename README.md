# APINATE [![Build Status](https://travis-ci.org/PhilipHarries/apinate.svg?branch=master)](https://travis-ci.org/PhilipHarries/apinate)

# Configuration-driven, single-binary arbitrary APIs

*apinate* is a simple HTTP API that can be configured to run any back end script or command, with parameters, and display the output in a variety of formats.  It runs as a single binary, and is driven by a single configuration file.

*apinate* can output JSON, YAML or raw text, and can output configurably templated HTML.

*apinate* is the simplest and easiest way of spinning up a functional API.

### Caution

Clearly a poorly configured *apinate* has the potential to pose a significant security risk as it can be used to run arbitrary commands with the privileges assigned to the user it runs as.  Do not, for example, link *apinate* to a command such as rm, or even cat, unless you have taken appropriate precautions.  If security is a concern, *apinate* should be linked to commands wrapped by scripts that validate inputs.

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

Configuration is from resource (url) to command.

Parameters may be passed to the command in two ways.  Typically, URL's are passed parameters via a query string such as:

"http://www.example.com/resource?param=something"

This is a standard and very configurable mechanism for passing information into an API, and is supported by apinate.

More than one parameter can be passed as follows:

"http://www.example.com/resource?param1=something&param2=somethingelse&param3=somethingelseagain"

These are ultimately passed to the script or command being run in the same format as they are specified in the URL, so that "param1=something" and "param2=somethingelse" from the above URL would be passed as the first and second parameters to the script.  The script therefore requires logic to split the parameters before making use of them.

In many cases it will not be possible to split the parameters, particularly if apinate is being used to create an API for a command or application that is not in your control.  In this case, parameters may be passed through directly by setting "altparams" to true.  The URL can then be called as:

"http://www.example.com/resource/something somethingelse somethingelseagain"
 - this will call the defined script and pass "something somethingelse somethingelseagain" as positional parameters.  If calling from a browser, the spaces should be translated into "%20" (the URL-encoding for a space)

Only one of "querystrings" or "altparams" may be specified for one resource - specifying both querystrings and altparams as true will result in an error and apinate will not start.

The api can be directed to output json, html, yaml or raw text.

address, port, altparams, querykeys, logfile and template directives are optional and default to:

address:   0.0.0.0
port:      8080
altparams: false
querykeys: [] (an empty array) 
logfile:   stdout
template:  plain.tmpl

logfile can take a filename, or the directive "stdout" will direct output to standard error.

querykeys takes an array of queries, which have a keyname, and an optional default value.  The format for defining these can be found in the respective configuration sections below.

#### HTML output

If HTML output is specified, template files can be added in ~/.apinate/templates or /usr/share/apinate/templates, and specify which to use in the mapping configuration with the template directive

### toml
```
address = "0.0.0.0"
port = 8080
contenttype = "json"
logfile = "~/.apinate/apinate.out"
[[mappings]]
  resource = "echo"
  command = "echo"
  altparams = true
[[mappings]]
  resource = "ping"
  command = "ping"
  altparams = true
[[mappings]]
  resource = hostname
  command = "hostname -f"
  altparams = false
  template = "new.tmpl"
[[mappings]]
  resource = myscript
  command = "/usr/local/bin/myscript"
  [[ mappings.querykeys ]]
    keyname = "key1"
    default = "value1"
  [[ mappings.querykeys ]]
    keyname = "key2"
```

### json
```
{
  "address": "0.0.0.0"
  "port": 8080
  "contenttype": "html"
  "mappings": [
    {
      "resource":  "echo",
      "command":   "echo",
      "altparams": true
    },
    {
      "resource": "ping",
      "command":  "ping",
      "altparams": false,
      "querykeys": [
        "keyname": "host",
        "default": "127.0.0.1"
      ]
    },
    {
      "resource":  "hostname",
      "command":   "hostname -f",
      "altparams": false
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
  - resource:  echo
    command:   echo
    altparams: true
  - resource:  ping
    command:   ping
    altparams: true
  - resource:  hostname
    command:   hostname -f
    altparams: false
```

# Usage

The api may be accessed simply via HTTP requests:

- http://hostname:8080/command
- http://hostname:8080/command/parameters

# Example

The following configuration will:
    report the system hostname at /hostname
    record the time taken to ping parameter-driven locations at /pint-time
    pass parameters val1, val2 and an operator (+,-,/,*) to a calculator script
 - this example is configured by json to output in html

### ~/.apinate.json
```
{
  "contenttype": "html",
  "mappings": [
    {
      "resource": "system-name",
      "command":  "hostname",
      "altparams":   false
    },
    {
      "resource": "ping-time",
      "command":  "ping -c 1 ",
      "altparams":   true
    },
    {
      "resource": "calc",
      "command":  "/usr/local/bin/mycalc",
      "querykeys": [
        {
          "keyname":  "val1",
          "default":  "0"
        },
        {
          "keyname":  "val2",
          "default":  "0"
        }  ,
        {
          "keyname":  "operator",
          "default":  "+"
        }  
      ]
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

curl http://localhost:8080/calc?val1=10&val2=5&operator=\* 
<html>
  <head></head>
  <body>
        <p>50</p>
  </body>
</html>
```



# Build / Packaging

build.sh will use fpm to build deb and rpm packages, call it with `build.sh <version>` where version should be the semver release version of apinate - e.g. 1.0.0

build.sh relies on fpm: https://github.com/jordansissel/fpm/wiki

