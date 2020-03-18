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
		return "Load Balancers: No issues"
	}

	o := &strings.Builder{}
	o.WriteString("Load Balancers\n")

	w := tablewriter.NewWriter(o)
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

		input := &elbv2.DescribeLoadBalancersInput{}
		fauxMarker := aws.String("marker")
		var arns []*string
		for {
			// only change the names if we aren't paginating
			if input.Marker == nil {
				input.Names = names
				if len(names) > 20 {
					input.Names = names[0:20]
					names = names[20:]
				}
			}
			if input.Marker == fauxMarker {
				input.Marker = nil
			}

			if len(input.Names) == 0 {
				break
			}

			page, err := c.DescribeLoadBalancers(input)
			if err != nil {
				errStr := err.Error()

				if !strings.HasPrefix(errStr, "LoadBalancerNotFound") {
					return nil, errs.Wrap(err)
				}

				// if instances are not found, pull them out of the input
				start := strings.Index(errStr, "'[")
				end := strings.LastIndex(errStr, "]'")
				if start == -1 || end == -1 || start == end {
					return nil, errs.New("couldn't find two ' chars in error message")
				}

				idsStr := errStr[start+2 : end]
				idsToRemove := strings.Split(idsStr, ", ")
				nonExisting := make(map[string]bool, len(idsToRemove))
				for _, ec2ID := range idsToRemove {
					nonExisting[ec2ID] = true
				}

				var newNames []*string
				for _, lbName := range input.Names {
					if !nonExisting[aws.StringValue(lbName)] {
						newNames = append(newNames, lbName)
					}
				}

				input.Names = newNames
				input.Marker = fauxMarker
				continue
			}

			for _, lb := range page.LoadBalancers {
				if _, ok := lbs[aws.StringValue(lb.LoadBalancerName)]; ok {
					arns = append(arns, lb.LoadBalancerArn)
				}
			}
			input.Marker = page.NextMarker

			if err != nil {
				return nil, errs.Wrap(err)
			}
		}

		if len(arns) > 0 {
			input := &elbv2.DescribeTagsInput{ResourceArns: arns}

			// for {
			page, err := c.DescribeTags(input)
			if err != nil {
				return nil, errs.Wrap(err)
			}

			for _, res := range page.TagDescriptions {
				if lb, ok := lbs[aws.StringValue(res.ResourceArn)]; ok {
					lb.Tags = make(map[string]string, len(res.Tags))
					for _, t := range res.Tags {
						lb.Tags[aws.StringValue(t.Key)] = aws.StringValue(t.Value)
					}
				}
			}

			// 	break
			// }
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
