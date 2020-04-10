package chanute

import (
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
}

func GenerateAggregateReport(envs []*Environment, options ...Option) (*AggregateReport, error) {
	var wg sync.WaitGroup
	wg.Add(len(envs))

	ar := &AggregateReport{
		LimitReports: map[string]*LimitReport{},
		// CostReports:  map[string]*CostReport{},
		Config: configFromOptions(options...),
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
				// if r.CostOptimization != nil {
				// 	ar.CostReports[e.Name] = r.CostOptimization
				// }
			}
			wg.Done()
		}(e)
	}
	wg.Wait()

	return ar, oErr
}
