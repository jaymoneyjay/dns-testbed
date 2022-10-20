package attack

import (
	"dns-testbed-go/testbed/component"
)

type Attack interface {
	WriteZoneFilesAndReturnEntryZone(param int, nameservers []*component.Nameserver) (string, error)
	Name() string
}
