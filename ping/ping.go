package ping

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/jaxxstorm/graphping/config"
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
)

type response struct {
	addr *net.IPAddr
	rtt  time.Duration
}

func RunPinger(interval int, prefix string, statsdClient statsd.Statter, group config.TargetGroups) error {

	// create a new pinger
	p := fastping.NewPinger()

	// create some interfaces to store results
	results := make(map[string]*response)
	index := make(map[string]string)

	// loop througb the targets in the struct and resolve the address
	for _, target := range group.Targets {
		ra, err := net.ResolveIPAddr("ip4:icmp", target.Address)
		log.Debug(fmt.Sprintf("Target %s resolves to %s", target.Address, ra))
		if err != nil {
			return errors.New(fmt.Sprintf("Can't resolve %s", target.Address))
		}

		// store the result of each ping poll
		results[ra.String()] = nil

		// map the ip address back to the label
		index[ra.String()] = target.Label

		// add the IP address to the list of ping endpoints
		p.AddIPAddr(ra)
	}

	onRecv, onIdle := make(chan *response), make(chan bool)
	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		onRecv <- &response{addr: addr, rtt: t}
	}
	p.OnIdle = func() {
		onIdle <- true
	}

	//determine what interval we should run at
	if group.Interval > 0 {
		log.Debug(fmt.Sprintf("Group Interval defined: %v", group.Interval))
		p.MaxRTT = time.Duration(group.Interval) * time.Second
	} else if interval > 0 {
		log.Debug(fmt.Sprintf("Global Interval defined: %v", interval))
		p.MaxRTT = time.Duration(interval) * time.Second
	} else {
		log.Debug(fmt.Sprintf("Using Default Interval: 60"))
		p.MaxRTT = 60 * time.Second
	}

	if (group.Prefix != "") {prefix = group.Prefix}

	// set the metric path
	metricPath := fmt.Sprintf("%s", prefix)
	log.Info("Starting pinger for ", group.Name, " with metric path: ", metricPath)
	if (prefix != "") {prefix = prefix + "."};
	p.RunLoop()

pingloop:
	for {
		select {
		case res := <-onRecv:
			if _, ok := results[res.addr.String()]; ok {
				results[res.addr.String()] = res
			}
		case <-onIdle:
			for host, r := range results {
				outputLabel := index[host]
				if r == nil {
					log.Debug(fmt.Sprintf("%s : unreachable", outputLabel))
					// send a metric for a failed ping
					err := statsdClient.Inc(fmt.Sprintf("%s%s.failed", prefix, outputLabel), 1, 1)
					if err != nil {
						log.Error(fmt.Sprintf("Error sending metric: %+v", err))
					}
				} else {
					log.Debug(fmt.Sprintf("%s : %v", outputLabel, r.rtt))
					// send a zeroed failed metric, because we succeeded!
					err := statsdClient.Inc(fmt.Sprintf("%s%s.failed", prefix, outputLabel), 0, 1)
					if err != nil {
						log.Error(fmt.Sprintf("Error sending metric: %+v", err))
					}
					err = statsdClient.TimingDuration(fmt.Sprintf("%s%s.timer", prefix, outputLabel), r.rtt, 1)
				}
				results[host] = nil
			}
		case <-p.Done():
			if err := p.Err(); err != nil {
				return fmt.Errorf("Can't start pinger: %s", err)
			}
			break pingloop
		}
	}
	p.Stop()

	return errors.New("failed")

}
