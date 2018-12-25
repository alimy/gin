package yaml

import (
	"github.com/alimy/gin/binding"
	"github.com/alimy/gin/render"
)

func init() {
	binding.Register(binding.MIMEYAML, &yamlBinding{})
	render.Register(render.YAMLRenderFactory, &YAMLRenderFactory{})
}
