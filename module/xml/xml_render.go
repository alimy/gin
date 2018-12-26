// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package xml

import (
	"encoding/xml"
	"github.com/alimy/gin/render"
	"net/http"
)

// XMLRender contains the given interface object.
type XMLRender struct {
	Data interface{}
}

type XMLRenderFactory struct{}

var xmlContentType = []string{"application/xml; charset=utf-8"}

// Render (XML) encodes the given interface object and writes data with custom ContentType.
func (r *XMLRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return xml.NewEncoder(w).Encode(r.Data)
}

// Setup set data and opts
func (r *XMLRender) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *XMLRender) Reset() {
	r.Data = nil
}

// WriteContentType (XML) writes XML ContentType for response.
func (*XMLRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, xmlContentType)
}

// Instance a new Render instance
func (XMLRenderFactory) Instance() render.RenderRecycler {
	return &XMLRender{}
}
