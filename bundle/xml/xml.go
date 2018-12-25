package xml

import (
	"github.com/alimy/gin/binding"
	"github.com/alimy/gin/render"
)

func init() {
	binding.Register(binding.MIMEXML, &xmlBinding{})
	render.Register(render.XMLRenderFactory, &XMLRenderFactory{})
}
