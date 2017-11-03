package humanreadable

// PodDetailsTemplate is go template for printing pod details
const PodDetailsTemplate = `{{$pod := .Pod -}}
Name:	{{.Pod.Metadata.Name}}
Namespace:	{{.Pod.Metadata.Namespace}}
Device:	{{.Pod.Status.Hostname}}
State:	{{.Status}}
Containers:{{range .Pod.Spec.Containers}}
  {{- $status := GetStatus $pod .Name}}
  {{.Name}}:
    ContainerID: {{$status.ContainerID}}
    Image: {{$status.Image}}
    State: {{$status.State}}
    {{- if .Pipe}}
    Pipe:
      stdout -> stdin: {{.Pipe.Stdout.Stdin.Name}}
		{{- end}}
	{{end}}
`
