package collector

import (
	"context"
	"errors"
	"os/exec"
	"regexp"
	"time"

	"github.com/prometheus-community/windows_exporter/log"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("tcp_ports", newtcp_portsCollector)
}

type tcp_portsCollector struct {
	ListenCount *prometheus.Desc
}

func newtcp_portsCollector() (Collector, error) {
	const subsystem = "tcp_ports"
	return &tcp_portsCollector{
		ListenCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "listen_count"),
			"Amount of listening tcp ports",
			nil,
			nil,
		),
	}, nil
}

func exec_netstat(proto string, timeout int) (string, error) {
	// setup a new context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	// create the command with the context
	cmd := exec.CommandContext(ctx, "netstat", "-ano", "-p", proto)

	out, err := cmd.Output()

	//check if the execution timed out
	if ctx.Err() == context.DeadlineExceeded {
		return "", errors.New("netstat command timed out")
	}

	//check execution result
	if err != nil {
		log.Debugf("netstat output: %s", out)
		return "", err
	}

	return string(out), nil
}

func parse_netstat_out(output string) (int, error) {

	re := regexp.MustCompile(`:(?P<port>\d+)\s+.*?LISTENING\s+(?P<pid>\d+)`)

	matches := re.FindAllStringSubmatch(output, -1)
	if matches == nil {
		return -1, errors.New("Error parsing netstat output")
	}

	return len(matches), nil
}

func (c *tcp_portsCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {

	const timeout = 3

	log.Debug("running netstat to determine listening IPv4 ports")
	netstat_v4_out, err := exec_netstat("TCP", timeout)
	if err != nil {
		log.Errorf("error running netstat: %s", err)
		return err
	}

	ipv4_listen_count, err := parse_netstat_out(netstat_v4_out)

	if err != nil {
		log.Errorf("error parsing netstat: %s", err)
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.ListenCount, prometheus.GaugeValue, float64(ipv4_listen_count))

	return nil

}
