package output

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	namespace = "baas"
	subsystem = "backup_restic"
)

type prometheusPusher struct {
	url      string
	hostname string
}

func newPrometheusPusher(url, hostname string) *prometheusPusher {
	return &prometheusPusher{
		url:      url,
		hostname: hostname,
	}
}

func (p *prometheusPusher) Update(collector prometheus.Collector) {
	push.New(p.url, "restic_backup").Collector(collector).
		Grouping("instance", p.hostname).
		Add()
}
