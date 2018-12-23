package json

import (
	"github.com/alimy/gin/binding"
	"github.com/alimy/gin/render"
)

func init() {
	binding.Register(binding.MIMEPOSTForm, jsonBinding{})

	render.Register(render.JSONRenderFactory, &JSONRenderFactory{})
	render.Register(render.IntendedJSONRenderFactory, &IndentedJSONRenderFactory{})
	render.Register(render.JsonpJSONRenderFactory, &JsonpJSONRenderFactory{})
	render.Register(render.SecureJSONRenderFactory, &SecureJSONRenderFactory{})
	render.Register(render.AsciiJSONRenderFactory, &AsciiJSONRenderFactory{})
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