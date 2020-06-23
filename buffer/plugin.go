package buffer

import (
	"encoding/xml"
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/presort"
)

type Plugin struct {
	ToolId  int
	Records int
	Output  output_connection.OutputConnection
}

type Config struct {
	Records int
}

func (plugin *Plugin) Init(toolId int, configXml string) bool {
	plugin.ToolId = toolId
	var config Config
	err := xml.Unmarshal([]byte(configXml), &config)
	if err != nil {
		api.OutputMessage(toolId, api.Error, err.Error())
		return false
	}
	plugin.Records = config.Records
	plugin.Output = output_connection.New(toolId, `Output`)
	return true
}

func (plugin *Plugin) PushAllRecords(_ int) bool {
	return false
}

func (plugin *Plugin) Close(_ bool) {}

func (plugin *Plugin) AddIncomingConnection(_ string, _ string) (api.IncomingInterface, *presort.PresortInfo) {
	return &Ii{
		ToolId:  plugin.ToolId,
		Output:  plugin.Output,
		Records: plugin.Records,
	}, nil
}

func (plugin *Plugin) AddOutgoingConnection(_ string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output.Add(connectionInterface)
	return true
}

func (plugin *Plugin) GetToolId() int {
	return plugin.ToolId
}
