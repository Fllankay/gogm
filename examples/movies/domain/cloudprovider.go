package domain

import "github.com/mindstand/gogm/v2"

type CloudProvider struct {
	// provides required node fields
	gogm.BaseUUIDNode

	Name       string        `gogm:"name=name"`
	ServerRoom []*ServerRoom `gogm:"direction=incoming;relationship=SERVERROOM_CLOUDPROVIDER"`
}
