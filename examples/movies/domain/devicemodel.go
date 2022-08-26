package domain

import "github.com/mindstand/gogm/v2"

type DeviceModel struct {
	// provides required node fields
	gogm.BaseUUIDNode

	DeviceHeight string `gogm:"name=deviceHeight"`
	LineNumber   string `gogm:"name=lineNumber"`
	ModelName    string `gogm:"name=modelName"`
	NodeNumber   string `gogm:"name=nodeNumber"`
	Remark       string `gogm:"name=remark"`
	ModelType    string `gogm:"name=modelType"`

	Device []*Device `gogm:"direction=incoming;relationship=DEVICE_DEVICEMODEL"`
}
