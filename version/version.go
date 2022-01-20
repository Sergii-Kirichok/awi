package version

import "fmt"

const (
	name    = "Avigilon Weight Integration"
	svcName = "AWI-Service"
	version = "1.0.13 (20220120.01)"
)

type Info struct {
	Version string
	Name    string
	SvcName string
}

func NewInfo() *Info {
	var data Info
	data.Name = name
	data.Version = version
	data.SvcName = svcName
	return &data
}

func (i *Info) String() string {
	return fmt.Sprintf("%s %s", i.Name, i.Version)
}
