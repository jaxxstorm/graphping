package main

import (
	"./config"
	"./ping"
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config-file, c", Usage: "Path to configuration file"},
		cli.StringFlag{Name: "statsd, s", Usage: "Address of statsd listener"},
	}

	app.Name = "graph-ping"
	app.Version = "0.1"
	app.Usage = "Ping a lost of endpoints and send the resulting metrics to statsd"
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Lee Briggs",
		},
	}

	app.Action = func(c *cli.Context) error {

		// Check we have a config file
		if !c.IsSet("config-file") {
			cli.ShowAppHelp(c)
			return cli.NewExitError("Error: No config file specified", -1)
		}

		if !c.IsSet("statsd") {
			cli.ShowAppHelp(c)
			return cli.NewExitError("Error: No statsd client specified", -1)
		}

		// if we can't parse it, error
		config, err := config.Parse(c.String("config-file"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error: Unable to parse config file - %s", err), -1)
		} else {

			// everything is fine, let's start pinging!
			// create a statsdClient

			statsdClient, err := statsd.NewClient(c.String("statsd"), config.Prefix)

			// Issues opening statsd
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Error: Unable to open statsd client - %s, err"), -1)
			}

			// defer till later
			defer statsdClient.Close()

			// loop through the groups and start a goroutine
			// for each group to ping the targets
			for _, groups := range config.Groups {
				go ping.RunPinger(config.Interval, statsdClient, groups)
			}

			// channel handling for interrupting app
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			signal.Notify(c, syscall.SIGTERM)
			done := make(chan bool, 1)

			go func() {
				for sig := range c {
					log.Printf("captured %v. ", sig)
					done <- true
				}
			}()
			<-done
		}
		return nil
	}

	// run!
	app.Run(os.Args)

}
