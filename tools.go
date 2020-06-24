package main

/*
   #include "tools.h"
*/
import "C"
import (
	"alteryx_streams_go/buffer"
	"alteryx_streams_go/combine_latest"
	"alteryx_streams_go/controller"
	"alteryx_streams_go/interval"
	"alteryx_streams_go/twitter"
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

//export Buffer
func Buffer(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &buffer.Plugin{}
	return C.long(api.ConfigurePlugin(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

//export Twitter
func Twitter(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &twitter.Plugin{}
	return C.long(api.ConfigurePlugin(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}
