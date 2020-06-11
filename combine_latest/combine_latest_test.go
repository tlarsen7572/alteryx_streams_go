package combine_latest_test

import (
	"alteryx_streams_go/combine_latest"
	"fmt"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"testing"
)

func TestIiRecordInfos(t *testing.T) {
	plugin := &combine_latest.Plugin{}
	plugin.Init(1, ``)

	leftInfoXml := initIi(plugin, `Left`)
	rightInfoXml := initIi(plugin, `Right`)

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

func initIi(plugin *combine_latest.Plugin, connectionName string) string {
	ii, _ := plugin.AddIncomingConnection(``, connectionName)
	g := recordinfo.NewGenerator()
	g.AddInt64Field(fmt.Sprintf(`%vId`, connectionName), ``)
	info := g.GenerateRecordInfo()
	infoXml, _ := info.ToXml(connectionName)
	ii.Init(infoXml)
	return infoXml
}
