// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package yaml

import (
	"github.com/alimy/gin/render"
	"net/http"

	"gopkg.in/yaml.v2"
)

// YAML contains the given interface object.
type YAMLRender struct {
	Data interface{}
}

type YAMLRenderFactory struct{}

var yamlContentType = []string{"application/x-yaml; charset=utf-8"}

// Render (YAML) marshals the given interface object and writes data with custom ContentType.
func (r *YAMLRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := yaml.Marshal(r.Data)
	if err != nil {
		return err
	}

	w.Write(bytes)
	return nil
}

// WriteContentType (YAML) writes YAML ContentType for response.
func (*YAMLRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, yamlContentType)
}

func (YAMLRenderFactory) Instance(data interface{}, opts ...interface{}) render.Render {
	return &YAMLRender{Data: data}
}
