package controller

import (
	"encoding/xml"
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/presort"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"time"
)

type Controller struct {
	ToolId  int
	Seconds int
	Output  output_connection.OutputConnection
	Info    recordinfo.RecordInfo
}

type config struct {
	Seconds int
}

const eventField = `Event`

func (plugin *Controller) Init(toolId int, configXml string) bool {
	plugin.ToolId = toolId
	var config config
	err := xml.Unmarshal([]byte(configXml), &config)
	if err != nil || config.Seconds == 0 {
		api.OutputMessage(toolId, api.Error, `invalid configuration`)
		return false
	}
	plugin.Seconds = config.Seconds

	generator := recordinfo.NewGenerator()
	generator.AddV_StringField(eventField, `Controller`, 10)
	plugin.Info = generator.GenerateRecordInfo()
	plugin.Output = output_connection.New(toolId, `Output`)
	return true
}

func (plugin *Controller) PushAllRecords(_ int) bool {
	err := plugin.Output.Init(plugin.Info)
	if err != nil {
		api.OutputMessage(plugin.ToolId, api.Error, err.Error())
		return false
	}

	api.OutputToolProgress(plugin.ToolId, 0)
	plugin.Output.UpdateProgress(0)

	if api.GetInitVar(plugin.ToolId, api.UpdateOnly) == `True` {
		return true
	}

	err = plugin.Info.SetStringField(eventField, `Start`)
	if err != nil {
		api.OutputMessage(plugin.ToolId, api.Error, err.Error())
		return false
	}
	record, err := plugin.Info.GenerateRecord()
	if err != nil {
		api.OutputMessage(plugin.ToolId, api.Error, err.Error())
		return false
	}

	plugin.Output.PushRecord(record)
	time.Sleep(time.Duration(plugin.Seconds) * time.Second)

	err = plugin.Info.SetStringField(eventField, "End")
	if err != nil {
		api.OutputMessage(plugin.ToolId, api.Error, err.Error())
		return false
	}
	record, err = plugin.Info.GenerateRecord()
	if err != nil {
		api.OutputMessage(plugin.ToolId, api.Error, err.Error())
		return false
	}

	plugin.Output.PushRecord(record)
	api.OutputMessage(plugin.ToolId, api.Complete, ``)
	plugin.Output.UpdateProgress(1)
	plugin.Output.Close()
	return true
}

func (plugin *Controller) Close(_ bool) {
	api.OutputMessage(plugin.ToolId, api.Info, `was asked to close`)
}

func (plugin *Controller) AddIncomingConnection(_ string, _ string) (api.IncomingInterface, *presort.PresortInfo) {
	return nil, nil
}

func (plugin *Controller) AddOutgoingConnection(_ string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output.Add(connectionInterface)
	return true
}

func (plugin *Controller) GetToolId() int {
	return plugin.ToolId
}
