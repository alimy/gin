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

// WriteContentType (XML) writes XML ContentType for response.
func (*XMLRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, xmlContentType)
}

func (XMLRenderFactory) Instance(data interface{}, opts ...interface{}) render.Render {
	return &XMLRender{Data: data}
}
