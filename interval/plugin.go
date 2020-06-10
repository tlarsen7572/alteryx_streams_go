package interval

import (
	"encoding/xml"
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/presort"
)

type Plugin struct {
	ToolId  int
	Seconds int
	Output  output_connection.OutputConnection
}

type config struct {
	Seconds int
}

func (plugin *Plugin) Init(toolId int, configXml string) bool {
	plugin.ToolId = toolId
	var config config
	err := xml.Unmarshal([]byte(configXml), &config)
	if err != nil || config.Seconds == 0 {
		api.OutputMessage(toolId, api.Error, `Configuration not valid for Interval`)
		return false
	}
	plugin.Seconds = config.Seconds
	plugin.Output = output_connection.New(toolId, `Output`)
	return true
}

func (plugin *Plugin) PushAllRecords(_ int) bool {
	return false
}

func (plugin *Plugin) Close(_ bool) {}

func (plugin *Plugin) AddIncomingConnection(_ string, _ string) (api.IncomingInterface, *presort.PresortInfo) {
	return &Ii{
		Output:       plugin.Output,
		ToolId:       plugin.ToolId,
		SleepSeconds: plugin.Seconds,
	}, nil
}

func (plugin *Plugin) AddOutgoingConnection(_ string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output.Add(connectionInterface)
	return true
}

func (plugin *Plugin) GetToolId() int {
	return plugin.ToolId
}
