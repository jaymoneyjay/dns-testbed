package attack

import "dns-testbed-go/testbed/component"

type template struct {
}

func NewTemplateAttack() *template {
	return &template{}
}

func (t *template) WriteZoneFilesAndReturnEntryZone(param int, nameservers []*component.Nameserver) (string, error) {
	for _, ns := range nameservers {
		err := ns.UpdateLocal("template.zone")
		if err != nil {
			return "", err
		}
	}
	return "target.com", nil
}

func (t *template) Name() string {
	return "template"
}
