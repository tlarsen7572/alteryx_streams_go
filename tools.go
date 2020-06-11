package main

/*
   #include "tools.h"
*/
import "C"
import (
	"alteryx_streams_go/combine_latest"
	"alteryx_streams_go/controller"
	"alteryx_streams_go/interval"
	"github.com/tlarsen7572/goalteryx/api"
	"unsafe"
)

func main() {}

//export Controller
func Controller(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &controller.Controller{}
	return C.long(api.ConfigurePlugin(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

//export Interval
func Interval(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &interval.Plugin{}
	return C.long(api.ConfigurePlugin(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

//export CombineLatest
func CombineLatest(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &combine_latest.Plugin{}
	return C.long(api.ConfigurePlugin(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}
