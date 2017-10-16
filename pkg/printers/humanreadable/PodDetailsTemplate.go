package humanreadable

// PodDetailsTemplate is go template for printing pod details
const PodDetailsTemplate = `Name: {{.Pod.Metadata.Name}}
Namespace:	{{.Pod.Metadata.Namespace}}
State: {{.PodStatus}}
Containers:{{range .Pod.Status.ContainerStatuses}}
  {{.ContainerID}}:
    Image: {{.Image}}
    State: {{.State}}
{{else}}
  (No container statuses available)
{{end}}
`
