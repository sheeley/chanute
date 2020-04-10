package chanute

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/olekukonko/tablewriter"
	"github.com/richardwilkes/toolbox/errs"
)

type RDSReport struct {
	Instances  []*RDSInstance
	Aggregated []*RDSAggregate
}

type RDSAggregate struct {
	Key                     string
	Instances               []*RDSInstance
	StorageProvisionedGB    int
	EstimatedMonthlySavings int
}

type RDSInstance struct {
	Region                  string
	Name                    string
	Type                    string
	MultiAZ                 bool
	StorageProvisionedGB    int
	DaysSinceLastConnection int
	EstimatedMonthlySavings int
	Tags                    map[string]string
}

func (r *RDSReport) AsciiReport() string {
	if len(r.Instances) == 0 {
		return "RDS: No issues"
	}

	o := &strings.Builder{}
	o.WriteString("RDS\n")

	w := tablewriter.NewWriter(o)
	w.SetHeader([]string{"Name", "MultiAZ", "Days Since Connection", "Storage Size (in GB)", "Monthly Cost"})

	if r.Aggregated == nil {
		for _, i := range r.Instances {
			w.Append([]string{i.Name, strconv.FormatBool(i.MultiAZ), strconv.Itoa(i.DaysSinceLastConnection), strconv.Itoa(i.StorageProvisionedGB), PrintDollars(i.EstimatedMonthlySavings)})
		}
		w.Render()
		return o.String()
	}

	for _, agg := range r.Aggregated {
		w.Append([]string{agg.Key, "", "", strconv.Itoa(agg.StorageProvisionedGB), PrintDollars(agg.EstimatedMonthlySavings)})

		if len(agg.Instances) > 0 {
			for _, i := range agg.Instances {
				w.Append([]string{i.Name, strconv.FormatBool(i.MultiAZ), strconv.Itoa(i.DaysSinceLastConnection), strconv.Itoa(i.StorageProvisionedGB), PrintDollars(i.EstimatedMonthlySavings)})
			}
			w.Append([]string{"", "", "", ""})
		}
	}

	w.Render()
	return o.String()
}

func rdsIdleInstances(config *Config, sess *session.Session, checks []*TrustedAdvisorCheck) (*RDSReport, error) {
	m := checksToMaps(checks)

	instances := make(map[string]*RDSInstance, len(checks))
	var names []*string
	for _, instance := range m {
		name := instance["DB Instance Name"]
		// names = append(names, aws.String(name))
		ri := &RDSInstance{
			Region:                  instance["Region"],
			Name:                    name,
			MultiAZ:                 instance["Multi-AZ"] == "Yes",
			Type:                    instance["Instance Type"],
			StorageProvisionedGB:    parseAmount(instance["Storage Provisioned (GB)"]),
			DaysSinceLastConnection: parseDays(instance["Days Since Last Connection"]),
			EstimatedMonthlySavings: parseAmount(instance["Estimated Monthly Savings (On Demand)"]),
		}
		instances[ri.Name] = ri
		names = append(names, &ri.Name)
	}

	if config.GetTags {
		tags, err := GetRDSTags(sess, names)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		for name, instanceTags := range tags {
			i, ok := instances[name]
			if !ok {
				continue
			}
			i.Tags = instanceTags
		}
	}

	r := &RDSReport{}

	for _, instance := range instances {
		r.Instances = append(r.Instances, instance)
	}
	sort.Slice(r.Instances, func(i, j int) bool {
		return r.Instances[i].EstimatedMonthlySavings > r.Instances[j].EstimatedMonthlySavings
	})

	if config.Aggregator != nil {
		aggregated := map[string]*RDSAggregate{}
		for _, i := range r.Instances {
			key := config.Aggregator(i.Tags)
			if key == "" {
				key = i.Name
			}
			if _, ok := aggregated[key]; !ok {
				aggregated[key] = &RDSAggregate{
					Key: key,
				}
			}
			if !config.HideResourceDetails {
				aggregated[key].Instances = append(aggregated[key].Instances, i)
			}
			aggregated[key].EstimatedMonthlySavings += i.EstimatedMonthlySavings
			aggregated[key].StorageProvisionedGB += i.StorageProvisionedGB
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

func GetRDSTags(sess *session.Session, names []*string) (TagMap, error) {
	c := rds.New(sess)

	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	tags := map[string]map[string]string{}
	for _, n := range names {
		arn := fmt.Sprintf("arn:aws:rds:%s:%s:db:%s", *sess.Config.Region, *result.Account, *n)
		resp, err := c.ListTagsForResource(&rds.ListTagsForResourceInput{
			ResourceName: &arn,
		})
		if err != nil {
			if strings.HasPrefix(err.Error(), "DBInstanceNotFound") {
				continue
			}
			fmt.Println(err)
			continue
		}

		tags[aws.StringValue(n)] = map[string]string{}
		for _, t := range resp.TagList {
			tags[aws.StringValue(n)][aws.StringValue(t.Key)] = aws.StringValue(t.Value)
		}
	}
	return tags, nil
}
