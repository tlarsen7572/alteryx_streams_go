package combine_latest

import (
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"unsafe"
)

type Ii struct {
	Parent        *Plugin
	ToolId        int
	initCallback  func(info recordinfo.RecordInfo)
	pushCallback  func(unsafe.Pointer) bool
	closeCallback func()
}

func (ii *Ii) Init(recordInfoIn string) bool {
	info, err := recordinfo.FromXml(recordInfoIn)
	if err != nil {
		api.OutputMessage(ii.ToolId, api.Error, err.Error())
	}
	ii.initCallback(info)
	return true
}

func (ii *Ii) PushRecord(record unsafe.Pointer) bool {
	return ii.pushCallback(record)
}

func (ii *Ii) UpdateProgress(percent float64) {
	ii.Parent.updateProgress(percent)
}

func (ii *Ii) Close() {
	ii.closeCallback()
}

func (ii *Ii) CacheSize() int {
	return 0
}
