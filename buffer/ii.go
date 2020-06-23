package buffer

import (
	"fmt"
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/recordblob"
	"github.com/tlarsen7572/goalteryx/recordcopier"
	"github.com/tlarsen7572/goalteryx/recordinfo"
)

type Ii struct {
	ToolId       int
	Output       output_connection.OutputConnection
	Records      int
	inInfo       recordinfo.RecordInfo
	copiers      []*recordcopier.RecordCopier
	outInfo      recordinfo.RecordInfo
	currentIndex int
}

func (ii *Ii) Init(recordInfoIn string) bool {
	var err error
	ii.inInfo, err = recordinfo.FromXml(recordInfoIn)
	if err != nil {
		api.OutputMessage(ii.ToolId, api.Error, err.Error())
		return false
	}

	generator := recordinfo.NewGenerator()
	for outputSet := 1; outputSet <= ii.Records; outputSet++ {
		for fieldIndex := 0; fieldIndex < ii.inInfo.NumFields(); fieldIndex++ {
			field, err := ii.inInfo.GetFieldByIndex(fieldIndex)
			if err != nil {
				api.OutputMessage(ii.ToolId, api.Error, err.Error())
				return false
			}
			generator.AddFieldUsingName(field, fmt.Sprintf(`%v %v`, field.Name, outputSet), `Buffer`)
		}
	}
	ii.outInfo = generator.GenerateRecordInfo()
	err = ii.Output.Init(ii.outInfo)
	if err != nil {
		api.OutputMessage(ii.ToolId, api.Error, err.Error())
		return false
	}

	outIndex := 0
	for outputSet := 1; outputSet <= ii.Records; outputSet++ {
		var indexMaps []recordcopier.IndexMap
		for fieldIndex := 0; fieldIndex < ii.inInfo.NumFields(); fieldIndex++ {
			indexMaps = append(indexMaps, recordcopier.IndexMap{
				DestinationIndex: outIndex,
				SourceIndex:      fieldIndex,
			})
			outIndex++
		}
		copier, err := recordcopier.New(ii.outInfo, ii.inInfo, indexMaps)
		if err != nil {
			api.OutputMessage(ii.ToolId, api.Error, err.Error())
			return false
		}
		ii.copiers = append(ii.copiers, copier)
	}

	return true
}

func (ii *Ii) PushRecord(record recordblob.RecordBlob) bool {
	err := ii.copiers[ii.currentIndex].Copy(record)
	if err != nil {
		return false
	}
	ii.currentIndex++

	if ii.currentIndex == ii.Records {
		record, err := ii.outInfo.GenerateRecord()
		if err != nil {
			api.OutputMessage(ii.ToolId, api.Error, err.Error())
			return false
		}
		ii.Output.PushRecord(record)
		ii.currentIndex = 0
	}
	return true
}

func (ii *Ii) UpdateProgress(percent float64) {
	api.OutputToolProgress(ii.ToolId, percent)
	ii.Output.UpdateProgress(percent)
}

func (ii *Ii) Close() {
	if ii.currentIndex > 0 {
		for fieldIndex := ii.currentIndex * ii.inInfo.NumFields(); fieldIndex < ii.outInfo.NumFields(); fieldIndex++ {
			field, _ := ii.outInfo.GetFieldByIndex(fieldIndex)
			err := ii.outInfo.SetFieldNull(field.Name)
			if err != nil {
				api.OutputMessage(ii.ToolId, api.Error, err.Error())
			}
		}
		record, err := ii.outInfo.GenerateRecord()
		if err != nil {
			api.OutputMessage(ii.ToolId, api.Error, err.Error())
		} else {
			ii.Output.PushRecord(record)
		}
	}

	api.OutputMessage(ii.ToolId, api.Complete, ``)
	ii.Output.Close()
}

func (ii *Ii) CacheSize() int {
	return 0
}
