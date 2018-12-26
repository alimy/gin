// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package json

import (
	"bytes"
	"fmt"
	"github.com/alimy/gin/render"
	"html/template"
	"net/http"

	"github.com/alimy/gin/module/json/internal/json"
)

// JSON contains the given interface object.
type JSONRender struct {
	Data interface{}
}

// IndentedJSON contains the given interface object.
type IndentedJSONRender struct {
	Data interface{}
}

// SecureJSON contains the given interface object and its prefix.
type SecureJSONRender struct {
	Prefix string
	Data   interface{}
}

// JsonpJSON contains the given interface object its callback.
type JsonpJSONRender struct {
	Callback string
	Data     interface{}
}

// AsciiJSON contains the given interface object.
type AsciiJSONRender struct {
	Data interface{}
}

// JSONRenderFactory contains the given interface object.
type JSONRenderFactory struct{}

// IndentedJSONRenderFactory contains the given interface object.
type IndentedJSONRenderFactory struct{}

// SecureJSONRenderFactory contains the given interface object.
type SecureJSONRenderFactory struct{}

// JsonpJSONRenderFactory contains the given interface object.
type JsonpJSONRenderFactory struct{}

// AsciiJSONRenderFactory contains the given interface object.
type AsciiJSONRenderFactory struct{}

// SecureJSONPrefix is a string which represents SecureJSON prefix.
type SecureJSONPrefix string

var jsonContentType = []string{"application/json; charset=utf-8"}
var jsonpContentType = []string{"application/javascript; charset=utf-8"}
var jsonAsciiContentType = []string{"application/json"}

// Render (JSON) writes data with custom ContentType.
func (r *JSONRender) Render(w http.ResponseWriter) (err error) {
	if err = WriteJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

// WriteContentType (JSON) writes JSON ContentType.
func (r *JSONRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, jsonContentType)
}

// Setup set data and opts
func (r *JSONRender) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *JSONRender) Reset() {
	r.Data = nil
}

// Instance a new Render instance
func (JSONRenderFactory) Instance() render.RenderRecycler {
	return &JSONRender{}
}

// WriteJSON marshals the given interface object and writes it with custom ContentType.
func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	render.WriteContentType(w, jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}

// Render (IndentedJSON) marshals the given interface object and writes it with custom ContentType.
func (r *IndentedJSONRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.MarshalIndent(r.Data, "", "    ")
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}

// WriteContentType (IndentedJSON) writes JSON ContentType.
func (*IndentedJSONRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, jsonContentType)
}

// Setup set data and opts
func (r *IndentedJSONRender) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *IndentedJSONRender) Reset() {
	r.Data = nil
}

// Instance a new Render instance
func (IndentedJSONRenderFactory) Instance() render.RenderRecycler {
	return &IndentedJSONRender{}
}

// Render (SecureJSON) marshals the given interface object and writes it with custom ContentType.
func (r *SecureJSONRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	// if the jsonBytes is array values
	if bytes.HasPrefix(jsonBytes, []byte("[")) && bytes.HasSuffix(jsonBytes, []byte("]")) {
		w.Write([]byte(r.Prefix))
	}
	w.Write(jsonBytes)
	return nil
}

// Setup set data and opts
func (r *SecureJSONRender) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
	if len(opts) == 1 {
		r.Prefix, _ = opts[0].(string)
	}
}

// Reset clean data and opts
func (r *SecureJSONRender) Reset() {
	r.Data = nil
	r.Prefix = ""
}

// WriteContentType (SecureJSON) writes JSON ContentType.
func (*SecureJSONRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, jsonContentType)
}

// Instance a new Render instance
func (SecureJSONRenderFactory) Instance() render.RenderRecycler {
	return &SecureJSONRender{}
}

// Render (JsonpJSON) marshals the given interface object and writes it and its callback with custom ContentType.
func (r *JsonpJSONRender) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	if r.Callback == "" {
		w.Write(ret)
		return nil
	}

	callback := template.JSEscapeString(r.Callback)
	w.Write([]byte(callback))
	w.Write([]byte("("))
	w.Write(ret)
	w.Write([]byte(")"))

	return nil
}

// Setup set data and opts
func (r *JsonpJSONRender) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
	if len(opts) == 1 {
		if callback, ok := opts[0].(string); ok {
			r.Callback = callback
		} else {
			r.Callback = ""
		}
	}
}

// Reset clean data and opts
func (r *JsonpJSONRender) Reset() {
	r.Data = nil
	r.Callback = ""
}

// WriteContentType (JsonpJSON) writes Javascript ContentType.
func (*JsonpJSONRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, jsonpContentType)
}

// Instance a new Render instance
func (JsonpJSONRenderFactory) Instance() render.RenderRecycler {
	return &JsonpJSONRender{}
}

// Render (AsciiJSON) marshals the given interface object and writes it with custom ContentType.
func (r *AsciiJSONRender) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	for _, r := range string(ret) {
		cvt := string(r)
		if r >= 128 {
			cvt = fmt.Sprintf("\\u%04x", int64(r))
		}
		buffer.WriteString(cvt)
	}

	w.Write(buffer.Bytes())
	return nil
}

// Setup set data and opts
func (r *AsciiJSONRender) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *AsciiJSONRender) Reset() {
	r.Data = nil
}

// WriteContentType (AsciiJSON) writes JSON ContentType.
func (*AsciiJSONRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, jsonAsciiContentType)
}

// Instance a new Render instance
func (AsciiJSONRenderFactory) Instance() render.RenderRecycler {
	return &AsciiJSONRender{}
}
