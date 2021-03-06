// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build go1.7

package json

import (
	"github.com/alimy/gin/render"
	"net/http"

	"github.com/alimy/gin/module/json/internal/json"
)

// PureJSON contains the given interface object.
type PureJSONRender struct {
	Data interface{}
}

type PureJsonRenderFactory struct{}

func init() {
	render.Register(render.PureJSONRenderType, &PureJsonRenderFactory{})
}

// Render (PureJSON) writes custom ContentType and encodes the given interface object.
func (r *PureJSONRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(r.Data)
}

// Setup set data and opts
func (r *PureJSONRender) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *PureJSONRender) Reset() {
	r.Data = nil
}

// WriteContentType (PureJSON) writes custom ContentType.
func (*PureJSONRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, jsonContentType)
}

// Instance a new Render instance
func (PureJsonRenderFactory) Instance() render.RenderRecycler {
	return &PureJSONRender{}
}
