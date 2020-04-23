package chanute

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/olekukonko/tablewriter"
)

type UnassociatedElasticIPAddressesReport struct {
	IPs []*UnassociatedElasticIPAddresses
	// Aggregated []*RedshiftAggregate
}

func (r *UnassociatedElasticIPAddressesReport) AsciiReport() string {
	if len(r.IPs) == 0 {
		return "EIPs: No issues"
	}
	o := &strings.Builder{}
	o.WriteString("EIPs\n")

	w := tablewriter.NewWriter(o)
	w.SetHeader([]string{"IP", "Region"})
	for _, ip := range r.IPs {
		w.Append([]string{ip.IPAddress, ip.Region})
	}
	w.Render()
	return o.String()
}

func (r *UnassociatedElasticIPAddressesReport) AggregateRows(a Aggregator) []*AggregateRow {
	var o []*AggregateRow
	for _, ip := range r.IPs {
		o = append(o, ip.AggregateRow(a))
	}
	return o
}

type UnassociatedElasticIPAddresses struct {
	Region    string
	IPAddress string
}

func (r *UnassociatedElasticIPAddresses) AggregateRow(a Aggregator) *AggregateRow {
	return &AggregateRow{
		Service:        "EIP",
		MonthlySavings: 7,
	}
}

func unassociatedElasticIPAddresses(cfg *Config, sess *session.Session, checks []*TrustedAdvisorCheck) (*UnassociatedElasticIPAddressesReport, error) {
	r := &UnassociatedElasticIPAddressesReport{}
	m := checksToMaps(checks)

	resources := make(map[string]*UnassociatedElasticIPAddresses, len(checks))
	for _, lim := range m {
		resource := &UnassociatedElasticIPAddresses{
			Region:    lim["Region"],
			IPAddress: lim["IP Address"],
		}
		resources[resource.IPAddress] = resource
	}
	return r, nil
}
