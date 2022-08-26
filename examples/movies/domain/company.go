package domain

import (
	"github.com/mindstand/gogm/v2"
	"time"
)

type Company struct {
	// provides required node fields
	gogm.BaseUUIDNode

	Name        string    `gogm:"name=name"`
	Address     string    `gogm:"name=address"`
	Contacts    string    `gogm:"name=contacts"`
	Telephone   string    `gogm:"name=telephone"`
	CompanyType string    `gogm:"name=companyType"`
	Remark      string    `gogm:"name=remark"`
	CreateTime  time.Time `gogm:"name=createTime"`
	UpdateTime  time.Time `gogm:"name=updateTime"`
}
