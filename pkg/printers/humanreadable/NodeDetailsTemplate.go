package humanreadable

// NodeDetailsTemplate is go template for printing node details
const NodeDetailsTemplate = `Hostname:	{{.Hostname}}
Uptime:	{{FormatUptime .Uptime}}
Arch/OS:	{{.Os}}/{{.Arch}}
Version:	{{.Version}}
Labels:{{range .Labels}}
	{{.Key}}={{.Value}}
{{- end}}
Addresses:{{range .Addresses}}
	{{.}}
{{- end}}
GrpcPort:	{{.GrpcPort}}
MachineID:	{{.MachineID}}
SystemUUID:	{{.SystemUUID}}
BootID:	{{.BootID}}
{{- if .Filesystems }}
Filesystems:
	Filesystem	Type	Size	Used	Available	Use%	Mounted on
	----------	----	----	----	---------	----	----------
{{- range .Filesystems}}
	{{.Filesystem}}	{{.TypeName}}	{{FormatBytes .Total}}	{{Subtract .Total .Free | FormatBytes}}	{{FormatBytes .Available}}	{{FormatPercent .Total .Free .Available}}	{{.MountDir}}
{{- end}}
{{- end}}
`
