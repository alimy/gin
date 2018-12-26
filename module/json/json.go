package json

import (
	"github.com/alimy/gin/binding"
	"github.com/alimy/gin/render"
)

func init() {
	binding.Register(binding.MIMEJSON, &jsonBinding{})

	render.Register(render.JSONRenderType, JSONRenderFactory{})
	render.Register(render.IntendedJSONRenderType, IndentedJSONRenderFactory{})
	render.Register(render.JsonpJSONRenderType, JsonpJSONRenderFactory{})
	render.Register(render.SecureJSONRenderType, SecureJSONRenderFactory{})
	render.Register(render.AsciiJSONRenderType, AsciiJSONRenderFactory{})
}

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// interface{} as a Number instead of as a float64.
var EnableDecoderUseNumber = false

// EnableJsonDecoderUseNumber sets true for binding.EnableDecoderUseNumberto to
// call the UseNumber method on the JSON Decoder instance.
func EnableJsonDecoderUseNumber() {
	EnableDecoderUseNumber = true
}
