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
	if guest.ByUUID {
		domain, err = conn.LookupDomainByUUIDString(guest.Domain)
	} else {
		domain, err = conn.LookupDomainByName(guest.Domain)
	}

	if err != nil {
		return err
	}
	guest.domain = *domain
	return nil
}
