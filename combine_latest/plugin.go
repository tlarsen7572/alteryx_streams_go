package combine_latest

import (
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/presort"
	"github.com/tlarsen7572/goalteryx/recordblob"
	"github.com/tlarsen7572/goalteryx/recordcopier"
	"github.com/tlarsen7572/goalteryx/recordinfo"
)

type Plugin struct {
	ToolId      int
	Output      output_connection.OutputConnection
	LeftIn      recordinfo.RecordInfo
	RightIn     recordinfo.RecordInfo
	Out         recordinfo.RecordInfo
	LeftCopier  *recordcopier.RecordCopier
	RightCopier *recordcopier.RecordCopier
	hasLeft     bool
	hasRight    bool
	progress    float64
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

func (plugin *Plugin) AddIncomingConnection(connectionType string, _ string) (api.IncomingInterface, *presort.PresortInfo) {
	if connectionType == `Left` {
		return &Ii{
			Parent:        plugin,
			ToolId:        plugin.ToolId,
			initCallback:  plugin.initLeft,
			pushCallback:  plugin.pushLeft,
			closeCallback: plugin.closeLeft,
		}, nil
	}
	return &Ii{
		Parent:        plugin,
		ToolId:        plugin.ToolId,
		initCallback:  plugin.initRight,
		pushCallback:  plugin.pushRight,
		closeCallback: plugin.closeRight,
	}, nil
}

func (plugin *Plugin) AddOutgoingConnection(_ string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output.Add(connectionInterface)
	return true
}

func (plugin *Plugin) GetToolId() int {
	return plugin.ToolId
}

func (plugin *Plugin) initLeft(info recordinfo.RecordInfo) {
	plugin.LeftIn = info
	if plugin.LeftIn != nil && plugin.RightIn != nil {
		plugin.initOutput()
	}
}

func (plugin *Plugin) initRight(info recordinfo.RecordInfo) {
	plugin.RightIn = info
	if plugin.LeftIn != nil && plugin.RightIn != nil {
		plugin.initOutput()
	}
}

func (plugin *Plugin) initOutput() {
	plugin.setUpOutInfo()
	err := plugin.Output.Init(plugin.Out)
	if err != nil {
		api.OutputMessage(plugin.ToolId, api.Error, err.Error())
	}
	plugin.setUpCopiers()
}

func (plugin *Plugin) pushLeft(record *recordblob.RecordBlob) bool {
	err := plugin.LeftCopier.Copy(record)
	if err != nil {
		api.OutputMessage(plugin.ToolId, api.Error, err.Error())
		return false
	}
	plugin.hasLeft = true
	return plugin.tryPushOut(record)
}

func (plugin *Plugin) pushRight(record *recordblob.RecordBlob) bool {
	err := plugin.RightCopier.Copy(record)
	if err != nil {
		api.OutputMessage(plugin.ToolId, api.Error, err.Error())
		return false
	}
	plugin.hasRight = true
	return plugin.tryPushOut(record)
}

func (plugin *Plugin) tryPushOut(record *recordblob.RecordBlob) bool {
	if !(plugin.hasLeft && plugin.hasRight) {
		return true
	}
	record, err := plugin.Out.GenerateRecord()
	if err != nil {
		api.OutputMessage(plugin.ToolId, api.Error, err.Error())
	}
	plugin.Output.PushRecord(record)
	return true
}

func (plugin *Plugin) closeLeft() {
	plugin.hasLeft = false
	plugin.tryFinalClose()
}

func (plugin *Plugin) closeRight() {
	plugin.hasRight = false
	plugin.tryFinalClose()
}

func (plugin *Plugin) tryFinalClose() {
	if plugin.hasLeft || plugin.hasRight {
		return
	}
	api.OutputMessage(plugin.ToolId, api.Complete, ``)
	plugin.Output.Close()
}

func (plugin *Plugin) setUpOutInfo() {
	generator := recordinfo.NewGenerator()
	for _, input := range []recordinfo.RecordInfo{plugin.LeftIn, plugin.RightIn} {
		for index := 0; index < input.NumFields(); index++ {
			field, _ := input.GetFieldByIndex(index)
			generator.AddField(field, `Combine Latest`)
		}
	}
	plugin.Out = generator.GenerateRecordInfo()
}

func (plugin *Plugin) setUpCopiers() {
	leftIndexMaps := []recordcopier.IndexMap{}
	outIndex := 0
	for index := 0; index < plugin.LeftIn.NumFields(); index++ {
		leftIndexMaps = append(leftIndexMaps, recordcopier.IndexMap{
			DestinationIndex: outIndex,
			SourceIndex:      index,
		})
		outIndex++
	}
	plugin.LeftCopier, _ = recordcopier.New(plugin.Out, plugin.LeftIn, leftIndexMaps)

	rightIndexMaps := []recordcopier.IndexMap{}
	for index := 0; index < plugin.RightIn.NumFields(); index++ {
		rightIndexMaps = append(rightIndexMaps, recordcopier.IndexMap{
			DestinationIndex: outIndex,
			SourceIndex:      index,
		})
		outIndex++
	}
	plugin.RightCopier, _ = recordcopier.New(plugin.Out, plugin.RightIn, rightIndexMaps)
}

func (plugin *Plugin) updateProgress(percent float64) {
	if percent > plugin.progress {
		plugin.progress = percent
	}
	api.OutputToolProgress(plugin.ToolId, plugin.progress)
	plugin.Output.UpdateProgress(plugin.progress)
}
