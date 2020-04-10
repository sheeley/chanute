package chanute

type Config struct {
	GetTags             bool
	HideResourceDetails bool
	Aggregator          Aggregator
	Checks              []Check
}

type Aggregator func(map[string]string) string

type Option func(*Config)

func WithCustomTagAggregator(a Aggregator) Option {
	return func(c *Config) {
		c.GetTags = true
		c.Aggregator = a
	}
}

func WithoutResourceDetails() Option {
	return func(c *Config) {
		c.HideResourceDetails = true
	}
}

func WithAggregationByTag(t string) Option {
	return WithCustomTagAggregator(func(tags map[string]string) string {
		return tags[t]
	})
}

func WithChecks(checks ...Check) Option {
	return func(c *Config) {
		c.Checks = checks
	}
}

func WithCostOptimizationChecks() Option {
	return func(c *Config) {
		c.Checks = append(c.Checks, costChecks...)
	}
}

func WithPerformanceChecks() Option {
	return func(c *Config) {
		c.Checks = append(c.Checks, performanceChecks...)
	}
}

func WithSecurityChecks() Option {
	return func(c *Config) {
		c.Checks = append(c.Checks, securityChecks...)
	}
}

func WithFaultToleranceChecks() Option {
	return func(c *Config) {
		c.Checks = append(c.Checks, faultToleranceChecks...)
	}
}

func WithServiceLimitChecks() Option {
	return func(c *Config) {
		c.Checks = append(c.Checks, serviceLimitChecks...)
	}
}

type TagMap map[string]map[string]string

func stringPtrSet(ss []*string) map[*string]bool {
	m := make(map[*string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}
