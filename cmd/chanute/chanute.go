package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sheeley/chanute"
)

func aggregator(tags map[string]string) string {
	for k, v := range tags {
		tags[strings.ToLower(k)] = strings.ToLower(v)
	}

	if team, ok := tags["team"]; ok {
		team = strings.TrimSpace(team)
		if team != "" {
			return team
		}
	}

	if o, ok := tags["organization"]; ok {
		spl := strings.Split(strings.ToLower(o), ":::")

		for _, s := range spl {
			kvSpl := strings.Split(s, "=")
			if len(kvSpl) != 2 {
				continue
			}
			if kvSpl[0] == "team" {
				return kvSpl[1]
			}
		}
	}

	return ""
	// return "No Team Tag"
}

func main() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
	r, err := chanute.GenerateReport(sess,
		chanute.WithCustomTagAggregator(aggregator),
		chanute.WithoutResourceDetails(),
		// chanute.WithChecks(
		// 	chanute.CheckEBS,
		// 	chanute.CheckEC2,
		// 	chanute.CheckRDS,
		// 	chanute.CheckLoadBalancers,
		// 	chanute.CheckRedshift,
		// ),
	)
	if err != nil {
		panic(err)
	}
	output := r.AsciiReport()
	fmt.Println(output)
}
