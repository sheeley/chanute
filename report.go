package chanute

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/support"
	"github.com/richardwilkes/toolbox/errs"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer = message.NewPrinter(language.English)

func PrintDollars(i int) string {
	if i == 0 {
		return "0"
	}

	return printer.Sprintf("$%d", i)
}

type Report struct {
	config        *Config
	EC2           *EC2Report
	LoadBalancers *LoadBalancerReport
	EBS           *EBSReport
	RDS           *RDSReport
	Redshift      *RedshiftReport
}

func (r *Report) AsciiReport() string {
	o := &strings.Builder{}

	if r.EC2 != nil {
		o.WriteString(r.EC2.AsciiReport())
		o.WriteString("\n")
	}
	if r.LoadBalancers != nil {
		o.WriteString(r.LoadBalancers.AsciiReport())
		o.WriteString("\n")
	}
	if r.EBS != nil {
		o.WriteString(r.EBS.AsciiReport())
		o.WriteString("\n")
	}
	if r.RDS != nil {
		o.WriteString(r.RDS.AsciiReport())
		o.WriteString("\n")
	}
	if r.Redshift != nil {
		o.WriteString(r.Redshift.AsciiReport())
	}

	return o.String()
}

type Config struct {
	GetTags             bool
	HideResourceDetails bool
	Aggregator          Aggregator
	Checks              []Check
}

type Check string

const (
	CheckEC2           = "Low Utilization Amazon EC2 Instances"
	CheckLoadBalancers = "Idle Load Balancers"
	CheckEBS           = "Underutilized Amazon EBS Volumes"
	CheckRDS           = "Amazon RDS Idle DB Instances"
	CheckRedshift      = "Underutilized Amazon Redshift Clusters"
	// "Amazon EC2 Reserved Instances Optimization":    true,
	// "Unassociated Elastic IP Addresses":             true,
	// "Amazon Route 53 Latency Resource Record Sets":  true,
	// "Amazon EC2 Reserved Instance Lease Expiration": true,
)

type Aggregator func(map[string]string) string

type Option func(*Config)

func WithCustomTagAggregator(a Aggregator) Option {
	return func(c *Config) {
		c.GetTags = true
		c.Aggregator = a
	}
}

func WithoutResourceDetails() Option {
	return func(c *Config) {
		c.HideResourceDetails = true
	}
}

func WithAggregationByTag(t string) Option {
	return WithCustomTagAggregator(func(tags map[string]string) string {
		return tags[t]
	})
}

func WithChecks(checks ...Check) Option {
	return func(c *Config) {
		c.Checks = checks
	}
}

var defaultChecks = map[Check]bool{
	CheckEC2:           true,
	CheckLoadBalancers: true,
	CheckEBS:           true,
	CheckRDS:           true,
	CheckRedshift:      true,
}

func GenerateReport(sess *session.Session, options ...Option) (*Report, error) {
	cfg := &Config{}
	for _, o := range options {
		o(cfg)
	}

	checks, err := ListNonOKTrustedAdvisorChecks(sess)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	activeChecks := defaultChecks
	if len(cfg.Checks) > 0 {
		activeChecks = make(map[Check]bool, len(cfg.Checks))
		for _, c := range cfg.Checks {
			activeChecks[c] = true
		}
	}

	var lookups = map[string][]*TrustedAdvisorCheck{}
	for _, check := range checks {
		if !activeChecks[Check(check.Name)] {
			continue
		}
		lookups[check.Name] = append(lookups[check.Name], check)
	}

	r := &Report{
		config: cfg,
	}

	var reportErr error
	for lookup, values := range lookups {
		switch lookup {
		case "Low Utilization Amazon EC2 Instances":
			r.EC2, reportErr = ec2LowUtilization(cfg, sess, values)
		case "Idle Load Balancers":
			r.LoadBalancers, reportErr = idleLoadBalancers(cfg, sess, values)
		case "Underutilized Amazon EBS Volumes":
			r.EBS, reportErr = ebsLowUtilization(cfg, sess, values)
		case "Amazon RDS Idle DB Instances":
			r.RDS, reportErr = rdsIdleInstances(cfg, sess, values)
		case "Underutilized Amazon Redshift Clusters":
			r.Redshift, reportErr = redshiftLowUtilization(cfg, sess, values)
		}
		if reportErr != nil {
			err = errs.Append(err, reportErr)
		}
	}

	return r, err
}

type TrustedAdvisorCheck struct {
	Name               string
	ID                 string
	Status             string
	Description        string
	Flagged, Processed int64

	// Check is used to get the high-level description of a check
	Check *support.TrustedAdvisorCheckDescription
	// Result is used to get detailed information about which resources are failing a check
	Result *support.TrustedAdvisorCheckResult
}

// ListNonOKTrustedAdvisorChecks queries Trusted Advisor and only returns checks that have a status of error or warning
// These are typically worth review, and opening a ticket to increase limits.
func ListNonOKTrustedAdvisorChecks(sess *session.Session) ([]*TrustedAdvisorCheck, error) {
	c := support.New(sess)
	o, err := c.DescribeTrustedAdvisorChecks(&support.DescribeTrustedAdvisorChecksInput{Language: aws.String("en")})
	if err != nil {
		return nil, err
	}

	var results []*TrustedAdvisorCheck
	for _, ch := range o.Checks {
		o, err2 := c.DescribeTrustedAdvisorCheckResult(&support.DescribeTrustedAdvisorCheckResultInput{CheckId: ch.Id})
		if err2 != nil {
			err = errs.Append(err, err2)
			continue
		}

		if o.Result == nil {
			err = errs.Append(err, errs.New("result or resources summary nil"))
			continue
		}

		var flagged int64
		var processed int64
		if o.Result.ResourcesSummary != nil {
			flagged = aws.Int64Value(o.Result.ResourcesSummary.ResourcesFlagged)
			processed = aws.Int64Value(o.Result.ResourcesSummary.ResourcesProcessed)
		}

		if aws.StringValue(o.Result.Status) == "ok" {
			continue
		}

		results = append(results, &TrustedAdvisorCheck{
			Name:        aws.StringValue(ch.Name),
			ID:          aws.StringValue(ch.Id),
			Status:      aws.StringValue(o.Result.Status),
			Flagged:     flagged,
			Processed:   processed,
			Description: aws.StringValue(ch.Description),

			Check:  ch,
			Result: o.Result,
		})
	}
	return results, err
}

func checksToMaps(checks []*TrustedAdvisorCheck) []map[string]string {
	var o []map[string]string
	for _, check := range checks {
		for _, res := range check.Result.FlaggedResources {
			mapped := map[string]string{}
			for idx, md := range check.Check.Metadata {
				if len(res.Metadata) > idx {
					mapped[aws.StringValue(md)] = aws.StringValue(res.Metadata[idx])
				}
			}
			o = append(o, mapped)
		}
	}
	return o
}
