package output

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// Output handles all output during an operation
type Output struct {
	outputers        []Outputter
	webhookPusher    *webhookPusher
	prometheusPusher *prometheusPusher
}

// Outputter bundles Webhook and Prometheus into one interface
type Outputter interface {
	Webhook
	Prometheus
	Console
}

// Trigger is the interface used for a single. Also intended for the use outside
// of the output module.
type Trigger interface {
	TriggerHook(data JsonMarshaller)
	TriggerProm(prom prometheus.Collector)
}

// Webhook defines an interface to send stats to a webhook url
type Webhook interface {
	// GetWebhookData should return a slice with all the objects that should be
	// sent to the webhook endpoint.
	GetWebhookData() []JsonMarshaller
}

// JsonMarshaller defines objects that can be marshalled to json.
type JsonMarshaller interface {
	ToJson() []byte
}

// Prometheus defines an interface to send stats to a prometheus push gateway
type Prometheus interface {
	// ToProm should return a slice of prometheus collectors that are ready to
	// be sent to the push gateway.
	ToProm() []prometheus.Collector
}

type Console interface {
	GetError() error
	GetStdOut() []string
	GetStdErrOut() []string
}

// New returns a new output object. It handles all the different outputs that
// should happen during a wrestic operation.
func New(webhookURL, promURL, hostname string) *Output {
	return &Output{
		outputers:        make([]Outputter, 0),
		webhookPusher:    newWebhookHandler(webhookURL),
		prometheusPusher: newPrometheusPusher(promURL, hostname),
	}
}

// Register adds an object to the slice of outputs to handle
func (o *Output) Register(out Outputter) {
	o.outputers = append(o.outputers, out)
}

// TriggerAll handles all the registered outputs and pushes the output to the
// webhook/prometheus endpoint.
func (o *Output) TriggerAll() {
	for _, out := range o.outputers {
		if out.GetError() != nil {
			fmt.Printf("Error occurred: %v\n command output:\n", out.GetError())
			fmt.Println(strings.Join(out.GetStdErrOut(), "\n"))
		}
		for _, prom := range out.ToProm() {
			if prom != nil {
				o.TriggerProm(prom)
			}
		}
		for _, hook := range out.GetWebhookData() {
			if hook != nil {
				fmt.Printf("Sending webhooks to %v: ", o.webhookPusher.url)
				o.TriggerHook(hook)
			}
		}
	}
}

// TriggerProm pushes a prometheus collector
func (o *Output) TriggerProm(prom prometheus.Collector) {
	o.prometheusPusher.Update(prom)
}

// TriggerHook pushes a single json
func (o *Output) TriggerHook(data JsonMarshaller) {
	err := o.webhookPusher.Push(data)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("done")
	}
}
