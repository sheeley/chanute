package chanute

import (
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"github.com/richardwilkes/toolbox/errs"
)

type EBSReport struct {
	Volumes    []*EBSVolume
	Aggregated []*EBSAggregate
	Errors     []string
}

func (r *EBSReport) AsciiReport() string {
	if len(r.Volumes) == 0 {
		return "EBS: No issues"
	}

	o := &strings.Builder{}
	o.WriteString("EBS\n")

	w := tablewriter.NewWriter(o)
	w.SetHeader([]string{"Name", "ID", "Size (in GB)", "Monthly Cost"})

	if r.Aggregated == nil {
		for _, v := range r.Volumes {
			w.Append([]string{v.Name, v.ID, strconv.Itoa(v.Size), PrintDollars(v.MonthlyStorageCost)})
		}
		w.Render()
		return o.String()
	}

	for _, agg := range r.Aggregated {
		w.Append([]string{agg.Key, "", strconv.Itoa(agg.Size), PrintDollars(agg.MonthlyStorageCost)})

		if len(agg.Volumes) > 0 {
			for _, v := range agg.Volumes {
				w.Append([]string{v.Name, v.ID, strconv.Itoa(v.Size), PrintDollars(v.MonthlyStorageCost)})
			}
			w.Append([]string{"", "", "", ""})
		}
	}

	w.Render()
	return o.String()
}

type EBSAggregate struct {
	Key                string
	Volumes            []*EBSVolume
	Size               int
	MonthlyStorageCost int
}

type EBSVolume struct {
	ID                 string
	Name               string
	Type               string
	Region             string
	MonthlyStorageCost int
	Size               int

	SnapshotID   string
	SnapshotName string
	SnapshotAge  string

	Tags map[string]string
}

func ebsLowUtilization(config *Config, sess *session.Session, checks []*TrustedAdvisorCheck) (*EBSReport, error) {
	m := checksToMaps(checks)

	r := &EBSReport{}
	volumes := make(map[string]*EBSVolume, len(checks))

	for _, volume := range m {
		volumeID := volume["Volume ID"]
		volumes[volumeID] = &EBSVolume{
			ID:                 volumeID,
			Name:               volume["Volume Name"],
			Type:               volume["Volume Type"],
			Size:               parseAmount(volume["Volume Size"]),
			MonthlyStorageCost: parseAmount(volume["Monthly Storage Cost"]),
			Region:             volume["Region"],

			SnapshotID:   volume["Snapshot ID"],
			SnapshotName: volume["Snapshot Name"],
			SnapshotAge:  volume["Snapshot Age"],
		}
	}

	if config.GetTags {
		c := ec2.New(sess)

		var ids []*string
		for _, v := range r.Volumes {
			ids = append(ids, aws.String(v.ID))
		}

		input := &ec2.DescribeVolumesInput{
			VolumeIds: ids,
		}

		for {
			page, err := c.DescribeVolumes(input)
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

				var newIDs []*string
				for _, iID := range input.VolumeIds {
					if !nonExisting[aws.StringValue(iID)] {
						newIDs = append(newIDs, iID)
					}
				}

				input.VolumeIds = newIDs
				continue
			}

			for _, v := range page.Volumes {
				if v2, ok := volumes[aws.StringValue(v.VolumeId)]; ok {
					v2.Tags = ec2TagsToMap(v.Tags)
				}
			}

			if page.NextToken == nil {
				break
			}
			input.NextToken = page.NextToken
		}
	}

	for _, v := range volumes {
		r.Volumes = append(r.Volumes, v)
	}

	sort.Slice(r.Volumes, func(i, j int) bool {
		return r.Volumes[i].MonthlyStorageCost > r.Volumes[j].MonthlyStorageCost
	})

	if config.Aggregator != nil {
		aggregated := map[string]*EBSAggregate{}
		for _, v := range r.Volumes {
			key := config.Aggregator(v.Tags)
			if key == "" {
				key = v.Name
				if key == "" {
					key = v.ID
				}
			}
			if _, ok := aggregated[key]; !ok {
				aggregated[key] = &EBSAggregate{Key: key}
			}
			if !config.HideResourceDetails {
				aggregated[key].Volumes = append(aggregated[key].Volumes, v)
			}
			aggregated[key].Size += v.Size
			aggregated[key].MonthlyStorageCost += v.MonthlyStorageCost
		}

		for _, agg := range aggregated {
			r.Aggregated = append(r.Aggregated, agg)
		}

		sort.Slice(r.Aggregated, func(i, j int) bool {
			return r.Aggregated[i].MonthlyStorageCost > r.Aggregated[j].MonthlyStorageCost
		})
	}

	return r, nil
}
