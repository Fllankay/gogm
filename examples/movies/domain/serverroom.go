package domain

import (
	"github.com/mindstand/gogm/v2"
	"time"
)

type ServerRoom struct {
	// provides required node fields
	gogm.BaseUUIDNode

	Address           string    `gogm:"name=address"`
	Area              string    `gogm:"name=area"`
	CommissioningDate time.Time `gogm:"name=commissioningDate"`
	Maintainer        string    `gogm:"name=maintainer"`
	Name              string    `gogm:"name=name"`
	Number            string    `gogm:"name=number"`
	Remark            string    `gogm:"name=remark"`
	Supplier          string    `gogm:"name=supplier"`
	Telephone         string    `gogm:"name=telephone"`

	CloudProvider *CloudProvider   `gogm:"direction=outgoing;relationship=SERVERROOM_CLOUDPROVIDER"`
	ServerCabinet []*ServerCabinet `gogm:"direction=incoming;relationship=SERVERCABINET_SERVERROOM"`
}

type SERVERROOM_CLOUDPROVIDER struct {
	// provides required node fields
	gogm.BaseUUIDNode

	Start *ServerRoom
	End   *CloudProvider
}
