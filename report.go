package chanute

import (
	"fmt"
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
	Config *Config

	CostOptimization *CostReport
	ServiceLimits    *LimitReport
}

func (r *Report) AsciiReport() string {
	o := &strings.Builder{}

	if r.CostOptimization != nil {
		o.WriteString(r.CostOptimization.AsciiReport())
		o.WriteString("\n")
	}
	if r.ServiceLimits != nil {
		o.WriteString(r.ServiceLimits.AsciiReport())
		o.WriteString("\n")
	}

	return o.String()
}

func GenerateReport(sess *session.Session, options ...Option) (*Report, error) {
	return generateReport(sess, configFromOptions(options...))
}

func configFromOptions(options ...Option) *Config {
	cfg := &Config{}
	for _, o := range options {
		o(cfg)
	}
	if len(cfg.Checks) == 0 {
		WithCostOptimizationChecks()(cfg)
	}
	return cfg
}

func generateReport(sess *session.Session, cfg *Config) (*Report, error) {

	activeChecks := make(map[Check]bool, len(cfg.Checks))
	for _, c := range cfg.Checks {
		activeChecks[c] = true
	}

	checks, err := ListNonOKTrustedAdvisorChecks(sess, activeChecks)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	var lookups = map[CheckType]map[Check][]*TrustedAdvisorCheck{}
	for _, check := range checks {
		chk := Check(check.Name)
		if !activeChecks[chk] {
			fmt.Printf("skipping %s\n", check.Name)
			continue
		}
		if chkType, ok := checkTypeLookup[chk]; ok {
			if _, ok = lookups[chkType]; !ok {
				lookups[chkType] = map[Check][]*TrustedAdvisorCheck{}
			}
			lookups[chkType][chk] = append(lookups[chkType][chk], check)
			continue
		}
		fmt.Printf("%s not supported\n", check.Name)
	}

	r := &Report{
		Config: cfg,
	}

	var reportErr error

	for chk, values := range lookups {
		switch chk {
		case CheckTypeCost:
			r.CostOptimization, reportErr = costReport(cfg, sess, values)
		case CheckTypeServiceLimit:
			r.ServiceLimits, reportErr = serviceLimits(cfg, sess, values)
		case CheckTypeFaultTolerance:
			// r.FaultTolerance, err = faultTolerance(cfg, sess, values)
		default:
			printUnhandled(chk, values)
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
func ListNonOKTrustedAdvisorChecks(sess *session.Session, activeChecks map[Check]bool) ([]*TrustedAdvisorCheck, error) {
	c := support.New(sess)
	o, err := c.DescribeTrustedAdvisorChecks(&support.DescribeTrustedAdvisorChecksInput{Language: aws.String("en")})
	if err != nil {
		return nil, err
	}

	var results []*TrustedAdvisorCheck
	for _, ch := range o.Checks {
		if len(activeChecks) > 0 && !activeChecks[Check(aws.StringValue(ch.Name))] {
			continue
		}

		cho, err2 := c.DescribeTrustedAdvisorCheckResult(&support.DescribeTrustedAdvisorCheckResultInput{CheckId: ch.Id})
		if err2 != nil {
			err = errs.Append(err, err2)
			continue
		}

		if cho.Result == nil {
			err = errs.Append(err, errs.New("result or resources summary nil"))
			continue
		}

		var flagged int64
		var processed int64
		if cho.Result.ResourcesSummary != nil {
			flagged = aws.Int64Value(cho.Result.ResourcesSummary.ResourcesFlagged)
			processed = aws.Int64Value(cho.Result.ResourcesSummary.ResourcesProcessed)
		}

		if aws.StringValue(cho.Result.Status) == "ok" {
			continue
		}

		results = append(results, &TrustedAdvisorCheck{
			Name:        aws.StringValue(ch.Name),
			ID:          aws.StringValue(ch.Id),
			Status:      aws.StringValue(cho.Result.Status),
			Flagged:     flagged,
			Processed:   processed,
			Description: aws.StringValue(ch.Description),

			Check:  ch,
			Result: cho.Result,
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
