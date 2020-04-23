package chanute

import (
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/olekukonko/tablewriter"
	"github.com/richardwilkes/toolbox/errs"
)

type RedshiftReport struct {
	Clusters   []*RedShiftCluster
	Aggregated []*RedshiftAggregate
}

type RedshiftAggregate struct {
	Key                     string
	EstimatedMonthlySavings int
	Clusters                []*RedShiftCluster
}

type RedShiftCluster struct {
	Type                    string
	Reason                  string
	EstimatedMonthlySavings int
	Status                  string
	Region                  string
	Name                    string
	Tags                    map[string]string
}

func (r *RedshiftReport) AsciiReport() string {
	if len(r.Clusters) == 0 {
		return "Redshift: No issues"
	}
	o := &strings.Builder{}
	o.WriteString("Redshift\n")

	w := tablewriter.NewWriter(o)
	w.SetHeader([]string{"Name", "Status", "Reason", "Monthly Cost"})

	if r.Aggregated == nil {
		for _, c := range r.Clusters {
			w.Append([]string{c.Name, c.Status, c.Reason, PrintDollars(c.EstimatedMonthlySavings)})
		}
		w.Render()
		return o.String()
	}

	for _, agg := range r.Aggregated {
		w.Append([]string{agg.Key, "", "", PrintDollars(agg.EstimatedMonthlySavings)})

		if len(agg.Clusters) > 0 {
			for _, c := range agg.Clusters {
				w.Append([]string{c.Name, c.Status, c.Reason, PrintDollars(c.EstimatedMonthlySavings)})
			}
			w.Append([]string{"", "", "", ""})
		}
	}

	w.Render()
	return o.String()
}

func redshiftLowUtilization(config *Config, sess *session.Session, checks []*TrustedAdvisorCheck) (*RedshiftReport, error) {
	m := checksToMaps(checks)

	clusters := make(map[string]*RedShiftCluster, len(checks))
	for _, instance := range m {
		c := &RedShiftCluster{
			Type:                    instance["Instance Type"],
			Reason:                  instance["Reason"],
			EstimatedMonthlySavings: parseAmount(instance["Estimated Monthly Savings"]),
			Status:                  instance["Status"],
			Region:                  instance["Region"],
			Name:                    instance["Cluster"],
		}
		clusters[c.Name] = c
	}

	if config.GetTags {
		allTags, err := GetRedshiftTags(sess)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		for id, tags := range allTags {
			if c2, ok := clusters[id]; ok {
				c2.Tags = tags
			}
		}
	}
	r := &RedshiftReport{}
	for _, c := range clusters {
		r.Clusters = append(r.Clusters, c)
	}

	sort.Slice(r.Clusters, func(i, j int) bool {
		return r.Clusters[i].EstimatedMonthlySavings > r.Clusters[j].EstimatedMonthlySavings
	})

	if config.Aggregator != nil {
		aggregated := map[string]*RedshiftAggregate{}
		for _, c := range r.Clusters {
			key := config.Aggregator(c.Tags)
			if key == "" {
				key = c.Name
			}
			if _, ok := aggregated[key]; !ok {
				aggregated[key] = &RedshiftAggregate{Key: key}
			}
			if !config.HideResourceDetails {
				aggregated[key].Clusters = append(aggregated[key].Clusters, c)
			}

			aggregated[key].EstimatedMonthlySavings += c.EstimatedMonthlySavings
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

func GetRedshiftTags(sess *session.Session) (TagMap, error) {
	c := redshift.New(sess)

	tags := map[string]map[string]string{}

	// var ids []*string
	// for _, v := range r.Volumes {
	// 	ids = append(ids, aws.String(v.ID))
	// }

	input := &redshift.DescribeClustersInput{
		// ClusterIdentifier: *string,
	}

	for {
		page, err := c.DescribeClusters(input)
		if err != nil {
			errStr := err.Error()

			if !strings.HasPrefix(errStr, "InvalidVolume.NotFound") {
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

			// TODO
			// var newIDs []*string
			// for _, iID := range input.VolumeIds {
			// 	if !nonExisting[aws.StringValue(iID)] {
			// 		newIDs = append(newIDs, iID)
			// 	}
			// }

			// input.VolumeIds = newIDs
			continue
		}

		for _, c := range page.Clusters {
			cid := aws.StringValue(c.ClusterIdentifier)
			tags[cid] = map[string]string{}
			for _, t := range c.Tags {
				tags[cid][aws.StringValue(t.Key)] = aws.StringValue(t.Value)
			}
		}

		if page.Marker == nil {
			break
		}
		input.Marker = page.Marker
	}

	return tags, nil
}

func (r *RedshiftReport) AggregateRows(a Aggregator) []*AggregateRow {
	var o []*AggregateRow
	for _, i := range r.Clusters {
		o = append(o, i.AggregateRow(a))
	}
	return o
}
func (r *RedShiftCluster) AggregateRow(a Aggregator) *AggregateRow {
	return &AggregateRow{
		Env:            "",
		Service:        "EC2",
		Key:            a(r.Tags),
		MonthlySavings: r.EstimatedMonthlySavings,
	}
}
