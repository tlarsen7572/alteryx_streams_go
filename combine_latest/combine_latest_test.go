package combine_latest_test

import (
	"alteryx_streams_go/combine_latest"
	"fmt"
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"testing"
	"unsafe"
)

type iiMock struct {
	recordsPassed int
	pushCallback  func(pointer unsafe.Pointer)
	info          recordinfo.RecordInfo
}

func (i *iiMock) Init(recordInfoIn string) bool {
	i.info, _ = recordinfo.FromXml(recordInfoIn)
	return true
}

func (i *iiMock) PushRecord(record unsafe.Pointer) bool {
	i.pushCallback(record)
	return true
}

func (i *iiMock) UpdateProgress(percent float64) {}

func (i *iiMock) Close() {}

func (i *iiMock) CacheSize() int {
	return 0
}

func TestPushHasNoErrors(t *testing.T) {
	plugin := &combine_latest.Plugin{}
	plugin.Init(1, ``)
	mock := &iiMock{}
	mock.pushCallback = func(record unsafe.Pointer) {
		mock.recordsPassed++
		t.Logf(`got record`)
	}
	plugin.AddOutgoingConnection(``, api.NewConnectionInterfaceStruct(mock))
	leftIi, _ := initIi(plugin, `Left`)
	rightIi, _ := initIi(plugin, `Right`)
	leftIi.PushRecord(unsafe.Pointer(&[]byte{1, 0, 0, 0, 0, 0, 0, 0, 0}[0]))
	rightIi.PushRecord(unsafe.Pointer(&[]byte{2, 0, 0, 0, 0, 0, 0, 0, 0}[0]))

	if mock.recordsPassed != 1 {
		t.Fatalf(`expected 1 record but got %v`, mock.recordsPassed)
	}
}

func TestIiRecordInfos(t *testing.T) {
	plugin := &combine_latest.Plugin{}
	plugin.Init(1, ``)

	_, leftInfoXml := initIi(plugin, `Left`)
	_, rightInfoXml := initIi(plugin, `Right`)

	pluginLeft, err := plugin.LeftIn.ToXml(`Left`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	if pluginLeft != leftInfoXml {
		t.Fatalf(`expected '%v' but got '%v'`, leftInfoXml, pluginLeft)
	}

	pluginRight, err := plugin.RightIn.ToXml(`Right`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	if pluginRight != rightInfoXml {
		t.Fatalf(`expected '%v' but got '%v'`, rightInfoXml, pluginRight)
	}

	outG := recordinfo.NewGenerator()
	outG.AddInt64Field(`LeftId`, `Combine Latest`)
	outG.AddInt64Field(`RightId`, `Combine Latest`)
	out := outG.GenerateRecordInfo()
	expectedXml, _ := out.ToXml(`Output`)
	actualXml, _ := plugin.Out.ToXml(`Output`)
	if expectedXml != actualXml {
		t.Fatalf("the outgoing record info was not in the expected format.  Expected\n%v\nbut got\n%v", expectedXml, actualXml)
	}
}

func initIi(plugin *combine_latest.Plugin, connectionName string) (api.IncomingInterface, string) {
	ii, _ := plugin.AddIncomingConnection(connectionName, ``)
	g := recordinfo.NewGenerator()
	g.AddInt64Field(fmt.Sprintf(`%vId`, connectionName), ``)
	info := g.GenerateRecordInfo()
	infoXml, _ := info.ToXml(connectionName)
	ii.Init(infoXml)
	return ii, infoXml
}
