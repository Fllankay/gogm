package domain

import (
	"github.com/mindstand/gogm/v2"
	"time"
)

type Device struct {
	// provides required node fields
	gogm.BaseUUIDNode

	AssetBelong  string    `gogm:"name=assetBelong"`
	AssetNumber  string    `gogm:"name=assetNumber"`
	RegisterTime time.Time `gogm:"name=registerTime"`
	Remark       string    `gogm:"name=remark"`
	Sn           string    `gogm:"name=sn"`
	UseType      string    `gogm:"name=useType"`

	ServerCabinet *ServerCabinet `gogm:"direction=outgoing;relationship=DEVICE_SERVERCABINET"`
	DeviceModel   *DeviceModel   `gogm:"direction=outgoing;relationship=DEVICE_DEVICEMODEL"`
}

type DEVICE_SERVERCABINET struct {
	// provides required node fields
	gogm.BaseUUIDNode

	Start    *Device
	End      *ServerCabinet
	Location int64 `gogm:"name=location"`
}

type DEVICE_DEVICEMODEL struct {
	// provides required node fields
	gogm.BaseUUIDNode

	Start *Device
	End   *DeviceModel
}
