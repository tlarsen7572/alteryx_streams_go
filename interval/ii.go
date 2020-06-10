package interval

import (
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"time"
	"unsafe"
)

type Ii struct {
	Output       output_connection.OutputConnection
	ToolId       int
	SleepSeconds int
	inInfo       recordinfo.RecordInfo
	outInfo      recordinfo.RecordInfo
	doLoop       bool
	counter      int
	loopDone     chan bool
}

func (ii *Ii) Init(recordInfoIn string) bool {
	ii.loopDone = make(chan bool)
	var err error
	ii.inInfo, err = recordinfo.FromXml(recordInfoIn)
	if err != nil {
		api.OutputMessage(ii.ToolId, api.Error, err.Error())
		return false
	}

	generator := recordinfo.NewGenerator()
	generator.AddInt64Field(`Count`, `Interval`)
	ii.outInfo = generator.GenerateRecordInfo()
	err = ii.Output.Init(ii.outInfo)
	if err != nil {
		api.OutputMessage(ii.ToolId, api.Error, err.Error())
		return false
	}
	return true
}

func (ii *Ii) PushRecord(record unsafe.Pointer) bool {
	event, isNull, err := ii.inInfo.GetStringValueFrom(`Event`, record)
	if err != nil {
		api.OutputMessage(ii.ToolId, api.Error, err.Error())
		return false
	}
	if isNull {
		return true
	}
	if event == `Start` {
		ii.handleStart()
		return true
	}
	if event == `End` {
		ii.handleEnd()
		return true
	}
	return true
}

func (ii *Ii) UpdateProgress(percent float64) {
	api.OutputToolProgress(ii.ToolId, percent)
	ii.Output.UpdateProgress(percent)
}

func (ii *Ii) Close() {
	<-ii.loopDone
	api.OutputMessage(ii.ToolId, api.Complete, ``)
	api.OutputToolProgress(ii.ToolId, 1)
	ii.Output.Close()
}

func (ii *Ii) CacheSize() int {
	return 0
}

func (ii *Ii) handleStart() {
	ii.doLoop = true
	ii.counter = 0
	go loop(ii)
}

func (ii *Ii) handleEnd() {
	ii.doLoop = false
}

func loop(ii *Ii) {
	for ii.doLoop {
		time.Sleep(time.Duration(ii.SleepSeconds) * time.Second)
		ii.counter++
		err := ii.outInfo.SetIntField(`Count`, ii.counter)
		if err != nil {
			api.OutputMessage(ii.ToolId, api.Error, err.Error())
			continue
		}
		record, err := ii.outInfo.GenerateRecord()
		if err != nil {
			api.OutputMessage(ii.ToolId, api.Error, err.Error())
			continue
		}
		ii.Output.PushRecord(record)
	}
	ii.loopDone <- true
}
