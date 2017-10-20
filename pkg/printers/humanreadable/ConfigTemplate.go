package humanreadable

// ConfigTemplate is template for printing out config in human readable format
var ConfigTemplate = `
Namespace: {{.Namespace }}
Endpoints:{{range .Endpoints}}
	Name: {{.Name }}
	URL: {{.URL }}
{{end}}
`
