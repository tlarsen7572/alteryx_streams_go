package combine_latest

import (
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"unsafe"
)

type Ii struct {
	Parent *Plugin
	Name   string
	ToolId int
}

func (ii *Ii) Init(recordInfoIn string) bool {
	info, err := recordinfo.FromXml(recordInfoIn)
	if err != nil {
		api.OutputMessage(ii.ToolId, api.Error, err.Error())
	}
	ii.Parent.initInput(ii.Name, info)
	return true
}

func (ii *Ii) PushRecord(record unsafe.Pointer) bool {
	panic("implement me")
}

func (ii *Ii) UpdateProgress(percent float64) {
	panic("implement me")
}

func (ii *Ii) Close() {
	panic("implement me")
}

func (ii *Ii) CacheSize() int {
	return 0
}
