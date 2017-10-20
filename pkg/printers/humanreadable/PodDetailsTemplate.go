package humanreadable

// PodDetailsTemplate is go template for printing pod details
const PodDetailsTemplate = `Name:	{{.Pod.Metadata.Name}}
Namespace:	{{.Pod.Metadata.Namespace}}
Device: {{.Pod.Status.Hostname}}
State:	{{.Status}}
{{if .Pod.Status}}
Containers:{{range .Pod.Status.ContainerStatuses}}
  {{.ContainerID}}:
    Image: {{.Image}}
    State: {{.State}}
{{end}}
{{else}}
  (No container statuses available)
{{end}}
`
