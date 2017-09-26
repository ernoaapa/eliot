package humanreadable

// PodsTableHeader is header for the table output
const PodsTableHeader = `NAMESPACE	NAME	CONTAINERS`

// PodsTableRowTemplate is golang template for printing table of pods information
const PodsTableRowTemplate = `{{.Metadata.Namespace}}	{{.Metadata.Name}}	{{len .Spec.Containers}}
`
