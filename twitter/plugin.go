package twitter

import (
	"encoding/xml"
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/presort"
)

type Plugin struct {
	ToolId int
	Output output_connection.OutputConnection
	Config *Config
}

type Config struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	Follow            string
	Track             string
}

func (plugin *Plugin) Init(toolId int, configXml string) bool {
	plugin.ToolId = toolId
	config := &Config{}
	err := xml.Unmarshal([]byte(configXml), config)
	if err != nil {
		api.OutputMessage(toolId, api.Error, err.Error())
		return false
	}
	plugin.Config = config
	plugin.Output = output_connection.New(toolId, `Output`)
	return true
}

func (plugin *Plugin) PushAllRecords(recordLimit int) bool {
	return false
}

func (plugin *Plugin) Close(hasErrors bool) {}

func (plugin *Plugin) AddIncomingConnection(connectionType string, connectionName string) (api.IncomingInterface, *presort.PresortInfo) {
	return &Ii{
		ToolId: plugin.ToolId,
		Config: plugin.Config,
		Output: plugin.Output,
	}, nil
}

func (plugin *Plugin) AddOutgoingConnection(connectionName string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output.Add(connectionInterface)
	return true
}

func (plugin *Plugin) GetToolId() int {
	return plugin.ToolId
}
