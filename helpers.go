package chanute

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

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
	fmt.Printf("%s: %s is an unhandled check\n", ct, c)
	m := checksToMaps(v)
	if len(m) == 0 {
		return
	}
	structName := strings.ReplaceAll(string(c), " ", "")
	var typeBody []string
	var mapBody []string
	for k := range m[0] {
		cleaned := strings.ReplaceAll(k, " ", "")
		typeBody = append(typeBody, fmt.Sprintf("\t%s string", cleaned))
		mapBody = append(mapBody, fmt.Sprintf(`%s: lim["%s"],`, cleaned, k))
	}

	r, n := utf8.DecodeRuneInString(structName)
	funcName := string(unicode.ToLower(r)) + structName[n:]

	_ = newTemplate.Execute(os.Stdout, map[string]string{
		"StructName": structName,
		"StructBody": strings.Join(typeBody, "\n"),
		"MapBody":    strings.Join(mapBody, "\n"),
		"FuncName":   funcName,
	})
}

var newTemplate = template.Must(template.New("").Parse(`
type {{.StructName}}Report struct {
	Resources []*{{.StructName}}
	// Aggregated []*RedshiftAggregate
}

func (r *{{.StructName}}Report) AsciiReport() string {
	return ""
}

func (r *{{.StructName}}Report) AggegateRows(a Aggregator) []*AggregateRow {
	var o []*AggregateRow
	for _, res := range r.Resources {
		o = append(o, res.AggregateRow(a))
	}
	return o
}

type {{.StructName}} struct {
	{{.StructBody}}
}

func (r *{{.StructName}}) AggregateRow(a Aggregator) *AggregateRow {
	return &AggregateRow{
		Service: "{{.StructName}}",
	}
}

func {{.FuncName}}(cfg *Config, sess *session.Session, checks []*TrustedAdvisorCheck) (*{{.StructName}}Report, error) {
	r := &{{.StructName}}Report{}
	m := checksToMaps(checks)

	resources := make(map[string]*{{.StructName}}, len(checks))
	for _, lim := range m {
		resource := &{{.StructName}}{
			{{.MapBody}}
		}
		resources[resource.Name] = resource
	}
	return r, nil
}
`))
