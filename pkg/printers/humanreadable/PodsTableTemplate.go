package humanreadable

const PodsTableTemplate = `NAMESPACE	NAME	CONTAINERS{{range .}}
{{.Metadata.namespace}}	{{.Metadata.name}}	{{len .Spec.Containers}}
{{else}}
No pods
{{end}}
`
