package chanute

import (
	"sort"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/richardwilkes/toolbox/errs"
)

type Environment struct {
	Name    string
	Session *session.Session
}

type AggregateReport struct {
	Config       *Config
	Reports      []*Report
	LimitReports map[string]*LimitReport
	CostReports  map[string]*CostReport
}

type AggregateSummary struct {
	allRows        []*AggregateRow
	aggregatedRows []*AggregateRow
	resourcesByKey map[string][]*AggregateRow
	TotalSavings   int
}

func (s *AggregateSummary) SummaryHeaders() []string {
	return []string{"Key", "Monthly Savings"}
}
func (s *AggregateSummary) SummaryRows() [][]string {
	var o [][]string
	for _, r := range s.aggregatedRows {
		o = append(o, []string{r.Key, strconv.Itoa(r.MonthlySavings)})
	}
	return o
}

func (s *AggregateSummary) Details() []*AggregateDetail {
	return nil
}

type AggregateRow struct {
	Service, Key, Env string
	MonthlySavings    int
}

type CostReporter interface {
	AggregateRow() *AggregateRow
	Detail() *AggregateDetail
}

type AggregateDetail struct {
	Key string
}

func (d *AggregateDetail) CSV() string {
	return "TODO"
}

func (r *AggregateReport) AggregatedCostSummary(a Aggregator) *AggregateSummary {
	sum := &AggregateSummary{}

	for _, report := range r.CostReports {
		sum.allRows = append(sum.allRows, report.AggregateRows(a)...)
	}

	aggTotal := map[string]*AggregateRow{}
	sum.resourcesByKey = map[string][]*AggregateRow{}
	for _, row := range sum.allRows {
		if _, ok := aggTotal[row.Key]; !ok {
			kr := &AggregateRow{Key: row.Key}
			sum.aggregatedRows = append(sum.aggregatedRows, kr)
			aggTotal[row.Key] = kr
		}
		sum.TotalSavings += row.MonthlySavings
		aggTotal[row.Key].MonthlySavings += row.MonthlySavings
		sum.resourcesByKey[row.Key] = append(sum.resourcesByKey[row.Key], row)
	}

	sort.Slice(sum.aggregatedRows, func(i, j int) bool {
		return sum.aggregatedRows[i].MonthlySavings > sum.aggregatedRows[j].MonthlySavings
	})

	return sum
}

type AggregateByKey struct {
	Key       string
	Total     int
	Resources []*string
}

func GenerateAggregateReport(envs []*Environment, options ...Option) (*AggregateReport, error) {
	var wg sync.WaitGroup
	wg.Add(len(envs))

	ar := &AggregateReport{
		LimitReports: map[string]*LimitReport{},
		CostReports:  map[string]*CostReport{},
		Config:       configFromOptions(options...),
	}

	if ar.Config.Aggregator == nil {
		return nil, errs.New("Aggregator is required")
	}

	var oErr error
	for _, e := range envs {
		go func(e *Environment) {
			r, err := generateReport(e.Session, ar.Config)
			if err != nil {
				oErr = errs.Append(oErr, err)
			} else {
				ar.Reports = append(ar.Reports, r)
				if r.ServiceLimits != nil {
					ar.LimitReports[e.Name] = r.ServiceLimits
				}
				if r.CostOptimization != nil {
					ar.CostReports[e.Name] = r.CostOptimization
				}
			}
			wg.Done()
		}(e)
	}
	wg.Wait()

	return ar, oErr
}
