package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

type ZoneInput struct {
	IP              string
	QName           string
	DefaultZoneFile string
}

type Zone struct {
	ID              string
	IP              string
	QName           string
	ZoneFileHost    string
	ZoneFileTarget  string
	Dir             string
	Logs            string
	Config          string
	DefaultZoneFile string
	Implementation  *implementation
}

func (c *Config) newZone(build string, input *ZoneInput) (*Zone, error) {
	var err error
	id := generateZoneID(input.QName)
	dir := filepath.Join(build, id)
	impl, err := bind("9.18.4")
	return &Zone{
		ID:              id,
		IP:              input.IP,
		QName:           input.QName,
		ZoneFileHost:    filepath.Join(build, "zones", fmt.Sprintf("%s.zone", id)),
		ZoneFileTarget:  filepath.Join("/zones", fmt.Sprintf("%s.zone", id)),
		Dir:             dir,
		Logs:            filepath.Join(dir, "logs"),
		Config:          filepath.Join(dir, "config"),
		DefaultZoneFile: input.DefaultZoneFile,
		Implementation:  impl,
	}, err
}

func generateZoneID(qname string) string {
	if qname == "." {
		return "root"
	}
	f := func(c rune) bool {
		return c == '.'
	}
	split := strings.FieldsFunc(qname, f)
	return strings.Join(split, "-")
}
