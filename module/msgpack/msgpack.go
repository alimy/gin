package msgpack

import (
	"github.com/alimy/gin/binding"
	"github.com/alimy/gin/render"
)

func init() {
	binding.Register(binding.MIMEMSGPACK, &msgpackBinding{})
	render.Register(render.MsgPackRenderFactory, &MsgPackRenderFactory{})
}
