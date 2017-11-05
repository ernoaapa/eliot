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
		Image:	{{.Image}}
    {{- if $status }}
		ContainerID:	{{$status.ContainerID}}
		State:	{{$status.State}}
    {{- end}}
    {{- if .Pipe}}
		Pipe:
			stdout -> stdin: {{.Pipe.Stdout.Stdin.Name}}
		{{- end}}
	{{end}}
`
