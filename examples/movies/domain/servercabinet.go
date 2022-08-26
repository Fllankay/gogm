package domain

import (
	"github.com/mindstand/gogm/v2"
	"time"
)

type ServerCabinet struct {
	// provides required node fields
	gogm.BaseUUIDNode

	Capacity      string    `gogm:"name=capacity"`
	ElectricPower string    `gogm:"name=electricPower"`
	EndDate       time.Time `gogm:"name=endDate"`
	Maintainer    string    `gogm:"name=maintainer"`
	Name          string    `gogm:"name=name"`
	Number        string    `gogm:"name=number"`
	Remark        string    `gogm:"name=remark"`
	StartDate     time.Time `gogm:"name=startDate"`

	ServerRoom *ServerRoom `gogm:"direction=outgoing;relationship=SERVERCABINET_SERVERROOM"`
	Device     []*Device   `gogm:"direction=incoming;relationship=DEVICE_SERVERCABINET"`
}

type SERVERCABINET_SERVERROOM struct {
	// provides required node fields
	gogm.BaseUUIDNode

	Start *ServerCabinet
	End   *ServerRoom
}
