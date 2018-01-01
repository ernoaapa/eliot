package humanreadable

// DeviceDetailsTemplate is go template for printing device details
const DeviceDetailsTemplate = `Hostname:	{{.Hostname}}
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
`
