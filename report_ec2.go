package chanute

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"github.com/richardwilkes/toolbox/errs"
)

type EC2Report struct {
	Instances  []*EC2Instance
	Aggregated []*EC2Aggregate
	Errors     []string
}

func (r *EC2Report) AsciiReport() string {
	if len(r.Instances) == 0 {
		return "EC2: No issues"
	}

	o := &strings.Builder{}
	o.WriteString("EC2\n")

	w := tablewriter.NewWriter(o)
	w.SetAutoMergeCells(true)
	w.SetHeader([]string{"Name", "ID", "Low Utilization Days", "Estimated Monthly Savings"})

	if r.Aggregated == nil {
		for _, i := range r.Instances {
			w.Append([]string{i.Name, i.ID, strconv.Itoa(i.LowUtilizationDays), strconv.Itoa(i.EstimatedMonthlySavings)})
		}
		w.Render()
		return o.String()
	}

	for _, agg := range r.Aggregated {
		w.Append([]string{agg.Key, "", "", PrintDollars(agg.EstimatedMonthlySavings)})

		if len(agg.Instances) > 0 {
			for _, i := range agg.Instances {
				w.Append([]string{i.Name, i.ID, strconv.Itoa(i.LowUtilizationDays), strconv.Itoa(i.EstimatedMonthlySavings)})
			}
			w.Append([]string{"", "", "", ""})
		}
	}

	w.Render()
	return o.String()
}

type EC2Aggregate struct {
	Key                     string
	Instances               []*EC2Instance
	EstimatedMonthlySavings int
}

type EC2Instance struct {
	Name     string
	ID       string
	Type     string
	RegionAZ string

	EstimatedMonthlySavings int

	LowUtilizationDays  int
	Network14DayAverage string
	CPU14DayAverage     string

	Day1, Day2, Day3, Day4, Day5, Day6, Day7, Day8, Day9, Day10, Day11, Day12, Day13, Day14 string

	Tags map[string]string
}

func ec2LowUtilization(config *Config, sess *session.Session, checks []*TrustedAdvisorCheck) (*EC2Report, error) {
	r := &EC2Report{}

	m := checksToMaps(checks)

	instances := make(map[string]*EC2Instance, len(checks))
	var ids []*string
	for _, instance := range m {

		instanceID := instance["Instance ID"]
		ids = append(ids, aws.String(instanceID))

		instances[instanceID] = &EC2Instance{
			Name:     instance["Instance Name"],
			ID:       instanceID,
			Type:     instance["Instance Type"],
			RegionAZ: instance["Region/AZ"],

			EstimatedMonthlySavings: parseAmount(instance["Estimated Monthly Savings"]),

			LowUtilizationDays: parseDays(instance["Number of Days Low Utilization"]),

			Network14DayAverage: instance["14-Day Average Network I/O"],
			CPU14DayAverage:     instance["14-Day Average CPU Utilization"],

			Day1:  instance["Day 1"],
			Day2:  instance["Day 2"],
			Day3:  instance["Day 3"],
			Day4:  instance["Day 4"],
			Day5:  instance["Day 5"],
			Day6:  instance["Day 6"],
			Day7:  instance["Day 7"],
			Day8:  instance["Day 8"],
			Day9:  instance["Day 9"],
			Day10: instance["Day 10"],
			Day11: instance["Day 11"],
			Day12: instance["Day 12"],
			Day13: instance["Day 13"],
			Day14: instance["Day 14"],
		}
	}

	if config.GetTags {
		c := ec2.New(sess)
		input := &ec2.DescribeInstancesInput{
			InstanceIds: ids,
		}

		for {
			page, err := c.DescribeInstances(input)
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
				for _, iID := range input.InstanceIds {
					if !nonExisting[aws.StringValue(iID)] {
						newIDs = append(newIDs, iID)
					}
				}

				input.InstanceIds = newIDs
				continue
			}

			for _, res := range page.Reservations {
				for _, i := range res.Instances {
					if ei, ok := instances[aws.StringValue(i.InstanceId)]; ok {
						ei.Tags = ec2TagsToMap(i.Tags)
					}
				}
			}

			if page.NextToken == nil {
				break
			}
			input.NextToken = page.NextToken
		}

	}

	for _, i := range instances {
		r.Instances = append(r.Instances, i)
	}

	sort.Slice(r.Instances, func(i, j int) bool {
		return r.Instances[i].EstimatedMonthlySavings > r.Instances[j].EstimatedMonthlySavings
	})

	if config.Aggregator != nil {
		aggregated := map[string]*EC2Aggregate{}
		for _, i := range r.Instances {
			key := config.Aggregator(i.Tags)
			if key == "" {
				key = i.Name
				if key == "" {
					key = i.ID
				}

			}
			if _, ok := aggregated[key]; !ok {
				aggregated[key] = &EC2Aggregate{
					Key: key,
				}
			}
			if !config.HideResourceDetails {
				aggregated[key].Instances = append(aggregated[key].Instances, i)
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

func ec2TagsToMap(t []*ec2.Tag) map[string]string {
	o := map[string]string{}
	for _, tag := range t {
		if _, ok := o[*tag.Key]; ok {
			fmt.Println("already a value!")
		}
		o[*tag.Key] = *tag.Value
	}
	return o
}
