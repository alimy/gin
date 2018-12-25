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
func (*JSONRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, jsonContentType)
}

func (JSONRenderFactory) Instance(data interface{}, opts ...interface{}) render.Render {
	return &JSONRender{Data: data}
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

func (IndentedJSONRenderFactory) Instance(data interface{}, opts ...interface{}) render.Render {
	return &IndentedJSONRender{Data: data}
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

// WriteContentType (SecureJSON) writes JSON ContentType.
func (*SecureJSONRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, jsonContentType)
}

func (SecureJSONRenderFactory) Instance(data interface{}, opts ...interface{}) render.Render {
	render := &SecureJSONRender{Data: data}
	if len(opts) == 1 {
		render.Prefix, _ = opts[0].(string)
	}
	return render
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

// WriteContentType (JsonpJSON) writes Javascript ContentType.
func (*JsonpJSONRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, jsonpContentType)
}

func (JsonpJSONRenderFactory) Instance(data interface{}, opts ...interface{}) render.Render {
	render := &JsonpJSONRender{Data: data}
	if len(opts) == 1 {
		if callback, ok := opts[0].(string); ok {
			render.Callback = callback
		} else {
			render.Callback = ""
		}
	}
	return render
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

// WriteContentType (AsciiJSON) writes JSON ContentType.
func (*AsciiJSONRender) WriteContentType(w http.ResponseWriter) {
	render.WriteContentType(w, jsonAsciiContentType)
}

func (AsciiJSONRenderFactory) Instance(data interface{}, opts ...interface{}) render.Render {
	return &AsciiJSONRender{Data: data}
}
