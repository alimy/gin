package protobuf

import (
	"github.com/alimy/gin/binding"
	"github.com/alimy/gin/render"
)

func init() {
	binding.Register(binding.MIMEPROTOBUF, &protobufBinding{})
	render.Register(render.ProtoBufRenderFactory, &ProtoBufRenderFactory{})
}
