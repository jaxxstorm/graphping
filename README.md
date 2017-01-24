# graphping

## Description

Graphping is a tool to ping a list of endpoints and send the results to [statsd](https://github.com/etsy/statsd)

It allows you to create graphs of latency in a similar manner to [smokeping](http://oss.oetiker.ch/smokeping/)

![](http://i.imgur.com/fEuGmTn.png)

Graphping is written in Go and takes advantages of many of Go's features:

  - Built in concurrency
  - Fast
  - Easy to build/install


## Usage

Graphping requires a few options to run. Here's the basic usage:

```
NAME:
   graph-ping - Ping a lost of endpoints and send the resulting metrics to statsd

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.1

AUTHOR(S):
   Lee Briggs

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config-file value, -c value  Path to configuration file
   --statsd value, -s value       Address of statsd listener
   --verbose                      Output metrics in logs
   --help, -h                     show help
   --version, -v                  print the version
```

### Config File

Graphping requires a config file declaring the endpoints you wish to ping. The config format is [hcl](https://github.com/hashicorp/hcl) meaning you can either provide a human readable HCL config file or a JSON config file. An example HCL file looks like this:

```
interval = 10 # A global interval. Can be overwritten per target group
prefix = "graphping" # A global prefix for statsd metrics


# Declare a target group with a name
target_group "search_engines" {
  # a custom ping interval for this group
  interval = 2
  # A prefix for the statsd metric for this group
  prefix = "search"
  # A name for the target. This becomes the statsd metric
  target "google" {
    address = "www.google.co.uk"
  }
  target "bing" {
    address = "www.bing.com"
  }
}

# You can specify multiple target groups
target_group "news_sites" {
  prefix = "uk"
  target "bbc" {
    address = "www.bbc.co.uk"
  }
}
```

### StatsD Listenser

You need to specify the address of the statsd listener you want to send metrics to. This is in string format, including port number.

For example:

```
graphping -c /path/to/config/file.hcl -s 127.0.0.1:8125
``` 

## Building

Currently building it is a bit crude. Simply set your `$GOPATH`: https://github.com/golang/go/wiki/GOPATH

Grab the external dependencies:

```
go get gopkg.in/urfave/cli.v1
go get github.com/Sirupsen/logrus
go get github.com/cactus/go-statsd-client/statsd
go get github.com/tatsushid/go-fastping
```

and then build it:

```
go build main.go
```

## Important Notes

* I am still learning Go, so there may be some absolute nonsense in here. Pull requests are very welcomne

## Acknowledgements

* Thanks to my colleague @cspargo for his initial idea and the first prototype of the code
