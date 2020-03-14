package chanute

import (
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/olekukonko/tablewriter"
	"github.com/richardwilkes/toolbox/errs"
)

type LoadBalancerReport struct {
	LoadBalancers []*LoadBalancer
	Aggregated    []*LoadBalancerAggregate
}

type LoadBalancerAggregate struct {
	Key                     string
	LoadBalancers           []*LoadBalancer
	EstimatedMonthlySavings int
}

type LoadBalancer struct {
	Region                  string
	Name                    string
	Reason                  string
	EstimatedMonthlySavings int

	Tags map[string]string
}

func (r *LoadBalancerReport) AsciiReport() string {
	if len(r.LoadBalancers) == 0 {
		return "EBS: No issues"
	}

	o := &strings.Builder{}
	o.WriteString("EBS\n")

	w := tablewriter.NewWriter(o)
	w.SetAutoMergeCells(true)
	w.SetHeader([]string{"Name", "Region", "Reason", "Monthly Cost"})

	if r.Aggregated == nil {
		for _, lb := range r.LoadBalancers {
			w.Append([]string{lb.Name, lb.Region, lb.Reason, PrintDollars(lb.EstimatedMonthlySavings)})
		}
		w.Render()
		return o.String()
	}

	for _, agg := range r.Aggregated {
		w.Append([]string{agg.Key, "", "", PrintDollars(agg.EstimatedMonthlySavings)})

		if len(agg.LoadBalancers) > 0 {
			for _, lb := range agg.LoadBalancers {
				w.Append([]string{lb.Name, lb.Region, lb.Reason, PrintDollars(lb.EstimatedMonthlySavings)})
			}
			w.Append([]string{"", "", "", ""})
		}
	}

	w.Render()
	return o.String()
}

func idleLoadBalancers(config *Config, sess *session.Session, checks []*TrustedAdvisorCheck) (*LoadBalancerReport, error) {
	m := checksToMaps(checks)

	lbs := make(map[string]*LoadBalancer, len(checks))
	var names []*string
	for _, instance := range m {
		lb := &LoadBalancer{
			Region:                  instance["Region"],
			Name:                    instance["Load Balancer Name"],
			Reason:                  instance["Reason"],
			EstimatedMonthlySavings: parseAmount(instance["Estimated Monthly Savings"]),
		}
		names = append(names, aws.String(lb.Name))
		lbs[lb.Name] = lb
	}

	if config.GetTags {
		c := elbv2.New(sess)

		var arns []*string
		err := c.DescribeLoadBalancersPages(&elbv2.DescribeLoadBalancersInput{Names: names}, func(o *elbv2.DescribeLoadBalancersOutput, last bool) bool {
			for _, lb := range o.LoadBalancers {
				arns = append(arns, lb.LoadBalancerArn)
			}
			return true
		})
		if err != nil {
			return nil, errs.Wrap(err)
		}
		if len(arns) == 0 {
			return nil, nil
		}

		input := &elbv2.DescribeTagsInput{ResourceArns: arns}

		for {
			page, err := c.DescribeTags(input)
			if err != nil {
				errStr := err.Error()

				if !strings.HasPrefix(errStr, "InvalidInstanceID.NotFound") {
					return nil, errs.Wrap(err)
				}

				// if instances are not found, pull them out of the input
				start := strings.Index(errStr, "'")
				end := strings.LastIndex(errStr, "'")
				if start == -1 || end == -1 || start == end {
					return nil, errs.New("couldn't find two ' chars in error message")
				}

				idsStr := errStr[start+1 : end]
				idsToRemove := strings.Split(idsStr, ", ")
				nonExisting := make(map[string]bool, len(idsToRemove))
				for _, ec2ID := range idsToRemove {
					nonExisting[ec2ID] = true
				}

				var newIDs []*string
				for _, iID := range input.ResourceArns {
					if !nonExisting[aws.StringValue(iID)] {
						newIDs = append(newIDs, iID)
					}
				}

				input.ResourceArns = newIDs
				continue
			}

			for _, res := range page.TagDescriptions {
				if lb, ok := lbs[aws.StringValue(res.ResourceArn)]; ok {
					lb.Tags = make(map[string]string, len(res.Tags))
					for _, t := range res.Tags {
						lb.Tags[aws.StringValue(t.Key)] = aws.StringValue(t.Value)
					}
				}
			}
			break
		}

	}

	r := &LoadBalancerReport{}
	for _, lb := range lbs {
		r.LoadBalancers = append(r.LoadBalancers, lb)
	}

	sort.Slice(r.LoadBalancers, func(i, j int) bool {
		return r.LoadBalancers[i].EstimatedMonthlySavings > r.LoadBalancers[j].EstimatedMonthlySavings
	})

	if config.Aggregator != nil {
		aggregated := map[string]*LoadBalancerAggregate{}
		for _, i := range r.LoadBalancers {
			key := config.Aggregator(i.Tags)
			if key == "" {
				key = i.Name
			}
			if _, ok := aggregated[key]; !ok {
				aggregated[key] = &LoadBalancerAggregate{
					Key: key,
				}
			}
			if !config.HideResourceDetails {
				aggregated[key].LoadBalancers = append(aggregated[key].LoadBalancers, i)
			}
			aggregated[key].EstimatedMonthlySavings += i.EstimatedMonthlySavings
		}
		for _, agg := range aggregated {
			r.Aggregated = append(r.Aggregated, agg)
		}

		sort.Slice(r.Aggregated, func(i, j int) bool {
			return r.Aggregated[i].EstimatedMonthlySavings > r.Aggregated[j].EstimatedMonthlySavings
		})
	}

	return r, nil
}
