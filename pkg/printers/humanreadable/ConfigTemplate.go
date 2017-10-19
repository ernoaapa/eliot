package humanreadable

// ConfigTemplate is template for printing out config in human readable format
var ConfigTemplate = `
Current context:
	Name: {{.GetCurrentContext.Name }}
	Namespace: {{.GetCurrentContext.Namespace }}
	Endpoint:
		Name: {{.GetCurrentEndpoint.Name }}
		URL: {{.GetCurrentEndpoint.URL }}
`
