package guest

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

type Guest struct {
	Connection string
	Domain     string
	ByUUID     bool
	QGATimeout int
	conn       libvirt.Connect
	domain     libvirt.Domain
	domainName string
}

func (guest *Guest) Connect() error {
	conn, err := libvirt.NewConnect(fmt.Sprintf("qemu+tcp://%s/system", guest.Connection))
	if err != nil {
		return err
	}
	guest.conn = *conn
	var (
		domain *libvirt.Domain
	)
	domain, err = conn.LookupDomainByUUIDString(guest.Domain)
	if err != nil {
		domain, err = conn.LookupDomainByName(guest.Domain)
	}
	if err != nil {
		return err
	}
	guest.domain = *domain
	guest.domainName, _ = domain.GetName()
	return nil
}

func (g Guest) IsSame(other Guest) bool {
	return g.Connection == other.Connection && g.Domain == other.Domain
}
func (g *Guest) GetName() string {
	if g.domainName == "" {
		g.domainName, _ = g.domain.GetName()
	}
	return g.domainName
}

func (g Guest) String() string {
	return fmt.Sprintf("<%s %s>", g.Connection, g.Domain)
}
