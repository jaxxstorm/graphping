package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/jaxxstorm/graphping/config"
	"github.com/jaxxstorm/graphping/ping"
	"gopkg.in/urfave/cli.v1"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config-file, c", Usage: "Path to configuration file"},
		cli.StringFlag{Name: "statsd, s", Usage: "Address of statsd listener"},
		cli.BoolFlag{Name: "verbose", Usage: "Output metrics in logs"},
	}

	app.Name = "graphping"
	app.Version = "0.1"
	app.Usage = "Ping a list of endpoints and send the resulting metrics to statsd"
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

		// if debug flag is set, override loglevel:
		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		// if we can't parse it, error
		config, err := config.Parse(c.String("config-file"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error: Unable to parse config file - %s", err), -1)
		} else {

			// everything is fine, let's start pinging!
			// create a statsdClient

			statsdClient, err := statsd.NewClient(c.String("statsd"), config.Prefix)
			log.Debug("Global StatsD Prefix: ", config.Prefix)

			// Issues opening statsd
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Error: Unable to open statsd client - %s, err"), -1)
			}

			// defer till later
			defer statsdClient.Close()

			done := make(chan bool, 1)

			// loop through the groups and start a goroutine
			// for each group to ping the targets
			for _, groups := range config.Groups {
				go func() {
					pingres := ping.RunPinger(config.Interval, statsdClient, groups)
					log.Warn(fmt.Sprintf("Pinger exited: %s", pingres))
					done <- true
				}()
			}

			// channel handling for interrupting app
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			signal.Notify(c, syscall.SIGTERM)

			go func() {
				for sig := range c {
					log.Warn(fmt.Sprintf("captured %v", sig))
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
