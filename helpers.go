package chanute

import (
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func TableString(headers []string, rows [][]string) string {
	o := &strings.Builder{}
	Table(o, headers, rows)
	return o.String()
}

func Table(o io.Writer, headers []string, rows [][]string) {
	if len(rows) == 0 {
		return
	}
	w := tablewriter.NewWriter(o)
	w.SetHeader(headers)
	w.AppendBulk(rows)
	w.Render()
}

// printUnhandled is used to generate types for Trusted Advisor checks
func printUnhandled(ct CheckType, checks map[Check][]*TrustedAdvisorCheck) {
	fmt.Println(ct + " is an unhandled check type")
	for c, v := range checks {
		printUnhandledCheck(ct, c, v)
	}
}

func printUnhandledCheck(ct CheckType, c Check, v []*TrustedAdvisorCheck) {
	fmt.Println(c + " is an unhandled check")
	m := checksToMaps(v)
	if len(m) == 0 {
		return
	}
	structName := string(c)
	var typeBody []string
	var mapBody []string
	for k := range m[0] {
		cleaned := strings.ReplaceAll(k, " ", "")
		typeBody = append(typeBody, fmt.Sprintf("\t%s string", cleaned))
		mapBody = append(mapBody, fmt.Sprintf(`c.%s = lim["%s"]`, cleaned, k))
	}
	fmt.Printf(`
type %s struct {
%s
}

%s
`, strings.ReplaceAll(structName, " ", ""), strings.Join(typeBody, "\n"), strings.Join(mapBody, "\n"))
}
