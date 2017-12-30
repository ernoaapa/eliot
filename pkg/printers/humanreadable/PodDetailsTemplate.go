package humanreadable

// PodDetailsTemplate is go template for printing pod details
const PodDetailsTemplate = `{{$pod := .Pod -}}
Name:	{{.Pod.Metadata.Name}}
Namespace:	{{.Pod.Metadata.Namespace}}
Device:	{{.Pod.Status.Hostname}}
State:	{{.Status}}
Restart Policy:	{{.Pod.Spec.RestartPolicy}}
Host Network:	{{.Pod.Spec.HostNetwork}}
Host PID:	{{.Pod.Spec.HostPID}}
Containers:{{range .Pod.Spec.Containers}}
  {{- $status := GetStatus $pod .Name}}
	{{.Name}}:
		Image:	{{.Image}}
    {{- if $status }}
		ContainerID:	{{$status.ContainerID}}
		State:	{{$status.State}}
		Restart Count:	{{$status.RestartCount}}
		Working Dir:	{{.WorkingDir}}
		{{- end}}
		Args:{{range .Args}}
			- {{.}}
		{{- end}}
		Env:{{range .Env}}
			- {{.}}
		{{- end}}
		Mounts:{{range .Mounts}}
			- type={{.Type}},source={{.Source}},destination={{.Destination}},options={{StringsJoin .Options ":"}}
		{{- end}}
    {{- if .Pipe}}
		Pipe:
			stdout -> stdin: {{.Pipe.Stdout.Stdin.Name}}
		{{- end}}
	{{end}}
`
