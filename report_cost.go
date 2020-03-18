package chanute

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/richardwilkes/toolbox/errs"
)

type CostReport struct {
	EC2           *EC2Report
	LoadBalancers *LoadBalancerReport
	EBS           *EBSReport
	RDS           *RDSReport
	Redshift      *RedshiftReport
}

func costReport(cfg *Config, sess *session.Session, lookups map[Check][]*TrustedAdvisorCheck) (*CostReport, error) {
	r := &CostReport{}
	var err error
	for lookup, values := range lookups {
		var reportErr error
		switch lookup {
		case CheckLowUtilizationAmazonEC2Instances:
			r.EC2, reportErr = ec2LowUtilization(cfg, sess, values)
		case CheckIdleLoadBalancers:
			r.LoadBalancers, reportErr = idleLoadBalancers(cfg, sess, values)
		case CheckUnderutilizedAmazonEBSVolumes:
			r.EBS, reportErr = ebsLowUtilization(cfg, sess, values)
		case CheckAmazonRDSIdleDBInstances:
			r.RDS, reportErr = rdsIdleInstances(cfg, sess, values)
		case CheckUnderutilizedAmazonRedshiftClusters:
			r.Redshift, reportErr = redshiftLowUtilization(cfg, sess, values)
		default:
			fmt.Println(lookup + " unhandled")
		}
		if reportErr != nil {
			err = errs.Append(err, reportErr)
		}
	}
	return r, err
}

func (r *CostReport) AsciiReport() string {
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
		o.WriteString("\n")
	}

	return o.String()
}
