package combine_latest

import (
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/presort"
	"github.com/tlarsen7572/goalteryx/recordinfo"
)

type Plugin struct {
	ToolId  int
	Output  output_connection.OutputConnection
	LeftIn  recordinfo.RecordInfo
	RightIn recordinfo.RecordInfo
	Out     recordinfo.RecordInfo
}

func (plugin *Plugin) Init(toolId int, _ string) bool {
	plugin.ToolId = toolId
	plugin.Output = output_connection.New(toolId, `Output`)
	return true
}

func (plugin *Plugin) PushAllRecords(_ int) bool {
	return false
}

func (plugin *Plugin) Close(_ bool) {}

func (plugin *Plugin) AddIncomingConnection(_ string, connectionName string) (api.IncomingInterface, *presort.PresortInfo) {
	return &Ii{
		Parent: plugin,
		Name:   connectionName,
		ToolId: plugin.ToolId,
	}, nil
}

func (plugin *Plugin) AddOutgoingConnection(_ string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output.Add(connectionInterface)
	return true
}

func (plugin *Plugin) GetToolId() int {
	return plugin.ToolId
}

func (plugin *Plugin) initInput(name string, info recordinfo.RecordInfo) {
	if name == `Left` {
		plugin.LeftIn = info
	} else {
		plugin.RightIn = info
	}
	if plugin.LeftIn != nil && plugin.RightIn != nil {
		plugin.initOutput()
	}
}

func (plugin *Plugin) initOutput() {
	generator := recordinfo.NewGenerator()
	for _, input := range []recordinfo.RecordInfo{plugin.LeftIn, plugin.RightIn} {
		for index := 0; index < input.NumFields(); index++ {
			field, _ := input.GetFieldByIndex(index)
			generator.AddField(field, `Combine Latest`)
		}
	}
	plugin.Out = generator.GenerateRecordInfo()
	err := plugin.Output.Init(plugin.Out)
	if err != nil {
		api.OutputMessage(plugin.ToolId, api.Error, err.Error())
	}
}
