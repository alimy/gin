// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package msgpack

import (
	"github.com/alimy/gin/render"
	"net/http"

	"github.com/ugorji/go/codec"
)

// MsgPack contains the given interface object.
type MsgPackRender struct {
	Data interface{}
}

type MsgPackRenderFactory struct{}

var msgpackContentType = []string{"application/msgpack; charset=utf-8"}

// WriteContentType (MsgPack) writes MsgPack ContentType.
func (*MsgPackRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, msgpackContentType)
}

// Render (MsgPack) encodes the given interface object and writes data with custom ContentType.
func (r *MsgPackRender) Render(w http.ResponseWriter) error {
	return WriteMsgPack(w, r.Data)
}

// Setup set data and opts
func (r *MsgPackRender) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *MsgPackRender) Reset() {
	r.Data = nil
}

// WriteMsgPack writes MsgPack ContentType and encodes the given interface object.
func WriteMsgPack(w http.ResponseWriter, obj interface{}) error {
	render.WriteContentType(w, msgpackContentType)
	var mh codec.MsgpackHandle
	return codec.NewEncoder(w, &mh).Encode(obj)
}

// Instance a new Render instance
func (MsgPackRenderFactory) Instance() render.RenderRecycler {
	return &MsgPackRender{}
}
