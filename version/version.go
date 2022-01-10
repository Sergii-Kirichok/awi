package version

const (
	name    = "AWI-Current"
	svcName = "AWI-Service"
	version = "0.0.1 (20220106.01)"
)

type Info struct {
	Version string
	Name    string
	SvcName string
}

func GetInfo() Info {
	var data Info
	data.Name = name
	data.Version = version
	data.SvcName = svcName
	return data
}
