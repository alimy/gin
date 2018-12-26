// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protobuf

import (
	"github.com/alimy/gin/render"
	"net/http"

	"github.com/golang/protobuf/proto"
)

// ProtoBufRender contains the given interface object.
type ProtoBufRender struct {
	Data interface{}
}

type ProtoBufRenderFactory struct{}

var protobufContentType = []string{"application/x-protobuf"}

// Render (ProtoBuf) marshals the given interface object and writes data with custom ContentType.
func (r *ProtoBufRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := proto.Marshal(r.Data.(proto.Message))
	if err != nil {
		return err
	}

	w.Write(bytes)
	return nil
}

// Setup set data and opts
func (r *ProtoBufRender) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *ProtoBufRender) Reset() {
	r.Data = nil
}

// WriteContentType (ProtoBuf) writes ProtoBuf ContentType.
func (*ProtoBufRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, protobufContentType)
}

// Instance a new Render instance
func (ProtoBufRenderFactory) Instance() render.RenderRecycler {
	return &ProtoBufRender{}
}
