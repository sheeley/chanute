package chanute

import (
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/olekukonko/tablewriter"
)

type LimitReport struct {
	Limits []*ServiceLimit
}

func (r *LimitReport) Title() string {
	return "Service Limits"
}

func (r *LimitReport) Headers(includeEnv ...bool) []string {
	var o []string
	if len(includeEnv) == 1 && includeEnv[0] {
		o = append(o, "Account")
	}
	o = append(o, "Status", "Service", "Limit Name", "Region", "Limit Amount", "Current Usage")
	return o
}

func (r *LimitReport) Rows(env ...string) [][]string {
	if len(r.Limits) == 0 {
		return nil
	}
	var o [][]string
	for _, l := range r.Limits {
		var row []string
		if len(env) == 1 && env[0] != "" {
			row = append(row, env[0])
		}
		row = append(row,
			l.Status,
			l.Service,
			l.LimitName,
			l.Region,
			strconv.Itoa(l.LimitAmount),
			strconv.Itoa(l.CurrentUsage),
		)
		o = append(o, row)
	}
	return o
}

func (r *LimitReport) AsciiReport() string {
	if len(r.Limits) == 0 {
		return "Service Limits: No issues"
	}
	o := &strings.Builder{}
	o.WriteString("Service Limits\n")

	w := tablewriter.NewWriter(o)

	w.SetHeader([]string{"Status", "Service", "Limit Name", "Region", "Limit Amount", "Current Usage"})
	for _, l := range r.Limits {
		w.Append([]string{
			l.Status,
			l.Service,
			l.LimitName,
			l.Region,
			strconv.Itoa(l.LimitAmount),
			strconv.Itoa(l.CurrentUsage),
		})
	}
	w.Render()
	return o.String()
}

type ServiceLimit struct {
	Service, Region, Status, LimitName string
	LimitAmount, CurrentUsage          int
}

func serviceLimits(config *Config, sess *session.Session, lookups map[Check][]*TrustedAdvisorCheck) (*LimitReport, error) {
	r := &LimitReport{}
	for _, checks := range lookups {
		m := checksToMaps(checks)

		for _, lim := range m {
			sts := lim["Status"]
			if sts == "Green" {
				continue
			}

			r.Limits = append(r.Limits, &ServiceLimit{
				Service:      lim["Service"],
				LimitName:    lim["Limit Name"],
				LimitAmount:  parseAmount(lim["Limit Amount"]),
				CurrentUsage: parseAmount(lim["Current Usage"]),
				Status:       sts,
				Region:       lim["Region"],
			})
		}
	}
	return r, nil
}
